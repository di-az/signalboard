package commute

import (
	"fmt"
	"signalboard/internal/config"
	"time"
)

type Config struct {
	UpdateRate       time.Duration
	GoogleMapsAPIKey string
}

func LoadConfig() (*Config, error) {
	updateRate, err := config.GetDuration("COMMUTE_UPDATE_RATE", "10m")
	if err != nil {
		return nil, fmt.Errorf("Invalid COMMUTE_UPDATE_RATE: %w", err)
	}

	return &Config{
		GoogleMapsAPIKey: config.MustGet("COMMUTE_GOOGLE_MAPS_API_KEY"),
		UpdateRate:       updateRate,
	}, nil
}
