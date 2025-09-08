package main

import (
	"fmt"
	"net/http"

	"open-meet/pkg/api"
	"open-meet/pkg/config"
)

func main() {
	// Load configuration from .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Initialize API service with config
	service, err := api.NewEngine(cfg)
	if err != nil {
		fmt.Printf("Failed to create service: %v\n", err)
		return
	}

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, service); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
