package postgres

import (
	"app/config"
	"app/pkg/crosscutting/errors"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // "disable", "require", "verify-full"

	MaxOpenConnections *int
	MaxIdleConnections *int
	ConnMaxLifetime    *time.Duration
	ConnMaxIdleTime    *time.Duration
}

func (conf PostgresConfig) GetMaxOpenConnections() int {
	if conf.MaxOpenConnections != nil {
		return *conf.MaxOpenConnections
	}
	return 25
}

func (conf PostgresConfig) GetMaxIdleConnections() int {
	if conf.MaxIdleConnections != nil {
		return *conf.MaxIdleConnections
	}
	return 5
}

func (conf PostgresConfig) GetConnMaxLifetime() time.Duration {
	if conf.ConnMaxLifetime != nil {
		return *conf.ConnMaxLifetime
	}
	return 5 * time.Minute
}

func (conf PostgresConfig) GetConnMaxIdleTime() time.Duration {
	if conf.ConnMaxIdleTime != nil {
		return *conf.ConnMaxIdleTime
	}
	return 1 * time.Minute
}

// LoadPostgresConfigFromEnv loads config from environment variables
func LoadPostgresConfigFromEnv() PostgresConfig {
	return PostgresConfig{
		Host:     config.GetDBHost(),
		Port:     config.GetDBPort(),
		User:     config.GetDBUser(),
		Password: config.GetDBPassword(),
		DBName:   config.GetDBName(),
		SSLMode:  config.GetDBSSLMode(),
	}
}

// DSN returns the connection string
func (conf PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBName, conf.SSLMode,
	)
}

// URL returns the connection URL format
func (conf PostgresConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, conf.SSLMode,
	)
}

// NewPostgresDB creates a new sql.DB instance for PostgreSQL
func NewPostgresDB(config PostgresConfig) (*sql.DB, error) {
	// For pgx driver use "pgx", for lib/pq use "postgres"
	db, err := sql.Open("pgx", config.DSN())
	if err != nil {
		return nil, errors.Lift(err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.GetMaxOpenConnections())
	db.SetMaxIdleConns(config.GetMaxIdleConnections())
	db.SetConnMaxLifetime(config.GetConnMaxLifetime())
	db.SetConnMaxIdleTime(config.GetConnMaxIdleTime())

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, errors.Lift(err)
	}

	return db, nil
}

// NewPostgresDBFromEnv creates sql.DB using environment variables
func NewPostgresDBFromEnv() (*sql.DB, error) {
	return NewPostgresDB(
		LoadPostgresConfigFromEnv(),
	)
}
