package main

import (
	"log"

	"github.com/krishnaditya65/auth-server/internal/platform/config"
	"github.com/krishnaditya65/auth-server/internal/platform/httpserver"
	"github.com/krishnaditya65/auth-server/internal/platform/nats"
	"github.com/krishnaditya65/auth-server/internal/platform/postgres"
	"github.com/krishnaditya65/auth-server/internal/platform/redis"
)

func main() {
	cfg := config.Load()

	db, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres init failed: %v", err)
	}
	defer db.Close()

	cache, err := redis.New(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis init failed: %v", err)
	}
	defer cache.Close()

	msgBus, err := nats.New(cfg.NATSURL)
	if err != nil {
		log.Fatalf("nats init failed: %v", err)
	}
	defer msgBus.Close()

	log.Println("postgres connected")
	log.Println("redis connected")
	log.Println("nats connected")

	server := httpserver.New(cfg.HTTPPort)

	log.Printf("starting auth server on port %s", cfg.HTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
