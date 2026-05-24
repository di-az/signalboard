package engine

import (
	"bytes"
	"commuteboard/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var TRAVEL_MODE = "DRIVE"
var ROUTING_PREFERENCE = "TRAFFIC_AWARE_OPTIMAL"
var GOOGLE_ENDPOINT = "https://routes.googleapis.com/distanceMatrix/v2:computeRouteMatrix"

type matrixRequest struct {
	Origins           []matrixOrigin      `json:"origins"`
	Destinations      []matrixDestination `json:"destinations"`
	TravelMode        string              `json:"travelMode"`
	RoutingPreference string              `json:"routingPreference"`
	// DepartureTime     string              `json:"departureTime"`
}

type matrixOrigin struct {
	Waypoint matrixWaypoint `json:"waypoint"`
}

type matrixDestination struct {
	Waypoint matrixWaypoint `json:"waypoint"`
}

type matrixWaypoint struct {
	Location matrixLocation `json:"location"`
}

type matrixLocation struct {
	LatLng matrixLatLng `json:"latLng"`
}

type matrixLatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RouteMatrixElement struct {
	OriginIndex      int    `json:"originIndex"`
	DestinationIndex int    `json:"destinationIndex"`
	DistanceMeters   int    `json:"distanceMeters"`
	Duration         string `json:"duration"`
}

func toMatrixOrigin(loc *domain.Location) (matrixOrigin, error) {
	origin := matrixOrigin{
		Waypoint: matrixWaypoint{
			Location: matrixLocation{
				LatLng: matrixLatLng{
					Latitude: loc.Latitude, Longitude: loc.Longitude,
				},
			},
		},
	}
	return origin, nil
}

func toMatrixDestination(loc *domain.Location) (matrixDestination, error) {
	destination := matrixDestination{
		Waypoint: matrixWaypoint{
			Location: matrixLocation{
				LatLng: matrixLatLng{
					Latitude: loc.Latitude, Longitude: loc.Longitude,
				},
			},
		},
	}
	return destination, nil
}

func (e *RouteEngine) computeRouteMatrix(
	ctx context.Context,
	routes []*domain.Route,
) ([]domain.RouteMeasurement, error) {
	log.Println("REQUESTING MATRIX")
	if len(routes) == 0 {
		return nil, nil
	}

	// Build origin/destination list
	var origins []matrixOrigin
	var destinations []matrixDestination

	for _, r := range routes {
		o, err := toMatrixOrigin(r.Origin)
		if err != nil {
			return nil, err
		}

		d, err := toMatrixDestination(r.Destination)
		if err != nil {
			return nil, err
		}

		origins = append(origins, o)
		destinations = append(destinations, d)
	}

	reqBody := matrixRequest{
		Origins:           origins,
		Destinations:      destinations,
		TravelMode:        TRAVEL_MODE,
		RoutingPreference: ROUTING_PREFERENCE,
		// DepartureTime:     time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", "  "); err != nil {
		return nil, err
	}

	log.Printf("Matrix request:\n%s\n", prettyJSON.String())

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		GOOGLE_ENDPOINT,
		bytes.NewReader(jsonData),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", e.apiKey)
	req.Header.Set("X-Goog-FieldMask", "originIndex,destinationIndex,duration,distanceMeters")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Matrix error response:\n%s\n", string(bodyBytes))
		return nil, fmt.Errorf("matrix request failed: %s", resp.Status)
	}

	var elements []RouteMatrixElement
	if err := json.Unmarshal(bodyBytes, &elements); err != nil {
		return nil, err
	}

	now := time.Now()
	log.Printf("ELEMENTS %d:\n%v\n", len(elements), elements)

	// Map response to in-memory routes
	routeMeasurements, err := mapMatrixElements(
		routes,
		elements,
		now,
	)
	if err != nil {
		return nil, err
	}

	// Debugging lines
	// log.Printf("ROUTES:\n")
	// for _, route := range routes {
	// 	log.Printf("%s %s", route.Origin.Name, route.Destination.Name)
	// }
	// log.Printf("ROUTE MEASURES %d", len(routeMeasurements))
	// for _, route := range routeMeasurements {
	// 	log.Printf(
	// 		"Route measure: %d %d %s %s\n",
	// 		route.RouteID,
	// 		route.DistanceMeters,
	// 		route.DurationSeconds,
	// 		route.RecordedAt,
	// 	)
	// }

	// Persist measurement
	if err := e.Store.UpdateMeasurements(ctx, routeMeasurements); err != nil {
		return nil, err
	}

	return routeMeasurements, nil
}
