package commute

import (
	"testing"
)

func TestSortRouteResponseSliceAscending(t *testing.T) {
	routes := []RouteResponse{
		{
			ID:              1,
			DurationMinutes: intPtr(30),
		},
		{
			ID:              2,
			DurationMinutes: intPtr(10),
		},
		{
			ID:              3,
			DurationMinutes: intPtr(20),
		},
	}

	SortRouteResponseSlice(routes)

	if *routes[0].DurationMinutes != 10 {
		t.Fatalf("expected first duration to be 10")
	}

	if *routes[1].DurationMinutes != 20 {
		t.Fatalf("expected second duration to be 20")
	}

	if *routes[2].DurationMinutes != 30 {
		t.Fatalf("expected third duration to be 30")
	}
}

func TestSortRouteResponseSliceNilLast(t *testing.T) {
	routes := []RouteResponse{
		{
			ID:              1,
			DurationMinutes: nil,
		},
		{
			ID:              2,
			DurationMinutes: intPtr(10),
		},
		{
			ID:              3,
			DurationMinutes: intPtr(20),
		},
	}

	SortRouteResponseSlice(routes)

	if *routes[0].DurationMinutes != 10 {
		t.Fatalf("expected first duration to be 10")
	}

	if *routes[1].DurationMinutes != 20 {
		t.Fatalf("expected second duration to be 20")
	}

	if routes[2].DurationMinutes != nil {
		t.Fatalf("expected nil duration to be last")
	}
}

func TestSortRouteResponseSliceAllNil(t *testing.T) {
	routes := []RouteResponse{
		{
			ID:              1,
			DurationMinutes: nil,
		},
		{
			ID:              2,
			DurationMinutes: nil,
		},
	}

	SortRouteResponseSlice(routes)

	if routes[0].DurationMinutes != nil {
		t.Fatalf("expected first duration to remain nil")
	}

	if routes[1].DurationMinutes != nil {
		t.Fatalf("expected second duration to remain nil")
	}
}

func TestSortRouteResponseSliceEmpty(t *testing.T) {
	var routes []RouteResponse

	SortRouteResponseSlice(routes)

	if len(routes) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestSortRouteResponseSliceSingle(t *testing.T) {
	routes := []RouteResponse{
		{
			ID:              1,
			DurationMinutes: intPtr(15),
		},
	}

	SortRouteResponseSlice(routes)

	if len(routes) != 1 {
		t.Fatalf("expected 1 route")
	}

	if *routes[0].DurationMinutes != 15 {
		t.Fatalf("expected duration unchanged")
	}
}

func intPtr(v int) *int {
	return &v
}
