package infrastructure

import (
	"auth-service/internal/application"
	"auth-service/internal/domain"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HTTPHandler maneja las peticiones HTTP del servicio de autenticación
type HTTPHandler struct {
	service *application.AuthService
}

// NewHTTPHandler crea un nuevo manejador HTTP
func NewHTTPHandler(service *application.AuthService) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

// LoginRequest representa la petición de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

// Login maneja el endpoint de autenticación
func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Crear credenciales del dominio
	credentials := &domain.LoginCredentials{
		Email:    req.Email,
		Password: req.Password,
	}

	// Intentar login
	token, err := h.service.Login(context.Background(), credentials)
	if err != nil {
		log.Printf("Login failed: %v", err)

		// Diferenciar errores de validación de errores de credenciales
		if err == domain.ErrInvalidEmail || err == domain.ErrInvalidPassword {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err == domain.ErrInvalidCredentials {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Responder con el token
	response := LoginResponse{
		Token:     token.Token,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck endpoint para verificar el estado del servicio
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// SetupRoutes configura las rutas del servidor
func (h *HTTPHandler) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
	return router
}
