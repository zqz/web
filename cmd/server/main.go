package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/zqz/web/backend/internal/config"
	"github.com/zqz/web/backend/internal/server"
)

func main() {
	logger := setupLogger()

	logger.Info().Msg("starting application")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	logger.Info().
		Str("env", cfg.Env).
		Int("port", cfg.Port).
		Msg("config loaded")

	ctx := context.Background()
	srv, err := server.New(ctx, cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to setup server")
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("server shutdown error")
		}
	}()

	logger.Info().Msg("database connected")
	logger.Info().Str("path", cfg.FilesPath).Msg("storage initialized")

	go func() {
		logger.Info().Str("address", cfg.Address()).Msg("starting HTTP server")
		if err := srv.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down server...")
}

func setupLogger() zerolog.Logger {
	if os.Getenv("ENV") == "production" {
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	return zerolog.New(output).With().Timestamp().Caller().Logger()
}
