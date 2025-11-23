package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Health check requested")

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"status": "ok"}`))
}
