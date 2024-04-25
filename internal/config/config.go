package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort     int
	Debug        bool
	DatabaseDSN  string
	MigrationDir string
	CarInfoUrl   string
	MaxWorkers   int
}

func Load() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return &cfg, err
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return &cfg, err
	}

	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		return &cfg, err
	}
	parseBool, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return &cfg, err
	}

	cfg.MaxWorkers = maxWorkers
	cfg.CarInfoUrl = os.Getenv("CAR_INFO_URL")
	cfg.DatabaseDSN = os.Getenv("DATABASE_DSN")
	cfg.MigrationDir = os.Getenv("MIGRATION_DIR")
	cfg.HttpPort = httpPort
	cfg.Debug = parseBool
	return &cfg, nil
}
