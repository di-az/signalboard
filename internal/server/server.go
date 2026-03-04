package server

import (
	"commuteboard/internal/engine"
	"commuteboard/internal/store"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
)

const PORT = ":3333"

type HttpServer struct {
	store  *store.RouteStore
	engine *engine.RouteEngine
}

func NewHttpServer(store *store.RouteStore, engine *engine.RouteEngine) *HttpServer {
	return &HttpServer{store: store, engine: engine}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/routes", s.GetRoutes)
	mux.HandleFunc("/routes/{id}", s.GetRouteByID)
	mux.HandleFunc("/health", s.CheckHealth)
	mux.HandleFunc("/engine/status", s.EngineStatus)

	server := &http.Server{
		Addr:    PORT,
		Handler: mux,
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

func (s *HttpServer) GetRoutes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting routes\n")
	var responses []CommuteResponse
	for _, route := range s.store.GetAll() {
		origin := s.engine.GetLocationByID(route.OriginID)
		destination := s.engine.GetLocationByID(route.DestinationID)
		if origin == nil || destination == nil {
			continue
		}

		newComm := NewCommuteResponse(*origin, *destination, route)
		responses = append(responses, newComm)
	}
	writeJSON(w, http.StatusOK, responses)
}

func (s *HttpServer) GetRouteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("Getting route for id: %v\n", id)

	route, err := s.store.GetByID(id)
	if err != nil {
		writeError(w, err)
		return
	}

	origin := s.engine.GetLocationByID(route.OriginID)
	destination := s.engine.GetLocationByID(route.DestinationID)
	if origin == nil || destination == nil {
		writeError(w, errors.New("missing origin or destination"))
		return
	}

	routeResp := NewCommuteResponse(*origin, *destination, route)
	writeJSON(w, http.StatusOK, routeResp)
}

func (s *HttpServer) CheckHealth(w http.ResponseWriter, r *http.Request) {
	healthMessage := map[string]string{"status": "ok"}
	writeJSON(w, http.StatusOK, healthMessage)
}

func (s *HttpServer) EngineStatus(w http.ResponseWriter, r *http.Request) {
	status := s.engine.Status()
	writeJSON(w, http.StatusOK, status)
}
