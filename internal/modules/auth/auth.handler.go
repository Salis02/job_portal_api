package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req RegisterRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	_, err := h.svc.Register(ctx, &req)
	if err != nil {
		http.Error(w, "failed create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	ua := r.UserAgent()
	ip := r.RemoteAddr
	ctx := r.Context()
	resp, err := h.svc.Login(ctx, &req, ua, ip)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// set http-only secure cookie for refresh token (optional) + return body
	// For simplicity, return tokens in body (Postman). In production, put refresh token in HttpOnly cookie.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// accept refresh token in body
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	ua := r.UserAgent()
	ip := r.RemoteAddr
	ctx := context.Background()
	resp, err := h.svc.Refresh(ctx, body.RefreshToken, ua, ip)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := h.svc.Logout(ctx, body.RefreshToken); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}

// Protected route example
func (h *Handler) Me(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := userID.(string)
	name, email, err := h.svc.repo.GetUserById(r.Context(), idStr)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(UserDTO{ID: idStr, Name: name, Email: email})
}
