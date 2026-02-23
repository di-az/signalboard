package engine

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
	"context"
	"log"
	"time"
)

type RouteEngine struct {
	Home       domain.Location
	Locations  []*domain.Location
	Store      *store.RouteStore
	UpdateRate time.Duration
	TickRate   time.Duration
}

type Status struct {
	TickRate   string `json:"tick_rate"`
	UpdateRate string `json:"update_rate"`
	Locations  int    `json:"locations"`
}

func NewRouteEngine(
	home domain.Location,
	locations []*domain.Location,
	store *store.RouteStore,
	updateRate time.Duration,
	tickRate time.Duration,
) *RouteEngine {
	return &RouteEngine{
		Home:       home,
		Locations:  locations,
		Store:      store,
		UpdateRate: updateRate,
		TickRate:   tickRate,
	}
}

func (e *RouteEngine) checkLocations() {
	log.Printf("engine tick at %s", time.Now())
	now := time.Now()

	for _, location := range e.Locations {
		// Skip if not in time range
		if !location.Schedule.ShouldRunNow(now) {
			continue
		}

		// Skip if recently updated
		if now.Sub(location.Schedule.LastUpdated) < e.UpdateRate {
			continue
		}

		route, err := getRoute(e.Home, *location)
		if err != nil {
			log.Printf("Error calculating route for %s: %v", location.Name, err)
			continue
		}

		location.Schedule.LastUpdated = now
		e.Store.Set(route)

		log.Printf("Route updated: %s -> %s (%d min)",
			e.Home.Name,
			location.Name,
			route.Minutes,
		)
	}
}

func (e *RouteEngine) Run(ctx context.Context) {
	ticker := time.NewTicker(e.TickRate)
	defer ticker.Stop()

	log.Printf("Route engine started\n")
	e.checkLocations()

	for {
		select {
		case <-ticker.C:
			e.checkLocations()
		case <-ctx.Done():
			return
		}
	}
}

func (e *RouteEngine) Status() Status {
	return Status{
		TickRate:   e.TickRate.String(),
		UpdateRate: e.UpdateRate.String(),
		Locations:  len(e.Locations),
	}
}

// TODO: Temporary Function
func getRoute(start, finish domain.Location) (domain.Route, error) {
	route := domain.Route{
		Start:        start,
		Finish:       finish,
		Minutes:      18,
		TrafficLevel: "Low",
		Timestamp:    time.Now(),
	}
	return route, nil
}
