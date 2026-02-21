package engine

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
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

func (e *RouteEngine) Run() {
	ticker := time.NewTicker(e.TickRate)
	defer ticker.Stop()

	// t := time.Now()
	// timeNow := time.Date(2026, 2, 21, 10, 30, 0, 0, t.Location())
	log.Printf("Route engine has started\n")

	for range ticker.C {
		log.Printf("Ticking")
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
			e.Store.Add(route)

			log.Printf("Route updated: %s -> %s (%d min)",
				e.Home.Name,
				location.Name,
				route.Minutes,
			)
		}
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
