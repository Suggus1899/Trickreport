package main

import (
	"context"
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

func main() {
	// ── Logger ────────────────────────────────────────────────────────
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"})

	// ── Config ────────────────────────────────────────────────────────
	cfg := config.Load()

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
	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// ── Migrations (embedded, versioned) ──────────────────────────────
	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	log.Info().Msg("Database migrations up to date")

	// ── Bootstrap initial admin (from ADMIN_EMAIL / ADMIN_PASSWORD) ───
	if err := bootstrap.EnsureAdmin(ctx, pool, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Error().Err(err).Msg("Admin bootstrap failed")
	}

	// ── Server ────────────────────────────────────────────────────────
	srv := httpServer.New(ctx, cfg, pool)
	srv.Start()
}
