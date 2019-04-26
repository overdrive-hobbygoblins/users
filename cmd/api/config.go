package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Web struct {
		DebugHost       string        `envconfig:"API_HOST" default:":9090"`
		ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"2s"`
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
	}
}

func loadConfig() (cfg config, err error) {
	const app = ""

	if err := envconfig.Process(app, &cfg); err != nil {
		_ = envconfig.Usage(app, &cfg)
		return cfg, err
	}

	return cfg, nil
}
