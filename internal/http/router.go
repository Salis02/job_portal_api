package http

import (
	"job-portal-api/internal/http/middleware"
	"net/http"

	"job-portal-api/internal/modules/healthcheck"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(hcHandler *healthcheck.Handler) http.Handler {
	router := httprouter.New()

	router.GET("/health", hcHandler.Check)

	// Wrap router with CORS
	return middleware.Cors(router)
}
