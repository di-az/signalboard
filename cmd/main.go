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
	"syscall"
)

func main() {
	appCfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	sqliteDB, err := db.NewSQLite(appCfg.SQLiteDB)
	if err != nil {
		log.Fatal(err)
	}

	commuteCfg, err := commute.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting engine")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := commute.NewCommuteStore(sqliteDB)

	commuteSource, err := commute.NewCommuteSource(
		ctx,
		store,
		commuteCfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	engine := engine.NewEngine(
		appCfg.TickRate,
		commuteSource,
	)
	server := server.NewHttpServer(engine)

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
