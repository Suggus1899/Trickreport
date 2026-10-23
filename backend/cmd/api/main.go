package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/trickreport/backend/internal/bootstrap"
	"github.com/trickreport/backend/internal/config"
	"github.com/trickreport/backend/internal/db"
	httpServer "github.com/trickreport/backend/internal/interfaces/http"
)

// @title           Trickreport API
// @version         1.0
// @description     Multi-tenant helpdesk / ticketing platform API.
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
// @description Bearer JWT token (or trickreport_token cookie).

var Version = "dev"
var Commit = "unknown"
var CommitDate = "unknown"

func main() {
	// ── Flags ─────────────────────────────────────────────────────────
	migrateUp := flag.Bool("migrate-up", false, "Apply all pending database migrations, then exit")
	migrateDown := flag.Bool("migrate-down", false, "Roll back the last database migration, then exit")
	seed := flag.Bool("seed", false, "Seed the database with demo data, then exit (refuses to run in production)")
	version := flag.Bool("version", false, "Print version and source commit, then exit")
	flag.Parse()
	if *version {
		fmt.Printf("Trickreport %s (%s, %s)\n", Version, Commit, CommitDate)
		return
	}

	// ── Config ────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Logger ────────────────────────────────────────────────────────
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// ── One-off maintenance commands ─────────────────────────────────
	// These apply/roll back migrations against the raw DSN and exit —
	// no server, no connection pool needed.
	if *migrateUp {
		if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
			log.Fatal().Err(err).Msg("Failed to run database migrations")
		}
		log.Info().Msg("Database migrations applied")
		return
	}
	if *migrateDown {
		if err := db.RollbackLastMigration(cfg.DatabaseURL); err != nil {
			log.Fatal().Err(err).Msg("Failed to roll back migration")
		}
		log.Info().Msg("Last migration rolled back")
		return
	}

	// ── Context ───────────────────────────────────────────────────────
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info().Msg("Shutdown signal received")
		cancel()
	}()

	// ── Database ──────────────────────────────────────────────────────
	pool, err := db.NewPool(ctx, cfg.DatabaseURL, db.DefaultPoolConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// ── Migrations (embedded, versioned) ──────────────────────────────
	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	log.Info().Msg("Database migrations up to date")

	if *seed {
		if err := bootstrap.Seed(ctx, pool, cfg.IsProduction()); err != nil {
			log.Fatal().Err(err).Msg("Seed failed")
		}
		log.Info().Msg("Seed complete")
		return
	}

	// ── Bootstrap initial admin (from ADMIN_EMAIL / ADMIN_PASSWORD) ───
	if err := bootstrap.EnsureAdmin(ctx, pool, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Error().Err(err).Msg("Admin bootstrap failed")
	}

	// ── Bootstrap the fixed system user (automation/worker attribution) ──
	if err := bootstrap.EnsureSystemUser(ctx, pool); err != nil {
		log.Error().Err(err).Msg("System user bootstrap failed")
	}

	// ── Server ────────────────────────────────────────────────────────
	srv := httpServer.New(ctx, cfg, pool)
	srv.Start()
}
