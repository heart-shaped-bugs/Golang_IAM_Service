package handlers

import (
	"encoding/json"
	"errors"
	"iam-service/api/http/models"
	auth "iam-service/internal/usecases"
	"net/http"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := models.RegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.authService.Register(req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			http.Error(w, auth.ErrUserAlreadyExists.Error(), http.StatusConflict)
		case errors.Is(err, auth.ErrEmailIsNotValid):
			http.Error(w, auth.ErrEmailIsNotValid.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := models.LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			http.Error(w, auth.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
