package engine

import (
	"fmt"
	"log"
	"time"

	"commuteboard/internal/domain"
)

func mapMatrixElements(
	routes []*domain.Route,
	elements []RouteMatrixElement,
	now time.Time,
) ([]domain.RouteMeasurement, error) {
	var measurements []domain.RouteMeasurement

	for _, el := range elements {
		// Ignore cross-routes
		if el.OriginIndex != el.DestinationIndex {
			continue
		}

		idx := el.OriginIndex

		// Check for valid response body matrix
		if idx < 0 || idx >= len(routes) {
			log.Printf("invalid index from matrix response: %d", idx)
			continue
		}
		log.Printf(
			"Matrix element: origin=%d dest=%d duration=%s distance=%d",
			el.OriginIndex,
			el.DestinationIndex,
			el.Duration,
			el.DistanceMeters,
		)

		duration, err := time.ParseDuration(el.Duration)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse duration %q: %w",
				el.Duration,
				err,
			)
		}

		r := routes[idx]

		measurements = append(measurements, domain.RouteMeasurement{
			RouteID:         r.ID,
			DistanceMeters:  el.DistanceMeters,
			DurationSeconds: duration,
			RecordedAt:      now,
		})
	}

	return measurements, nil
}
