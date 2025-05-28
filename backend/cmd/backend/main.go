package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/zqz/web/backend/internal/service"
)

func main() {
	boil.DebugMode = true
	env := os.Getenv("ZQZ_ENV")
	log := logger(env)

	s, err := service.NewProdServer(&log, env)
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

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	z := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	z.FormatFieldValue = func(i interface{}) string {
		if str, ok := i.(string); ok {
			// Replace escaped newlines in stack traces
			if strings.Contains(str, "\\n") {
				return strings.ReplaceAll(str, "\\n", "\n")
			}
		}
		return fmt.Sprintf("%s", i)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return zlog.Output(z)
}
