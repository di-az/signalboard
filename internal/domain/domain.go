package domain

import (
	"time"
)

type TimeRange struct {
	Start time.Duration
	End   time.Duration
}

type Schedule struct {
	Days map[time.Weekday][]TimeRange
}

type Location struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
}

type Route struct {
	ID              int
	Origin          *Location
	Destination     *Location
	DistanceMeters  int
	DurationSeconds time.Duration
	RecordedAt      time.Time
	Schedule        Schedule
}

type RouteMeasurement struct {
	RouteID         int
	DistanceMeters  int
	DurationSeconds time.Duration
	RecordedAt      time.Time
}

func (s Schedule) ShouldRunNow(t time.Time) bool {
	dayRanges, ok := s.Days[t.Weekday()]
	if !ok {
		return false
	}

	nowMinutes := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute

	for _, timeRange := range dayRanges {
		if nowMinutes >= timeRange.Start && nowMinutes <= timeRange.End {
			// log.Printf("Need to run now for: %v %v-%v", t.Weekday(), timeRange.Start, timeRange.End)
			return true
		}
	}

	return false
}

func (c Route) IsFresh(now time.Time, updateRate time.Duration) bool {
	timeToUpdate := c.RecordedAt.Add(updateRate)
	if timeToUpdate.Before(now) {
		return false
	}
	return true
}
