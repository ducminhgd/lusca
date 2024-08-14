package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	Cfg    = Load()
)

type Config struct {
	API       APIConfig
	Databases DBConfig `env:",prefix=DB_"`
}

type APIConfig struct {
	Host string `env:"SERVER_HOST, default=0.0.0.0"`
	Port int    `env:"SERVER_PORT, default=8000"`
}

type DBConfig struct {
	UseReplication  bool          `env:"USE_REPLICATION, default=false"`
	Source          DBAttributes  `env:",prefix=SOURCE_"`
	Replica         DBAttributes  `env:",prefix=REPLICA_"`
	LogLevel        int           `env:"LOG_LEVEL, default=2"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME, default=300s"`
	ConnMaxLifeTime time.Duration `env:"CONN_MAX_LIFE_TIME, default=300s"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS, default=5"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS, default=10"`
}

type DBAttributes struct {
	Dsn      string `env:"DSN"`
	Host     string `env:"HOST, default=127.0.0.1"`
	Port     int    `env:"PORT, default=5432"`
	Database string `env:"DATABASE, default=lusca"`
	Username string `env:"USERNAME, default=postgres_user"`
	Password string `env:"PASSWORD, default=postgres_pw"`
	SSLMode  string `env:"SSL_MODE, default=disable"`
}

func (c DBAttributes) GetPostgresDSN() string {
	if c.Dsn != "" {
		return c.Dsn
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", c.Username, c.Password, c.Database, c.Host, c.Port, c.SSLMode)
}

func Load() Config {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		logger.ErrorContext(ctx, "Failed to load config", slog.Any("error", err))
	}
	return c
}
