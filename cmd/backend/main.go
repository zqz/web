package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zqz/upl/server"
)

func main() {
	var logger zerolog.Logger

	if os.Getenv("ZQZ_ENV") == "production" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	s, err := server.Init(&logger, "./config.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}
	defer s.Close()

	err = s.Run()
	if err != nil {
		logger.Fatal().Err(err).Msg("error running server")
	}

	logger.Info().Msg("ending")
}
