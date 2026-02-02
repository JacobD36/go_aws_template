package infrastructure

import (
	"context"
	"employee-service/internal/application"
	"employee-service/internal/domain"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HTTPHandler maneja las peticiones HTTP
type HTTPHandler struct {
	service *application.EmployeeService
}

// NewHTTPHandler crea un nuevo manejador HTTP
func NewHTTPHandler(service *application.EmployeeService) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

type CreateEmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateEmployee maneja la creación de un empleado
func (h *HTTPHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	employee, err := h.service.CreateEmployee(context.Background(), req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("Error creating employee: %v", err)
		// Diferenciar errores de validación del dominio
		if err == domain.ErrInvalidPassword || err == domain.ErrInvalidName || err == domain.ErrInvalidEmail {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver versión pública sin password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee.ToPublic())
}

// GetEmployees obtiene todos los empleados
func (h *HTTPHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.GetAllEmployees(context.Background())
	if err != nil {
		log.Printf("Error getting employees: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir a versión pública sin passwords
	publicEmployees := make([]*domain.EmployeePublic, len(employees))
	for i, emp := range employees {
		publicEmployees[i] = emp.ToPublic()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicEmployees)
}

// SetupRoutes configura las rutas del servidor
func (h *HTTPHandler) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/employees", h.CreateEmployee).Methods("POST")
	router.HandleFunc("/employees", h.GetEmployees).Methods("GET")
	return router
}
