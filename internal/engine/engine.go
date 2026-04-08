package engine

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
	"context"
	"log"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type RouteEngine struct {
	mu     sync.RWMutex
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

func (e *RouteEngine) updateRoutes(ctx context.Context, force bool) {
	log.Printf("engine tick at %s", time.Now())
	now := time.Now()
	if !force {
		e.lastTick.Store(now)
	}

	var activeRoutes []*domain.Route

	for _, route := range e.Routes {
		// Skip if not in time range
		if !route.Schedule.ShouldRunNow(now) {
			continue
		}

		// Skip if route is fresh
		if !force && route.IsFresh(now, e.UpdateRate) {
			continue
		}

		activeRoutes = append(activeRoutes, route)
	}

	if len(activeRoutes) == 0 {
		return
	}

	measurements, err := e.computeRouteMatrix(ctx, activeRoutes)
	if err != nil {
		log.Printf("error computing matrix: %v\n", err)
		return
	}

	// Update in-memory routes
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, m := range measurements {
		route, ok := e.Routes[m.RouteID]
		if !ok {
			continue
		}
		route.DistanceMeters = &m.DistanceMeters
		route.DurationSeconds = &m.DurationSeconds
		route.RecordedAt = &m.RecordedAt
	}
}

func (e *RouteEngine) Run(ctx context.Context) {
	e.running.Store(true)
	defer e.running.Store(false)

	ticker := time.NewTicker(e.TickRate)
	defer ticker.Stop()

	log.Printf("Route engine started\n")
	e.updateRoutes(ctx, false)

	for {
		select {
		case <-ticker.C:
			e.updateRoutes(ctx, false)
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

func (e *RouteEngine) GetRoutes() []domain.Route {
	e.mu.RLock()
	defer e.mu.RUnlock()

	routes := make([]domain.Route, 0, len(e.Routes))
	for _, r := range e.Routes {
		routes = append(routes, *r)
	}
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].ID < routes[j].ID
	})
	return routes
}

func (e *RouteEngine) RefreshNow(ctx context.Context) error {
	e.updateRoutes(ctx, true)
	return nil
}
