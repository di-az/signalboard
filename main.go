package main

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/engine"
	"commuteboard/internal/store"
	"fmt"
	"time"
)

const UpdateRate = 1 * time.Minute
const tickRate = 10 * time.Second

var home = domain.Location{
	Name:      "Home",
	Latitude:  "20.745317326696103",
	Longitude: "-103.44431208289149",
	// Schedule:  Schedule{Times: []string{"08:00-10:00"}},
}

var work = domain.Location{
	Name:      "Work",
	Latitude:  "20.688900217575455",
	Longitude: "-103.42880959994349",
	Schedule: domain.Schedule{
		Days: map[time.Weekday][]domain.TimeRange{
			time.Tuesday: {
				{Start: 8 * time.Hour, End: 10 * time.Hour},
			},
			time.Thursday: {
				{Start: 8 * time.Hour, End: 10 * time.Hour},
			},
			time.Saturday: {
				{Start: 1 * time.Hour, End: 23 * time.Hour},
			},
		},
	},
}

var piano = domain.Location{
	Name:      "Piano",
	Latitude:  "20.688900217575455",
	Longitude: "-103.42880959994349",
	Schedule: domain.Schedule{
		Days: map[time.Weekday][]domain.TimeRange{
			time.Saturday: {
				{Start: 9 * time.Hour, End: 18 * time.Hour},
			},
		},
	},
}

func main() {
	fmt.Println("Running server")

	locations := []*domain.Location{&work, &piano}

	store := store.NewRouteStore()

	engine := engine.NewRouteEngine(home, locations, store, UpdateRate, tickRate)

	go engine.Run()

	select {}
}
