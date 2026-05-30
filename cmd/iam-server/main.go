package main

import (
	"log/slog"
	"net/http"
	"os"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
	authpostgres "github.com/krishnaditya65/auth-server/internal/auth/infra/postgres"
	authmiddleware "github.com/krishnaditya65/auth-server/internal/auth/middleware"
	authhttp "github.com/krishnaditya65/auth-server/internal/auth/transport/http"
	authorizationapp "github.com/krishnaditya65/auth-server/internal/authorization/app"
	authorizationpostgres "github.com/krishnaditya65/auth-server/internal/authorization/infra/postgres"
	authorizationmiddleware "github.com/krishnaditya65/auth-server/internal/authorization/middleware"
	authorizationhttp "github.com/krishnaditya65/auth-server/internal/authorization/transport/http"

	identitypostgres "github.com/krishnaditya65/auth-server/internal/identity/infra/postgres"
	tenantapp "github.com/krishnaditya65/auth-server/internal/tenant/app"
	tenantpostgres "github.com/krishnaditya65/auth-server/internal/tenant/infra/postgres"

	postgresuser "github.com/krishnaditya65/auth-server/internal/identity/infra/postgresuser"
	sessionpostgres "github.com/krishnaditya65/auth-server/internal/session/infra/postgres"

	identityapp "github.com/krishnaditya65/auth-server/internal/identity/app"
	identityhttp "github.com/krishnaditya65/auth-server/internal/identity/transport/http"

	"github.com/krishnaditya65/auth-server/internal/platform/config"
	"github.com/krishnaditya65/auth-server/internal/platform/httpserver"
	"github.com/krishnaditya65/auth-server/internal/platform/nats"
	platformpostgres "github.com/krishnaditya65/auth-server/internal/platform/postgres"
	pgtx "github.com/krishnaditya65/auth-server/internal/platform/postgres/tx"
	"github.com/krishnaditya65/auth-server/internal/platform/redis"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := platformpostgres.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("postgres init failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	cache, err := redis.New(cfg.RedisAddr)
	if err != nil {
		slog.Error("redis init failed", "err", err)
		os.Exit(1)
	}
	defer cache.Close()

	msgBus, err := nats.New(cfg.NATSURL)
	if err != nil {
		slog.Error("nats init failed", "err", err)
		os.Exit(1)
	}
	defer msgBus.Close()

	slog.Info("infrastructure connected", "postgres", cfg.DatabaseURL != "", "env", cfg.AppEnv)

	txManager := pgtx.NewManager(db)

	// repositories
	identityRepo := identitypostgres.NewRepository(db)
	credentialRepo := authpostgres.NewRepository(db)
	tenantRepo := tenantpostgres.NewRepository(db)
	userRepo := postgresuser.NewRepository(db)
	sessionRepo := sessionpostgres.NewRepository(db)
	permissionRepo := authorizationpostgres.NewPermissionRepository(db)
	rolePermissionRepo := authorizationpostgres.NewRolePermissionRepository(db)
	roleRepo := authorizationpostgres.NewRoleRepository(db)
	userRoleRepo := authorizationpostgres.NewUserRoleRepository(db)

	// services
	slugService := tenantapp.NewSlugService(tenantRepo)
	bootstrapService := authorizationapp.NewBootstrapService(roleRepo)
	permissionBootstrapService := authorizationapp.NewPermissionBootstrapService(
		permissionRepo,
		rolePermissionRepo,
	)

	// use cases
	registerUseCase := authapp.NewRegisterUseCase(
		txManager,
		identityRepo,
		credentialRepo,
		tenantRepo,
		userRepo,
		slugService,
		userRoleRepo,
		bootstrapService,
		permissionBootstrapService,
	)

	loginUseCase := authapp.NewLoginUseCase(
		identityRepo,
		credentialRepo,
		userRepo,
		userRoleRepo,
		sessionRepo,
	)

	refreshUseCase := authapp.NewRefreshUseCase(sessionRepo)

	logoutUseCase := authapp.NewLogoutUseCase(sessionRepo)

	identityGetUserUseCase := identityapp.NewGetUserUseCase(userRepo)
	listUsersUseCase := identityapp.NewListUsersUseCase(userRepo)
	createUserUseCase := identityapp.NewCreateUserUseCase(
		txManager,
		identityRepo,
		credentialRepo,
		userRepo,
		roleRepo,
		userRoleRepo,
	)

	createRoleUseCase := authorizationapp.NewCreateRoleUseCase(roleRepo)
	listRolesUseCase := authorizationapp.NewListRolesUseCase(roleRepo)
	assignPermissionUseCase := authorizationapp.NewAssignPermissionToRoleUseCase(
		roleRepo,
		permissionRepo,
		rolePermissionRepo,
	)
	listRolePermissionsUseCase := authorizationapp.NewListRolePermissionsUseCase(
		roleRepo,
		rolePermissionRepo,
	)

	// handlers
	authHandler := authhttp.NewHandler(
		registerUseCase,
		loginUseCase,
		refreshUseCase,
		logoutUseCase,
	)

	identityHandler := identityhttp.NewHandler(
		identityGetUserUseCase,
		listUsersUseCase,
		createUserUseCase,
	)

	authorizationHandler := authorizationhttp.NewHandler(
		createRoleUseCase,
		listRolesUseCase,
		assignPermissionUseCase,
		listRolePermissionsUseCase,
	)

	// middleware
	authMiddleware := authmiddleware.NewAuthenticationMiddleware(
		sessionRepo,
		identityRepo,
		userRoleRepo,
		rolePermissionRepo,
	)

	server := httpserver.New(cfg.HTTPPort)
	r := server.Router()

	r.Use(httpserver.CORS)

	// public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/refresh", authHandler.Refresh)

	// authenticated routes
	withAuth := func(next http.Handler) http.Handler {
		return authMiddleware.Authenticate(next)
	}
	withPerm := func(p string) func(http.Handler) http.Handler {
		return authorizationmiddleware.RequirePermission(p)
	}

	r.With(withAuth).Get("/me", authHandler.Me)
	r.With(withAuth).Post("/logout", authHandler.Logout)

	r.With(withAuth, withPerm("users:read")).Get("/users", identityHandler.ListUsers)
	r.With(withAuth, withPerm("users:read")).Get("/users/{userID}", identityHandler.GetUser)
	r.With(withAuth, withPerm("users:create")).Post("/users", identityHandler.CreateUser)

	r.With(withAuth, withPerm("roles:create")).Post("/roles", authorizationHandler.CreateRole)
	r.With(withAuth, withPerm("roles:read")).Get("/roles", authorizationHandler.ListRoles)
	r.With(withAuth, withPerm("roles:update")).Post("/roles/{roleID}/permissions", authorizationHandler.AssignPermission)
	r.With(withAuth, withPerm("roles:read")).Get("/roles/{roleID}/permissions", authorizationHandler.ListRolePermissions)

	slog.Info("starting auth server", "port", cfg.HTTPPort, "env", cfg.AppEnv)

	if err := server.Start(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
