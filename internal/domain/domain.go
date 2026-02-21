package domain

import "time"

type TimeRange struct {
	Start time.Duration
	End   time.Duration
}

type Schedule struct {
	Days        map[time.Weekday][]TimeRange
	LastUpdated time.Time
}

type Location struct {
	ID        string
	Name      string
	Latitude  string
	Longitude string
	Schedule  Schedule
}

type Route struct {
	Start        Location
	Finish       Location
	Minutes      int
	TrafficLevel string
	Timestamp    time.Time
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
