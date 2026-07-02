package commute

import (
	"context"
	"fmt"
	"signalboard/internal/db"
)

type CommuteStore struct {
	db *CommuteDb
}

func NewCommuteStore(db *db.SQLite) *CommuteStore {
	return &CommuteStore{db: NewCommuteDb(db)}
}

func (s *CommuteStore) GetAllRoutes(ctx context.Context) ([]*Route, error) {
	routeRows, err := s.db.GetRouteRows(ctx)
	if err != nil {
		return nil, err
	}

	locationRows, err := s.db.GetLocations(ctx)
	if err != nil {
		return nil, err
	}
	locations := locationRowsToDomain(locationRows)

	locationMap := make(map[int]*Location, len(locations))
	for _, loc := range locations {
		locationMap[loc.ID] = loc
	}

	routes := make([]*Route, 0, len(routeRows))
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

		r := &Route{
			ID:          row.ID,
			Origin:      origin,
			Destination: destination,
		}

		if row.DistanceMeters != nil {
			r.DistanceMeters = row.DistanceMeters
		}

		if row.DurationSeconds != nil {
			r.DurationSeconds = row.DurationSeconds
		}

		if row.RecordedAt != nil {
			r.RecordedAt = row.RecordedAt
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func (s *CommuteStore) GetRoutesWithSchedules(ctx context.Context) ([]*Route, error) {
	routes, err := s.GetAllRoutes(ctx)
	if err != nil {
		return nil, err
	}

	scheduleRows, err := s.db.GetSchedules(ctx)
	if err != nil {
		return nil, err
	}

	scheduleMap, err := ScheduleRowsToDomain(scheduleRows)
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

func locationRowsToDomain(rows []LocationRow) []*Location {
	result := make([]*Location, 0, len(rows))

	for _, row := range rows {
		result = append(result, &Location{
			ID:        row.ID,
			Name:      row.Name,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		})
	}

	return result
}

func (s *CommuteStore) UpdateMeasurements(
	ctx context.Context,
	routes []RouteMeasurement,
) error {

	rows := toMeasurementRows(routes)

	return s.db.UpdateRouteMeasurements(ctx, rows)
}

func toMeasurementRows(routes []RouteMeasurement) []RouteMeasurementRow {
	rows := make([]RouteMeasurementRow, 0, len(routes))

	for _, r := range routes {
		rows = append(rows, RouteMeasurementRow{
			ID:              r.RouteID,
			DistanceMeters:  r.DistanceMeters,
			DurationSeconds: int(r.DurationSeconds.Seconds()),
			RecordedAt:      r.RecordedAt,
		})
	}

	return rows
}
