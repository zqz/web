package server

import (
	"database/sql"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/friendsofgo/errors"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type config struct {
	Port        int    `env:"PORT" envDefault:"3000"`
	Env         string `env:"ZQZ_ENV" envDefault:"development"`
	DatabaseURL string `env:"DATABASE_URL"`
	FilesPath   string `env:"FILES_PATH" envDefault:"./files"`
}

func (c config) isProduction() bool {
	return c.Env == "production"
}

func (c config) isDevelopment() bool {
	return len(c.Env) == 0 || c.Env == "development"
}

func loadConfig() (config, error) {
	cfgd := config{Env: os.Getenv("ZQZ_ENV")}

	if cfgd.isDevelopment() {
		err := godotenv.Load()
		if err != nil {
			return cfgd, err
		}
	}

	if err := env.Parse(&cfgd); err != nil {
		return cfgd, err
	}

	return cfgd, nil
}

func openDatabase(url string) (*sql.DB, error) {
	openStr, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", openStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to establish connection to db")
	}

	return db, nil
}
