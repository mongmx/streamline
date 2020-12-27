package config

import (
	"sync"

	"os"

	"github.com/mongmx/streamline/config/postgres"
	"github.com/mongmx/streamline/config/redis"
)

// Cfg is a configuration variable for the app.
var Cfg Config
var once sync.Once

// Config is a general configuration.
type Config struct {
	Mode     string
	Port     string
	Debug    string
	Postgres postgres.Config
	Redis    redis.Config
}

// LoadEnv loads configuration from env variables.
func LoadEnv() {
	once.Do(func() {
		Cfg = Config{
			Mode:     os.Getenv("MODE"),
			Port:     os.Getenv("PORT"),
			Debug:    os.Getenv("DEBUG"),
			Postgres: postgres.LoadEnv(),
			Redis:    redis.LoadEnv(),
		}
	})
}
