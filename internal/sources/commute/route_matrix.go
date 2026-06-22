package commute

import (
	"bytes"
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

func toMatrixOrigin(loc *Location) (matrixOrigin, error) {
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

func toMatrixDestination(loc *Location) (matrixDestination, error) {
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

// Computing Route Matrix using Google Maps Route Matrix.
// https://developers.google.com/maps/documentation/routes/compute-route-matrix-over
func (s *CommuteSource) computeRouteMatrix(
	ctx context.Context,
	routes []*Route,
) ([]RouteMeasurement, error) {
	// log.Println("Helper: NOT REQUESTING MATRIX")
	// return nil, nil
	log.Println("Updating routes!")
	if len(routes) == 0 {
		return nil, nil
	}

	// Build origin/destination list
	var destinations []matrixDestination
	// var origins []matrixOrigin

	// Pick first route as single origin for Google Request pricing optimization
	origin, err := toMatrixOrigin(routes[0].Origin)
	if err != nil {
		return nil, err
	}

	for _, r := range routes {
		d, err := toMatrixDestination(r.Destination)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, d)
	}

	reqBody := matrixRequest{
		Origins: []matrixOrigin{
			origin,
		},
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

	// log.Printf("Matrix request:\n%s\n", prettyJSON.String())

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
	req.Header.Set("X-Goog-Api-Key", s.apiKey)
	req.Header.Set("X-Goog-FieldMask", "originIndex,destinationIndex,duration,distanceMeters")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// DEBUG: Debugging lines
	// log.Printf("DEBUG:\n")
	// log.Print(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Matrix error response:\n%s\n", string(bodyBytes))
		return nil, fmt.Errorf("matrix request failed: %s", resp.Status)
	}

	var elements []RouteMatrixElement
	if err := json.Unmarshal(bodyBytes, &elements); err != nil {
		return nil, err
	}

	now := time.Now()

	// Map response to in-memory routes
	routeMeasurements, err := mapMatrixElements(
		routes,
		elements,
		now,
	)
	if err != nil {
		return nil, err
	}

	// DEBUG: Debugging lines
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
	if err := s.Store.UpdateMeasurements(ctx, routeMeasurements); err != nil {
		return nil, err
	}

	return routeMeasurements, nil
}
