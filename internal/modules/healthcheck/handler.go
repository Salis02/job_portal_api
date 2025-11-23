package healthcheck

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.Status())
}
