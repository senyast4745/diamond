package main

import (
	"github.com/caarlos0/env/v6"
)

type config struct {
	Listen   string `env:"LISTEN" envDefault:"localhost:7171"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFmt   string `env:"LOG_FMT" envDefault:"console"`

	DbAddr string `env:"DB_HOST" envDefault:"mongodb://localhost:27017/"`
	DbName string `env:"DB_NAME" envDefault:"diamond"`
}

func initConfig() (*config, error) {
	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
