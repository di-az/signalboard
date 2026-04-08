package server

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// if data == nil {
	// 	return errors.New("empty response body")
	// }

	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	// log.Println(err.Error())
	_ = writeJSON(w, status, err.Error())
}
