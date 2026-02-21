package store

import (
	"commuteboard/internal/domain"
	"sync"
)

type RouteStore struct {
	mu     sync.RWMutex
	routes []domain.Route
}

func NewRouteStore() *RouteStore {
	return &RouteStore{}
}

func (s *RouteStore) Add(route domain.Route) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.routes = append(s.routes, route)

	// Optional: keep only last 100
	// if len(s.routes) > 100 {
	// 	s.routes = s.routes[len(s.routes)-100:]
	// }
}

func (s *RouteStore) GetAll() []domain.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy to avoid race issues
	cpy := make([]domain.Route, len(s.routes))
	copy(cpy, s.routes)
	return cpy
}
