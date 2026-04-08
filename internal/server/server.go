package server

import (
	"commuteboard/internal/engine"
	"commuteboard/internal/store"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
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
	mux.HandleFunc("/routes/active", s.GetActiveRoutes)
	mux.HandleFunc("/routes/refresh", s.RefreshRoutes)
	mux.HandleFunc("/health", s.CheckHealth)
	mux.HandleFunc("/engine/status", s.EngineStatus)

	server := &http.Server{
		Addr:    PORT,
		Handler: withCORS(mux),
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
	log.Println("GET /routes")

	routes := s.engine.GetRoutes()
	response := make([]RouteResponse, 0, len(routes))

	for _, route := range routes {
		response = append(response, NewRouteResponse(route))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *HttpServer) GetActiveRoutes(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /routes/active")

	now := time.Now()
	active := make([]RouteResponse, 0)

	routes := s.engine.GetRoutes()
	for _, route := range routes {
		if route.Schedule.ShouldRunNow(now) {
			active = append(active, NewRouteResponse(route))
		}
	}

	SortRouteResponseSlice(active)
	writeJSON(w, http.StatusOK, active)
}

func (s *HttpServer) RefreshRoutes(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /routes/refresh")

	err := s.engine.RefreshNow(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("failed to refresh routes"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "refresh triggered",
	})
}

// func (s *HttpServer) GetRouteByID(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	log.Printf("Getting route for id: %v\n", id)
//
// 	route, err := s.store.GetByID(id)
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}
//
// 	origin := s.engine.GetLocationByID(route.OriginID)
// 	destination := s.engine.GetLocationByID(route.DestinationID)
// 	if origin == nil || destination == nil {
// 		writeError(w, errors.New("missing origin or destination"))
// 		return
// 	}
//
// 	routeResp := NewCommuteResponse(*origin, *destination, route)
// 	writeJSON(w, http.StatusOK, routeResp)
// }

func (s *HttpServer) CheckHealth(w http.ResponseWriter, r *http.Request) {
	healthMessage := map[string]string{"status": "ok"}
	writeJSON(w, http.StatusOK, healthMessage)
}

func (s *HttpServer) EngineStatus(w http.ResponseWriter, r *http.Request) {
	status := s.engine.Status()
	writeJSON(w, http.StatusOK, status)
}
