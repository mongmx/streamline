package postgres

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"os"
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
func NewPostgres(cfg Config) string {
	if cfg.Pass == "" {
		return fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.SSL,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.Pass, cfg.SSL,
	)
}

// MigrateUp call migration up
func MigrateUp(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/postgres",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
