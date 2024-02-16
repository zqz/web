package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type databaseConfig struct {
	Host string `json:"host"`
	Name string `json:"name"`
	Pass string `json:"pass"`
	Port int    `json:"port"`
	User string `json:"user"`
}

type config struct {
	Port     int            `json:"port"`
	DBConfig databaseConfig `json:"database"`
}

func (dc databaseConfig) openstring() string {
	s := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		dc.Host, dc.Port, dc.User, dc.Name,
	)

	if dc.Pass != "" {
		s += " password=" + dc.Pass
	}

	return s
}

func (dc databaseConfig) loadDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", dc.openstring())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func parseConfig(path string) (config, error) {
	cfg := config{}

	if path == "" {
		return cfg, errors.New("path is empty")
	}

	b, err := os.ReadFile(path)

	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
