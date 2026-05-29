package main

import (
	"log"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
	authpostgres "github.com/krishnaditya65/auth-server/internal/auth/infra/postgres"
	authhttp "github.com/krishnaditya65/auth-server/internal/auth/transport/http"

	identitypostgres "github.com/krishnaditya65/auth-server/internal/identity/infra/postgres"

	"github.com/krishnaditya65/auth-server/internal/platform/config"
	"github.com/krishnaditya65/auth-server/internal/platform/httpserver"
	"github.com/krishnaditya65/auth-server/internal/platform/nats"
	platformpostgres "github.com/krishnaditya65/auth-server/internal/platform/postgres"
	pgtx "github.com/krishnaditya65/auth-server/internal/platform/postgres/tx"
	"github.com/krishnaditya65/auth-server/internal/platform/redis"
)

func main() {
	cfg := config.Load()

	db, err := platformpostgres.New(cfg.DatabaseURL)
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

	// transaction manager
	txManager := pgtx.NewManager(db)

	// repositories
	identityRepo := identitypostgres.NewRepository(db)
	credentialRepo := authpostgres.NewRepository(db)

	// usecases
	registerUseCase := authapp.NewRegisterUseCase(
		txManager,
		identityRepo,
		credentialRepo,
	)

	// handlers
	authHandler := authhttp.NewHandler(
		registerUseCase,
	)

	server := httpserver.New(cfg.HTTPPort)

	server.Handle(
		"POST",
		"/register",
		authHandler.Register,
	)

	log.Printf("starting auth server on port %s", cfg.HTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
