package main

import (
	"log"

	"github.com/nabilfikrisp/go-crud/config"
	"github.com/nabilfikrisp/go-crud/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
