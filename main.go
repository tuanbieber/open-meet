package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})
	_ = http.ListenAndServe(":8080", nil)
}
