package server

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// if data == nil {
	// 	return errors.New("empty response body")
	// }

	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	// log.Println(err.Error())
	_ = WriteJSON(w, status, err.Error())
}
