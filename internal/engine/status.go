package engine

import (
	"time"
)

type Status struct {
	Sources  []SourceStatus `json:"sources"`
	TickRate string         `json:"tick_rate"`

	Running  bool   `json:"running"`
	LastTick string `json:"last_tick"`
}

type SourceStatus struct {
	Name string `json:"name"`
}

func (e *Engine) Status() Status {
	lastTick, _ := e.lastTick.Load().(time.Time)

	var sourceStatuses []SourceStatus

	for _, source := range e.Sources {
		sourceStatuses = append(sourceStatuses, SourceStatus{Name: source.Name()})
	}

	var lastTickStr string
	if !lastTick.IsZero() {
		lastTickStr = lastTick.Format(time.RFC3339)
	}

	return Status{
		Running:  e.running.Load(),
		TickRate: e.TickRate.String(),
		LastTick: lastTickStr,
		Sources:  sourceStatuses,
	}
}
