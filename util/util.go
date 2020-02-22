package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	body := map[string]string{
		"error": message,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Fatal(err)
	}
}
