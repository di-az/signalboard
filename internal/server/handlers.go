package server

import "net/http"

func (s *HttpServer) CheckHealth(w http.ResponseWriter, r *http.Request) {
	healthMessage := map[string]string{"status": "ok"}
	WriteJSON(w, http.StatusOK, healthMessage)
}

func (s *HttpServer) EngineStatus(w http.ResponseWriter, r *http.Request) {
	status := s.engine.Status()
	WriteJSON(w, http.StatusOK, status)
}
