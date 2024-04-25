package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Elvilius/auto-catalog/internal/config"
	handler "github.com/Elvilius/auto-catalog/internal/handlers"
	"github.com/Elvilius/auto-catalog/internal/repo"
	"github.com/Elvilius/auto-catalog/internal/service"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		slog.Error("failed to connect to db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	err = goose.Up(db.DB, cfg.MigrationDir)
	if err != nil {
		slog.Error("failed to apply migrations: %v", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	httpService := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
	}

	repo := repo.NewRepo(db)
	service := service.NewService(repo, cfg)

	handler.Register(mux, service)

	go func() {
		if err := httpService.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run http server: %v", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server...")

	shutdownCtx, cancelShutdown := context.WithCancel(context.Background())
	defer cancelShutdown()

	if err := httpService.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error while shutting down server: %s", err)
		os.Exit(1)
	}

	slog.Info("Server gracefully stopped")
}
