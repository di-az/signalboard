package commute

import (
	"testing"
	"time"
)

func TestScheduleShouldRunNow(t *testing.T) {
	tests := []struct {
		name     string
		schedule Schedule
		now      time.Time
		expected bool
	}{
		{
			name: "within single time range",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local), // Monday 12:00
			expected: true,
		},
		{
			name: "before time range",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 7, 59, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "after time range",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 18, 0, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "wrong weekday",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 21, 12, 0, 0, 0, time.Local), // Tuesday
			expected: false,
		},
		{
			name: "exactly at start boundary",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 8, 0, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "exactly at end boundary",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   17 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 17, 0, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "multiple time ranges second range active",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   10 * time.Hour,
						},
						{
							Start: 14 * time.Hour,
							End:   18 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 15, 0, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "multiple time ranges no active range",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{
					time.Monday: {
						{
							Start: 8 * time.Hour,
							End:   10 * time.Hour,
						},
						{
							Start: 14 * time.Hour,
							End:   18 * time.Hour,
						},
					},
				},
			},
			now:      time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "empty schedule",
			schedule: Schedule{
				Days: map[time.Weekday][]TimeRange{},
			},
			now:      time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schedule.ShouldRunNow(tt.now)

			if result != tt.expected {
				t.Fatalf("expected %v got %v", tt.expected, result)
			}
		})
	}
}

func TestRouteIsFresh(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.Local)

	tests := []struct {
		name       string
		route      Route
		updateRate time.Duration
		expected   bool
	}{
		{
			name: "fresh route",
			route: Route{
				RecordedAt: timePtr(now.Add(-5 * time.Minute)),
			},
			updateRate: 10 * time.Minute,
			expected:   true,
		},
		{
			name: "stale route",
			route: Route{
				RecordedAt: timePtr(now.Add(-15 * time.Minute)),
			},
			updateRate: 10 * time.Minute,
			expected:   false,
		},
		{
			name: "exactly at freshness boundary",
			route: Route{
				RecordedAt: timePtr(now.Add(-10 * time.Minute)),
			},
			updateRate: 10 * time.Minute,
			expected:   true,
		},
		{
			name: "nil recorded at",
			route: Route{
				RecordedAt: nil,
			},
			updateRate: 10 * time.Minute,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.route.IsFresh(now, tt.updateRate)

			if result != tt.expected {
				t.Fatalf("expected %v got %v", tt.expected, result)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
