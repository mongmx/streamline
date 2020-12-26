package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Config redis.
type Config struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	SSL     string
	SSLCert string
	SSLKey  string
}

// LoadEnv - load configuration from env.
func LoadEnv() Config {
	return Config{
		Host:    os.Getenv("POSTGRES_HOST"),
		Port:    os.Getenv("POSTGRES_PORT"),
		User:    os.Getenv("POSTGRES_USER"),
		Pass:    os.Getenv("POSTGRES_PASS"),
		DBName:  os.Getenv("POSTGRES_DB"),
		SSL:     os.Getenv("POSTGRES_SSL"),
		SSLCert: os.Getenv("POSTGRES_SSL_CERT"),
		SSLKey:  os.Getenv("POSTGRES_SSL_KEY"),
	}
}

// NewPostgres creates new connection to postgres and return the connection
func NewPostgres(cfg Config) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.Pass, cfg.SSL,
	)
	return sql.Open("postgres", conn)
}
