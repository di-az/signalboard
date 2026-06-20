package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"signalboard/internal/config"
	"signalboard/internal/db"
	"signalboard/internal/engine"
	"signalboard/internal/server"
	"signalboard/internal/sources/commute"
	"signalboard/internal/store"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	sqliteDB, err := db.NewSQLite(cfg.SQLiteDB)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting engine")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.NewRouteStore(sqliteDB)

	commuteSource, err := commute.NewCommuteSource(
		ctx,
		store,
		cfg.UpdateRate,
		cfg.TickRate,
		cfg.GoogleMapsAPIKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	engine := engine.NewEngine(
		cfg.TickRate,
		commuteSource,
	)
	server := server.NewHttpServer(store, engine, commuteSource)

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
