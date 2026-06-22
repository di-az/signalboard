package commute

import (
	"errors"
	"net/http"
	"signalboard/internal/server"
	"time"
)

func (s *CommuteSource) GetRoutesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := s.GetRoutes()
		response := make([]RouteResponse, 0, len(routes))

		for _, route := range routes {
			response = append(response, NewRouteResponse(route))
		}
		server.WriteJSON(w, http.StatusOK, response)
	}
}

func (s *CommuteSource) GetActiveRoutesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		active := make([]RouteResponse, 0)

		routes := s.GetRoutes()
		for _, route := range routes {
			if route.Schedule.ShouldRunNow(now) {
				active = append(active, NewRouteResponse(route))
			}
		}

		SortRouteResponseSlice(active)
		server.WriteJSON(w, http.StatusOK, active)
	}
}

func (s *CommuteSource) RefreshRoutes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Refresh(r.Context())
		if err != nil {
			server.WriteError(
				w,
				http.StatusInternalServerError,
				errors.New("failed to refresh routes"),
			)
			return
		}

		server.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "refresh triggered",
		})
	}
}
