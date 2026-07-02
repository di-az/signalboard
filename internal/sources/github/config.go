package github

import (
	"signalboard/internal/config"
	"time"
)

const (
	DefaultHistoryDays = 30
	DefaultTimeout     = 10 * time.Second
)

type Config struct {
	Username string
	Token    string

	HistoryDays int
	Timeout     time.Duration
}

func LoadConfig() (*Config, error) {
	username := config.GetString("GITHUB_USERNAME", "di-az")
	token := config.GetString("GITHUB_TOKEN", "di-az")

	return &Config{
		Username:    username,
		Token:       token,
		HistoryDays: DefaultHistoryDays,
		Timeout:     DefaultTimeout,
	}, nil
}
