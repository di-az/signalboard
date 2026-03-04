package engine

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type RouteEngine struct {
	Routes map[int]*domain.Route
	Store  *store.RouteStore

	UpdateRate time.Duration
	TickRate   time.Duration
	running    atomic.Bool
	lastTick   atomic.Value
	client     *http.Client
	apiKey     string
}

type Status struct {
	Running    bool   `json:"running"`
	TickRate   string `json:"tick_rate"`
	UpdateRate string `json:"update_rate"`
	Locations  int    `json:"locations"`
	LastTick   int    `json:"last_tick"`
}

func NewRouteEngine(
	ctx context.Context,
	store *store.RouteStore,
	updateRate time.Duration,
	tickRate time.Duration,
	apiKey string,
) (*RouteEngine, error) {
	routes, err := store.GetRoutesWithSchedules(ctx)
	if err != nil {
		return nil, err
	}

	routeMap := make(map[int]*domain.Route, len(routes))
	for _, r := range routes {
		routeMap[r.ID] = r
	}

	engine := &RouteEngine{
		Routes:     routeMap,
		Store:      store,
		UpdateRate: updateRate,
		TickRate:   tickRate,
		client:     &http.Client{Timeout: 5 * time.Second},
		apiKey:     apiKey,
	}

	return engine, nil
}

func (e *RouteEngine) checkLocations(ctx context.Context) {
	// log.Printf("engine tick at %s", time.Now())
	now := time.Now()
	e.lastTick.Store(now)

	var activeRoutes []*domain.Route

	for _, route := range e.Routes {
		// Skip if not in time range
		if !route.Schedule.ShouldRunNow(now) {
			continue
		}

		// Skip if route is fresh
		if route.IsFresh(now, e.UpdateRate) {
			continue
		}

		activeRoutes = append(activeRoutes, route)
	}

	if len(activeRoutes) == 0 {
		// log.Println("NO ACTIVE")
		return
	}

	// log.Println("HERE")
	// for _, r := range activeRoutes {
	// 	log.Println(*r)
	// }

	// routes, err := e.computeRouteMatrix(ctx, activeDestinations)
	// if err != nil {
	// 	log.Printf("error computing matrix: %v\n", err)
	// 	return
	// }
	//
	// for _, comm := range routes {
	// 	// activeDestinations[i].Schedule.LastUpdated = now
	// 	e.lastMeasured[comm.DestinationID] = now
	// 	log.Printf("Setting routes: %v\n", comm)
	// 	err := e.Store.Set(ctx, comm)
	// 	if err != nil {
	// 		log.Printf("error saving to db: %v\n", err)
	// 		return
	// 	}
	// }
}

func (e *RouteEngine) Run(ctx context.Context) {
	e.running.Store(true)
	defer e.running.Store(false)

	ticker := time.NewTicker(e.TickRate)
	defer ticker.Stop()

	log.Printf("Route engine started\n")
	e.checkLocations(ctx)

	for {
		select {
		case <-ticker.C:
			e.checkLocations(ctx)
		case <-ctx.Done():
			log.Println("Route engine shutting down")
			return
		}
	}
}

func (e *RouteEngine) Status() Status {
	return Status{
		Running:    e.running.Load(),
		TickRate:   e.TickRate.String(),
		UpdateRate: e.UpdateRate.String(),
		Locations:  len(e.Routes),
	}
}
