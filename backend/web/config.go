package web

import (
	"database/sql"

	"github.com/caarlos0/env/v10"
	"github.com/friendsofgo/errors"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type Config struct {
	Port        int    `env:"PORT" envDefault:"3000"`
	Env         string `env:"ZQZ_ENV" envDefault:"development"`
	DatabaseURL string `env:"DATABASE_URL"`
	FilesPath   string `env:"FILES_PATH" envDefault:"./files"`
}

func (c Config) isProduction() bool {
	return c.Env == "production"
}

func (c Config) isDevelopment() bool {
	return len(c.Env) == 0 || c.Env == "development"
}

func LoadConfig(currentEnv string) (Config, error) {
	cfgd := Config{Env: currentEnv}

	if cfgd.isDevelopment() {
		godotenv.Load()
	}

	if err := env.Parse(&cfgd); err != nil {
		return cfgd, err
	}

	return cfgd, nil
}

func OpenDatabase(url string) (*sql.DB, error) {
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
