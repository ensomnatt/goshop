package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"users_service/internal/config"
	"users_service/internal/usecase"
	"users_service/pkg/logger"
)

type UsersHandler struct {
	UC  *usecase.UsersUseCase
	log *logger.Logger
	cfg *config.Config
}

func NewUsersHandler(uc *usecase.UsersUseCase, cfg *config.Config, log *logger.Logger) *UsersHandler {
	return &UsersHandler{
		UC:  uc,
		cfg: cfg,
		log: log,
	}
}

func (h *UsersHandler) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)

	return http.ListenAndServe(fmt.Sprintf(":%s", h.cfg.Port), mux)
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.Errorf("Failed to register: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.UC.Register(req.Name, req.Email, req.Password, r.Context())
	if err != nil {
		h.log.Errorf("Failed to register: %v", err)
		http.Error(w, "Failed to register", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.Errorf("Failed to login: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, expirationTime, err := h.UC.Login(req.Email, req.Password, r.Context())
	if err != nil {
		h.log.Errorf("Failed to login: %v", err)
		http.Error(w, "Failed to login", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token:     token,
		ExpiresAt: expirationTime,
	})
}
