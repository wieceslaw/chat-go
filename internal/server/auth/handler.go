package auth

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service UserService
}

func NewAuthHanlder(service UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", h.handleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.service.Register(r.Context(), &RegisterUser{
		Name:     req.Username,
		Password: req.Password,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := RegisterResponseDto{
		Message: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := h.service.Login(r.Context(), &LoginData{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := LoginResponseDto{
		AuthToken: string(*token),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
