package server

import (
	"commuteboard/internal/domain"
	"time"
)

type CommuteResponse struct {
	OriginID        int       `json:"origin_id"`
	OriginName      string    `json:"origin_name"`
	DestinationID   int       `json:"destination_id"`
	DestinationName string    `json:"destination_name"`
	DurationMinutes int       `json:"duration_minutes"`
	DistanceKM      float64   `json:"distance_km"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewCommuteResponse(origin domain.Location, destination domain.Location, route domain.Route) CommuteResponse {
	return CommuteResponse{
		OriginID:        origin.ID,
		OriginName:      origin.Name,
		DestinationID:   destination.ID,
		DestinationName: destination.Name,
		DurationMinutes: int(route.DurationSeconds.Minutes()),
		DistanceKM:      float64(route.DistanceMeters) / 1000,
		UpdatedAt:       route.RecordedAt,
	}
}
