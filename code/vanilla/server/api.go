package server

import (
	"encoding/json"
	"net/http"
)

func ApiRoutes() {
	http.HandleFunc("/", handleHello)
}

type HelloResponse struct {
	Message string `json:"message"`
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{Message: "Hello, World!"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
