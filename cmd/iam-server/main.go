package main

import (
	"log"
	"net/http"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
	authpostgres "github.com/krishnaditya65/auth-server/internal/auth/infra/postgres"
	authhttp "github.com/krishnaditya65/auth-server/internal/auth/transport/http"
	authorizationapp "github.com/krishnaditya65/auth-server/internal/authorization/app"

	authorizationpostgres "github.com/krishnaditya65/auth-server/internal/authorization/infra/postgres"

	identitypostgres "github.com/krishnaditya65/auth-server/internal/identity/infra/postgres"
	tenantapp "github.com/krishnaditya65/auth-server/internal/tenant/app"
	tenantpostgres "github.com/krishnaditya65/auth-server/internal/tenant/infra/postgres"

	postgresuser "github.com/krishnaditya65/auth-server/internal/identity/infra/postgresuser"
	sessionpostgres "github.com/krishnaditya65/auth-server/internal/session/infra/postgres"

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
	tenantRepo := tenantpostgres.NewRepository(db)

	userRepo := postgresuser.NewRepository(db)
	sessionRepo := sessionpostgres.NewRepository(
		db,
	)

	// usecases
	slugService := tenantapp.NewSlugService(
		tenantRepo,
	)

	roleRepo := authorizationpostgres.NewRoleRepository(db)

	userRoleRepo := authorizationpostgres.NewUserRoleRepository(db)

	bootstrapService := authorizationapp.NewBootstrapService(
		roleRepo,
	)

	registerUseCase := authapp.NewRegisterUseCase(
		txManager,
		identityRepo,
		credentialRepo,
		tenantRepo,
		userRepo,
		slugService,
		userRoleRepo,
		bootstrapService,
	)

	loginUseCase := authapp.NewLoginUseCase(
		identityRepo,
		credentialRepo,
		userRepo,
		userRoleRepo,
		sessionRepo,
	)

	meUseCase := authapp.NewMeUseCase(
		sessionRepo,
		identityRepo,
		userRepo,
		userRoleRepo,
	)

	// handlers
	authHandler := authhttp.NewHandler(
		registerUseCase,
		loginUseCase,
		meUseCase,
	)
	server := httpserver.New(cfg.HTTPPort)

	server.Handle(
		"GET",
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		},
	)

	server.Handle(
		"POST",
		"/register",
		authHandler.Register,
	)

	server.Handle(
		"POST",
		"/login",
		authHandler.Login,
	)

	server.Handle(
		"GET",
		"/me",
		authHandler.Me,
	)

	log.Printf("starting auth server on port %s", cfg.HTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
