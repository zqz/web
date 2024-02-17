package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/zqz/upl/server"
)

func main() {
	env := os.Getenv("ZQZ_ENV")
	log := logger(env)

	s, err := server.Init(&log, env)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
	defer s.Close()

	err = s.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error running server")
	}

	log.Info().Msg("ending")
}

func logger(env string) zerolog.Logger {
	if env == "production" {
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return zlog.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
