package commute

import (
	"context"
	"log"
	"net/http"
	"signalboard/internal/sources"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type CommuteSource struct {
	mu     sync.RWMutex
	Store  *CommuteStore
	Routes map[int]*Route

	UpdateRate time.Duration
	TickRate   time.Duration
	running    atomic.Bool
	lastTick   atomic.Value
	client     *http.Client
	apiKey     string
}

func NewCommuteSource(
	ctx context.Context,
	store *CommuteStore,
	updateRate time.Duration,
	tickRate time.Duration,
	apiKey string,
) (*CommuteSource, error) {
	routes, err := store.GetRoutesWithSchedules(ctx)
	if err != nil {
		return nil, err
	}

	routeMap := make(map[int]*Route, len(routes))
	for _, r := range routes {
		routeMap[r.ID] = r
	}

	source := &CommuteSource{
		Routes:     routeMap,
		Store:      store,
		UpdateRate: updateRate,
		TickRate:   tickRate,
		client:     &http.Client{Timeout: 5 * time.Second},
		apiKey:     apiKey,
	}

	return source, nil
}

func (s *CommuteSource) Name() string {
	return "commute"
}

func (s *CommuteSource) Refresh(ctx context.Context) error {
	s.updateRoutes(ctx, false)
	return nil
}

func (s *CommuteSource) ForceRefresh(ctx context.Context) error {
	s.updateRoutes(ctx, true)
	return nil
}

func (s *CommuteSource) Endpoints() []sources.Endpoint {
	return []sources.Endpoint{
		{Method: "GET", Path: "", Handler: s.GetRoutesHandler()},
		{Method: "GET", Path: "/active", Handler: s.GetActiveRoutesHandler()},
		{Method: "POST", Path: "/refresh", Handler: s.RefreshRoutes()},
	}
}

func (s *CommuteSource) GetRoutes() []Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	routes := make([]Route, 0, len(s.Routes))
	for _, r := range s.Routes {
		routes = append(routes, *r)
	}
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].ID < routes[j].ID
	})
	return routes
}

func (s *CommuteSource) updateRoutes(ctx context.Context, force bool) {
	log.Printf("engine tick at %s", time.Now())
	now := time.Now()
	if !force {
		s.lastTick.Store(now)
	}

	var activeRoutes []*Route

	for _, route := range s.Routes {
		// Skip if not in time range
		if !route.Schedule.ShouldRunNow(now) {
			continue
		}

		// Skip if route is fresh
		if !force && route.IsFresh(now, s.UpdateRate) {
			continue
		}

		activeRoutes = append(activeRoutes, route)
	}

	if len(activeRoutes) == 0 {
		return
	}

	measurements, err := s.computeRouteMatrix(ctx, activeRoutes)
	if err != nil {
		log.Printf("error computing matrix: %v\n", err)
		return
	}

	// Update in-memory routes
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range measurements {
		route, ok := s.Routes[m.RouteID]
		if !ok {
			continue
		}
		route.DistanceMeters = &m.DistanceMeters
		route.DurationSeconds = &m.DurationSeconds
		route.RecordedAt = &m.RecordedAt
	}
}
