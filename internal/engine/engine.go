package engine

import (
	"context"
	"log"
	"signalboard/internal/sources"
	"time"
)

type Engine struct {
	TickRate time.Duration
	Sources  []sources.Source
}

func NewEngine(
	tickRate time.Duration,
	sources ...sources.Source,
) *Engine {
	return &Engine{
		TickRate: tickRate,
		Sources:  sources,
	}
}

func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// initial refresh
	for _, source := range e.Sources {
		if err := source.Refresh(ctx); err != nil {
			log.Printf(
				"source %s refresh failed: %v",
				source.Name(),
				err,
			)
		}
	}

	for {
		select {
		case <-ticker.C:
			for _, source := range e.Sources {
				if err := source.Refresh(ctx); err != nil {
					log.Printf(
						"source %s refresh failed: %v",
						source.Name(),
						err,
					)
				}
			}

		case <-ctx.Done():
			log.Println("Engine shutting down")
			return
		}
	}
}
