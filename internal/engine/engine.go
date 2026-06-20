package engine

import (
	"context"
	"log"
	"signalboard/internal/sources"
	"sync/atomic"
	"time"
)

type Engine struct {
	TickRate time.Duration
	Sources  []sources.Source
	running  atomic.Bool
	lastTick atomic.Value
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
	e.running.Store(true)
	defer e.running.Store(false)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	now := time.Now()
	e.lastTick.Store(now)

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
			e.lastTick.Store(time.Now())
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
