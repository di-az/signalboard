package commute

import (
	"testing"
	"time"
)

func TestMapMatrixElements(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local)

	routes := []*Route{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	tests := []struct {
		name          string
		elements      []RouteMatrixElement
		expectedCount int
		expectError   bool
	}{
		{
			name: "maps diagonal routes only",
			elements: []RouteMatrixElement{
				{
					OriginIndex:      0,
					DestinationIndex: 0,
					DistanceMeters:   5000,
					Duration:         "600s",
				},
				{
					OriginIndex:      0,
					DestinationIndex: 1,
					DistanceMeters:   9999,
					Duration:         "999s",
				},
				{
					OriginIndex:      1,
					DestinationIndex: 0,
					DistanceMeters:   9999,
					Duration:         "999s",
				},
				{
					OriginIndex:      1,
					DestinationIndex: 1,
					DistanceMeters:   10000,
					Duration:         "1200s",
				},
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "invalid duration format",
			elements: []RouteMatrixElement{
				{
					OriginIndex:      0,
					DestinationIndex: 0,
					DistanceMeters:   5000,
					Duration:         "invalid",
				},
			},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name: "out of bounds index ignored",
			elements: []RouteMatrixElement{
				{
					OriginIndex:      99,
					DestinationIndex: 99,
					DistanceMeters:   5000,
					Duration:         "600s",
				},
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "cross routes ignored",
			elements: []RouteMatrixElement{
				{
					OriginIndex:      0,
					DestinationIndex: 1,
					DistanceMeters:   5000,
					Duration:         "600s",
				},
				{
					OriginIndex:      1,
					DestinationIndex: 0,
					DistanceMeters:   5000,
					Duration:         "600s",
				},
			},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			measurements, err := mapMatrixElements(
				routes,
				tt.elements,
				now,
			)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(measurements) != tt.expectedCount {
				t.Fatalf(
					"expected %d measurements got %d",
					tt.expectedCount,
					len(measurements),
				)
			}
		})
	}
}

func TestMapMatrixElementsMapsCorrectRoutes(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local)

	routes := []*Route{
		{
			ID: 100,
		},
		{
			ID: 200,
		},
	}

	elements := []RouteMatrixElement{
		{
			OriginIndex:      0,
			DestinationIndex: 0,
			DistanceMeters:   5000,
			Duration:         "600s",
		},
		{
			OriginIndex:      1,
			DestinationIndex: 1,
			DistanceMeters:   10000,
			Duration:         "1200s",
		},
	}

	measurements, err := mapMatrixElements(routes, elements, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(measurements) != 2 {
		t.Fatalf("expected 2 measurements got %d", len(measurements))
	}

	if measurements[0].RouteID != 100 {
		t.Fatalf("expected route ID 100 got %d", measurements[0].RouteID)
	}

	if measurements[1].RouteID != 200 {
		t.Fatalf("expected route ID 200 got %d", measurements[1].RouteID)
	}

	if measurements[0].DistanceMeters != 5000 {
		t.Fatalf(
			"expected distance 5000 got %d",
			measurements[0].DistanceMeters,
		)
	}

	if measurements[1].DistanceMeters != 10000 {
		t.Fatalf(
			"expected distance 10000 got %d",
			measurements[1].DistanceMeters,
		)
	}
}
