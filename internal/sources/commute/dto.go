package commute

import (
	"sort"
	"time"
)

type RouteResponse struct {
	ID              int        `json:"id"`
	Origin          string     `json:"origin"`
	Destination     string     `json:"destination"`
	DurationMinutes *int       `json:"duration_minutes"`
	DistanceKM      *float64   `json:"distance_km"`
	RecordedAt      *time.Time `json:"recorded_at"`
	ActiveNow       bool       `json:"active_now"`
}

func NewRouteResponse(route Route) RouteResponse {
	r := RouteResponse{
		ID:          route.ID,
		Origin:      route.Origin.Name,
		Destination: route.Destination.Name,
		ActiveNow:   route.Schedule.ShouldRunNow(time.Now()),
	}

	if route.DurationSeconds != nil {
		min := int(route.DurationSeconds.Minutes())
		r.DurationMinutes = &min
	}
	if route.DistanceMeters != nil {
		km := float64(*route.DistanceMeters) / 1000
		r.DistanceKM = &km
	}
	if route.RecordedAt != nil {
		r.RecordedAt = route.RecordedAt
	}

	return r
}

func SortRouteResponseSlice(routes []RouteResponse) {
	sort.Slice(routes, func(i, j int) bool {
		a := routes[i].DurationMinutes
		b := routes[j].DurationMinutes

		// Both nil → equal
		if a == nil && b == nil {
			return false
		}

		// Nil goes last
		if a == nil {
			return false
		}
		if b == nil {
			return true
		}

		return *a < *b
	})
}
