package main

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/engine"
	"commuteboard/internal/server"
	"commuteboard/internal/store"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const UpdateRate = 1 * time.Minute
const tickRate = 10 * time.Second

var home = domain.Location{
	ID:        "1",
	Name:      "Home",
	Latitude:  "20.745317326696103",
	Longitude: "-103.44431208289149",
	// Schedule:  Schedule{Times: []string{"08:00-10:00"}},
}

var work = domain.Location{
	ID:        "2",
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
			time.Monday: {
				{Start: 1 * time.Hour, End: 23 * time.Hour},
			},
		},
	},
}

var piano = domain.Location{
	ID:        "3",
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	locations := []*domain.Location{&work, &piano}
	store := store.NewRouteStore()
	engine := engine.NewRouteEngine(home, locations, store, UpdateRate, tickRate)
	server := server.NewHttpServer(store, engine)

	// Run HTTP server
	go func() {
		if err := server.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Run engine
	go engine.Run(ctx)

	<-ctx.Done()
	log.Println("Shutting down application...")
}
