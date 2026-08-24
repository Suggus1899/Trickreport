package main

import (
	"context"
	"flag"
	"io"
	"os"
	"os/signal"
	"path/filepath"
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

func main() {
	// ── Flags ─────────────────────────────────────────────────────────
	migrateUp := flag.Bool("migrate-up", false, "Apply all pending database migrations, then exit")
	migrateDown := flag.Bool("migrate-down", false, "Roll back the last database migration, then exit")
	seed := flag.Bool("seed", false, "Seed the database with demo data, then exit (refuses to run in production)")
	flag.Parse()

	// ── Config ────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Logger ────────────────────────────────────────────────────────
	// Single log file: backend/trickreport.log (no logs/ folder).
	// In development: also write to stdout with console formatting.
	// In production: write JSON to the file only.
	logPath := filepath.Join(".", "trickreport.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback to stdout-only if the file can't be opened.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"})
		log.Warn().Err(err).Msg("failed to open log file, falling back to stdout only")
	} else {
		defer logFile.Close()
		fileWriter := zerolog.New(logFile).With().Timestamp().Logger()
		if cfg.IsProduction() {
			// Production: JSON to file only.
			log.Logger = fileWriter
		} else {
			// Development: console to stdout + JSON to file (multi-writer).
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
			log.Logger = log.Output(io.MultiWriter(consoleWriter, fileWriter))
		}
	}

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
