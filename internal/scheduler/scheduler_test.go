package scheduler

import (
	"testing"
	"time"
)

func TestSchedulerFirstRun(t *testing.T) {
	start := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)

	scheduler := NewScheduler([]Schedule{
		{
			Plugin:   "weather",
			Widget:   "activity",
			Start:    start,
			End:      start.Add(time.Hour),
			Interval: 30 * time.Minute,
		},
	})

	jobs := scheduler.Due(start)

	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
}

func TestSchedulerInterval(t *testing.T) {
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		lastRun time.Time
		want    int
	}{
		{
			name:    "before interval",
			lastRun: now.Add(-15 * time.Minute),
			want:    0,
		},
		{
			name:    "at interval",
			lastRun: now.Add(-30 * time.Minute),
			want:    1,
		},
		{
			name:    "after interval",
			lastRun: now.Add(-45 * time.Minute),
			want:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := NewScheduler([]Schedule{
				{
					Plugin:   "github",
					Widget:   "activity",
					Start:    now.Add(-time.Hour),
					End:      now.Add(time.Hour),
					Interval: 30 * time.Minute,
					LastRun:  tt.lastRun,
				},
			})

			jobs := scheduler.Due(now)

			if len(jobs) != tt.want {
				t.Fatalf("expected %d jobs, got %d", tt.want, len(jobs))
			}
		})
	}
}

func TestSchedulerMultipleSchedules(t *testing.T) {
	start := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)

	scheduler := NewScheduler([]Schedule{
		{
			Plugin:   "weather",
			Widget:   "activity",
			Start:    start,
			End:      start.Add(3 * time.Hour),
			Interval: 30 * time.Minute,
		},
		{
			Plugin:   "clock",
			Widget:   "current",
			Start:    start.Add(time.Hour),
			End:      start.Add(4 * time.Hour),
			Interval: time.Hour,
		},
	})

	jobs := scheduler.Due(start.Add(2 * time.Hour))

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].Plugin != scheduler.Schedules[0].Plugin {
		t.Errorf("expected %s job, got %q", scheduler.Schedules[0].Plugin, jobs[0].Plugin)
	}

	if jobs[1].Plugin != scheduler.Schedules[1].Plugin {
		t.Errorf("expected %s job, got %q", scheduler.Schedules[1].Plugin, jobs[1].Plugin)
	}
}

func TestSchedulerInvalidInterval(t *testing.T) {
	start := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)

	scheduler := NewScheduler([]Schedule{
		{
			Plugin:   "weather",
			Widget:   "activity",
			Start:    start,
			End:      start.Add(time.Hour),
			Interval: 0,
		},
	})

	jobs := scheduler.Due(start)

	if len(jobs) != 0 {
		t.Fatalf("expected no jobs, got %d", len(jobs))
	}
}

func TestSchedulerDoesNotReturnJobTwice(t *testing.T) {
	start := time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)

	scheduler := NewScheduler([]Schedule{
		{
			Plugin:   "weather",
			Widget:   "activity",
			Start:    start,
			End:      start.Add(time.Hour),
			Interval: 30 * time.Minute,
		},
	})

	first := scheduler.Due(start)

	if len(first) != 1 {
		t.Fatalf("expected 1 job, got %d", len(first))
	}

	second := scheduler.Due(start)

	if len(second) != 0 {
		t.Fatalf("expected 0 jobs, got %d", len(second))
	}
}
