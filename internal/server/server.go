package server

import (
	"commuteboard/internal/store"
	"context"
	"log"
	"net"
	"net/http"
)

const PORT = ":3333"

type HttpServer struct {
	store *store.RouteStore
}

func NewHttpServer(store *store.RouteStore) *HttpServer {
	return &HttpServer{store: store}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/routes", s.GetRoutes)
	mux.HandleFunc("/routes/{id}", s.GetRouteByID)
	mux.HandleFunc("/health", s.CheckHealth)

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
	var responseRoutes []RouteResponse
	for _, route := range s.store.GetAll() {
		response := NewRouteResponse(route)
		responseRoutes = append(responseRoutes, response)
	}
	writeJSON(w, http.StatusOK, responseRoutes)
}

func (s *HttpServer) GetRouteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("Getting route for id: %v\n", id)
	route, err := s.store.GetByID(id)
	routeResp := NewRouteResponse(route)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, routeResp)
}

func (s *HttpServer) CheckHealth(w http.ResponseWriter, r *http.Request) {
	healthMessage := map[string]string{"status": "ok"}
	writeJSON(w, http.StatusOK, healthMessage)
}
