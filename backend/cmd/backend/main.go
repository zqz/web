package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/zqz/web/backend/web"
)

func main() {
	boil.DebugMode = true
	env := os.Getenv("ZQZ_ENV")
	log := logger(env)

	s, err := web.Init(&log, env)
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
