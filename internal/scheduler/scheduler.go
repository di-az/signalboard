package scheduler

import (
	"signalboard/internal/plugin"
	"time"
)

// Scheduler determines which scheduled widgets are due for execution.
type Scheduler struct {
	Schedules []Schedule
}

// Schedule defines a widget execution parameters.
type Schedule struct {
	Plugin   string
	Widget   string
	Start    time.Time
	End      time.Time
	Interval time.Duration
	Config   plugin.Request
	LastRun  time.Time
}

// Job represents a widget execution that is due.
type Job struct {
	Plugin string
	Widget string
	Config plugin.Request
}

func NewScheduler(schedules []Schedule) *Scheduler {
	return &Scheduler{
		Schedules: schedules,
	}
}

func (s *Scheduler) Due(now time.Time) []Job {
	var jobs []Job

	for i := range s.Schedules {
		schedule := &s.Schedules[i]

		if now.Before(schedule.Start) || now.After(schedule.End) {
			continue
		}

		if !isDue(*schedule, now) {
			continue
		}

		jobs = append(jobs, Job{
			Plugin: schedule.Plugin,
			Widget: schedule.Widget,
			Config: schedule.Config,
		})

		schedule.LastRun = now
	}

	return jobs
}

func isDue(schedule Schedule, now time.Time) bool {
	if schedule.Interval <= 0 {
		return false
	}

	if schedule.LastRun.IsZero() {
		return true
	}

	return now.Sub(schedule.LastRun) >= schedule.Interval
}
