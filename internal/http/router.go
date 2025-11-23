package http

import (
	"job-portal-api/internal/http/middleware"
	"net/http"

	"job-portal-api/internal/modules/auth"
	"job-portal-api/internal/modules/healthcheck"

	"github.com/julienschmidt/httprouter"
)

// assume we already have pg pool and created repo/service/handler

func NewRouter(hcHandler *healthcheck.Handler,
	authRepo *auth.Repo,
	jwtMgr *auth.JWTManager) http.Handler {
	router := httprouter.New()
	authSvc := auth.NewService(authRepo, jwtMgr)
	authHandler := auth.NewHandler(authSvc)

	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/refresh", authHandler.Refresh)
	router.POST("/api/auth/logout", authHandler.Logout)

	// protected
	meHandler := authHandler.Me
	router.Handler("GET", "/api/auth/me", auth.JWTMiddleware(jwtMgr)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// convert httprouter style: we can adapt; if using julienschmidt/httprouter then Register like earlier.
		meHandler(w, r, nil) // adjust accordingly for your router signature
	})))

	router.GET("/health", hcHandler.Check)

	// Wrap router with CORS
	return middleware.Cors(router)
}
