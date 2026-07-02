package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config contains application-wide configuration settings
type Config struct {
	// SQLiteDB: Path to the SQLite database
	SQLiteDB string

	// TickRate: How often the engine executes
	TickRate time.Duration

	// // GoogleMapsAPIKey: Key used to use the Google Routes API
	// //TODO: REmove this from GeneralConfigModule
	// GoogleMapsAPIKey string
	// // UpdateRate: Rate at which the maps route will perform a new query to obtain the commute time
	// UpdateRate time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	tickRate, err := GetDuration("TICK_RATE", "1m")
	if err != nil {
		return nil, fmt.Errorf("Invalid TICK_RATE: %w", err)
	}

	SQLiteDB := GetString("SQLITE_PATH", "./routes.db")

	return &Config{
		SQLiteDB: SQLiteDB,
		TickRate: tickRate,
	}, nil
}

func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required env var: %s", key))
	}
	return val
}

func GetString(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func GetDuration(key, def string) (time.Duration, error) {
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
