package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"signalboard/internal/engine"
)

const PORT = ":3333"

type HttpServer struct {
	engine *engine.Engine
}

func NewHttpServer(
	engine *engine.Engine,
) *HttpServer {
	return &HttpServer{engine: engine}
}

func (s *HttpServer) registerSourceRoutes(mux *http.ServeMux) {
	for _, source := range s.engine.Sources {
		base := "/" + source.Name()

		for _, endpoint := range source.Endpoints() {
			path := base + endpoint.Path

			mux.HandleFunc(path, endpoint.Handler)
		}
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.CheckHealth)
	mux.HandleFunc("/engine/status", s.EngineStatus)
	s.registerSourceRoutes(mux)

	server := &http.Server{
		Addr: PORT,
		Handler: withLogging(
			withCORS(mux),
		),
	}

	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		return err
	}
	log.Printf("HTTP server started\n")

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		_ = server.Shutdown(context.Background())
	}()

	err = server.Serve(ln)

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
