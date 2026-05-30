package main

import (
	"log"
	"net/http"

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
	permissionRepo :=
		authorizationpostgres.NewPermissionRepository(
			db,
		)

	rolePermissionRepo :=
		authorizationpostgres.NewRolePermissionRepository(
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
	//services
	permissionBootstrapService :=
		authorizationapp.NewPermissionBootstrapService(
			permissionRepo,
			rolePermissionRepo,
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
		permissionBootstrapService,
	)

	loginUseCase := authapp.NewLoginUseCase(
		identityRepo,
		credentialRepo,
		userRepo,
		userRoleRepo,
		sessionRepo,
	)

	identityGetUserUseCase :=
		identityapp.NewGetUserUseCase(
			userRepo,
		)

	listUsersUseCase :=
		identityapp.NewListUsersUseCase(
			userRepo,
		)

	createUserUseCase :=
		identityapp.NewCreateUserUseCase(
			txManager,
			identityRepo,
			credentialRepo,
			userRepo,
			roleRepo,
			userRoleRepo,
		)

	createRoleUseCase :=
		authorizationapp.NewCreateRoleUseCase(
			roleRepo,
		)

	listRolesUseCase :=
		authorizationapp.NewListRolesUseCase(
			roleRepo,
		)

	assignPermissionUseCase :=
		authorizationapp.NewAssignPermissionToRoleUseCase(
			roleRepo,
			permissionRepo,
			rolePermissionRepo,
		)

	listRolePermissionsUseCase :=
		authorizationapp.NewListRolePermissionsUseCase(
			roleRepo,
			rolePermissionRepo,
		)
	refreshUseCase :=
		authapp.NewRefreshUseCase(
			sessionRepo,
		)

		// handlers
	authHandler := authhttp.NewHandler(
		registerUseCase,
		loginUseCase,
		refreshUseCase,
	)

	identityHandler :=
		identityhttp.NewHandler(
			identityGetUserUseCase,
			listUsersUseCase,
			createUserUseCase,
		)
	authorizationHandler :=
		authorizationhttp.NewHandler(
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

	r.Get(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Write([]byte("ok"))
		},
	)

	r.Post(
		"/register",
		authHandler.Register,
	)

	r.Post(
		"/login",
		authHandler.Login,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(
				next,
			)
		},
	).Get(
		"/me",
		authHandler.Me,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(
				next,
			)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"users:read",
				)(
				next,
			)
		},
	).Get(
		"/users/{userID}",
		identityHandler.GetUser,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(
				next,
			)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"users:read",
				)(
				next,
			)
		},
	).Get(
		"/users",
		identityHandler.ListUsers,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(
				next,
			)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"users:create",
				)(
				next,
			)
		},
	).Post(
		"/users",
		identityHandler.CreateUser,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(next)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"roles:create",
				)(next)
		},
	).Post(
		"/roles",
		authorizationHandler.CreateRole,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(next)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"roles:read",
				)(next)
		},
	).Get(
		"/roles",
		authorizationHandler.ListRoles,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(next)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"roles:update",
				)(next)
		},
	).Post(
		"/roles/{roleID}/permissions",
		authorizationHandler.AssignPermission,
	)

	r.With(
		func(next http.Handler) http.Handler {
			return authMiddleware.Authenticate(next)
		},
		func(next http.Handler) http.Handler {
			return authorizationmiddleware.
				RequirePermission(
					"roles:read",
				)(next)
		},
	).Get(
		"/roles/{roleID}/permissions",
		authorizationHandler.ListRolePermissions,
	)

	r.Post(
		"/refresh",
		authHandler.Refresh,
	)

	log.Printf("starting auth server on port %s", cfg.HTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
