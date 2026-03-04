package main

import (
	"commuteboard/internal/db"
	"commuteboard/internal/engine"
	"commuteboard/internal/store"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const UpdateRate = 1 * time.Minute
const tickRate = 10 * time.Minute

// var home = domain.Location{
// 	ID:        "1",
// 	Name:      "Home",
// 	Latitude:  20.745317326696103,
// 	Longitude: -103.44431208289149,
// 	// Schedule:  Schedule{Times: []string{"08:00-10:00"}},
// }
//
// var work = domain.Location{
// 	ID:        "2",
// 	Name:      "Work",
// 	Latitude:  20.688900217575455,
// 	Longitude: -103.42880959994349,
// 	Schedule: domain.Schedule{
// 		Days: map[time.Weekday][]domain.TimeRange{
// 			time.Tuesday: {
// 				{Start: 8 * time.Hour, End: 10 * time.Hour},
// 			},
// 			time.Thursday: {
// 				{Start: 8 * time.Hour, End: 10 * time.Hour},
// 			},
// 			time.Monday: {
// 				{Start: 0 * time.Hour, End: 23 * time.Hour},
// 			},
// 		},
// 	},
// }
//
// var piano = domain.Location{
// 	ID:        "3",
// 	Name:      "Piano",
// 	Latitude:  20.688900217575455,
// 	Longitude: -103.42880959994349,
// 	Schedule: domain.Schedule{
// 		Days: map[time.Weekday][]domain.TimeRange{
// 			time.Saturday: {
// 				{Start: 9 * time.Hour, End: 18 * time.Hour},
// 			},
// 		},
// 	},
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_MAPS_API_KEY not set")
	}

	sqliteDB, err := db.NewSQLite("routes.db")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting engine")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.NewRouteStore(sqliteDB)

	engine, err := engine.NewRouteEngine(ctx, store, UpdateRate, tickRate, apiKey)
	if err != nil {
		log.Fatal(err)
	}
	// server := server.NewHttpServer(store, engine)
	//
	// // Run HTTP server
	// go func() {
	// 	if err := server.Run(ctx); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("server error: %v", err)
	// 	}
	// }()

	// Run engine
	go engine.Run(ctx)

	<-ctx.Done()
	log.Println("Shutting down application...")
}
