package store

import (
	"commuteboard/internal/db"
	"commuteboard/internal/domain"
	"context"
	"fmt"
	"time"
)

type RouteStore struct {
	db *db.SQLite
}

func NewRouteStore(db *db.SQLite) *RouteStore {
	return &RouteStore{db: db}
}

func (s *RouteStore) GetAllRoutes(ctx context.Context) ([]*domain.Route, error) {
	routeRows, err := s.db.GetRouteRows(ctx)
	if err != nil {
		return nil, err
	}

	locationRows, err := s.db.GetLocations(ctx)
	if err != nil {
		return nil, err
	}
	locations := locationRowsToDomain(locationRows)

	locationMap := make(map[int]*domain.Location, len(locations))
	for _, loc := range locations {
		locationMap[loc.ID] = loc
	}

	routes := make([]*domain.Route, 0, len(routeRows))
	for _, row := range routeRows {
		origin := locationMap[row.OriginID]
		destination := locationMap[row.DestinationID]

		if origin == nil || destination == nil {
			return nil, fmt.Errorf(
				"route %d references missing location (origin=%d destination=%d)",
				row.ID,
				row.OriginID,
				row.DestinationID,
			)
		}

		r := &domain.Route{
			ID:              row.ID,
			Origin:          origin,
			Destination:     destination,
			DistanceMeters:  row.DistanceMeters,
			DurationSeconds: time.Duration(row.DurationSeconds) * time.Second,
			RecordedAt:      row.RecordedAt,
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func (s *RouteStore) GetRoutesWithSchedules(ctx context.Context) ([]*domain.Route, error) {
	routes, err := s.GetAllRoutes(ctx)
	if err != nil {
		return nil, err
	}

	scheduleRows, err := s.db.GetSchedules(ctx)
	if err != nil {
		return nil, err
	}

	scheduleMap, err := db.ScheduleRowsToDomain(scheduleRows)
	if err != nil {
		return nil, err
	}

	for _, r := range routes {
		if sched, ok := scheduleMap[r.ID]; ok {
			r.Schedule = sched
		}
	}

	return routes, nil
}

func locationRowsToDomain(rows []db.LocationRow) []*domain.Location {
	result := make([]*domain.Location, 0, len(rows))

	for _, row := range rows {
		result = append(result, &domain.Location{
			ID:        row.ID,
			Name:      row.Name,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		})
	}

	return result
}

func (s *RouteStore) UpdateMeasurements(
	ctx context.Context,
	routes []domain.RouteMeasurement,
) error {

	rows := toMeasurementRows(routes)

	return s.db.UpdateRouteMeasurements(ctx, rows)
}

func toMeasurementRows(routes []domain.RouteMeasurement) []db.RouteMeasurementRow {
	rows := make([]db.RouteMeasurementRow, 0, len(routes))

	for _, r := range routes {
		rows = append(rows, db.RouteMeasurementRow{
			ID:              r.RouteID,
			DistanceMeters:  r.DistanceMeters,
			DurationSeconds: int(r.DurationSeconds.Seconds()),
			RecordedAt:      r.RecordedAt,
		})
	}

	return rows
}
