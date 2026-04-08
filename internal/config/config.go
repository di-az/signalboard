package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLiteDB string

	GoogleMapsAPIKey string
	UpdateRate       time.Duration
	TickRate         time.Duration
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	// Required
	cfg.GoogleMapsAPIKey = mustGet("GOOGLE_MAPS_API_KEY")

	// Optional with defaults
	cfg.SQLiteDB = getString("SQLITE_PATH", "./routes.db")

	cfg.UpdateRate, err = getDuration("UPDATE_RATE", "10m")
	if err != nil {
		return nil, fmt.Errorf("invalid UPDATE_RATE: %w", err)
	}

	cfg.TickRate, err = getDuration("TICK_RATE", "10s")
	if err != nil {
		return nil, fmt.Errorf("invalid TICK_RATE: %w", err)
	}

	return cfg, nil
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required env var: %s", key))
	}
	return val
}

func getString(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getDuration(key, def string) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for key %s: %s", key, val)
	}

	return duration, nil
}
