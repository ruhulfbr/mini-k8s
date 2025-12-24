package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{
		Status:  "success",
		Message: "Hello from Docker API container",
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Server running on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
