package store

import (
	"commuteboard/internal/domain"
	"sync"
)

type RouteStore struct {
	mu     sync.RWMutex
	routes map[string]domain.Route
}

func NewRouteStore() *RouteStore {
	return &RouteStore{routes: make(map[string]domain.Route)}
}

func (s *RouteStore) Set(route domain.Route) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := route.Finish.Name
	s.routes[key] = route
}

func (s *RouteStore) GetAll() []domain.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Route, 0, len(s.routes))
	for _, r := range s.routes {
		result = append(result, r)
	}
	return result

}
