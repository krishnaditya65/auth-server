package main

import (
	"log"

	"github.com/krishnaditya65/auth-server/internal/platform/config"
	"github.com/krishnaditya65/auth-server/internal/platform/httpserver"
)

func main() {
	cfg := config.Load()

	server := httpserver.New(cfg.HTTPPort)

	log.Printf("starting auth server on port %s", cfg.HTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
