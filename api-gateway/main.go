package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Employee struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIGateway struct {
	employeeServiceURL string
	authServiceURL     string
}

func NewAPIGateway() *APIGateway {
	employeeServiceURL := os.Getenv("EMPLOYEE_SERVICE_URL")
	if employeeServiceURL == "" {
		employeeServiceURL = "http://localhost:8081"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8082"
	}

	return &APIGateway{
		employeeServiceURL: employeeServiceURL,
		authServiceURL:     authServiceURL,
	}
}

func (gw *APIGateway) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Forward request to employee service (delegar validación al servicio)
	jsonData, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/employees", gw.employeeServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Error calling employee service: %v", err)
		http.Error(w, "Error communicating with employee service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (gw *APIGateway) GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%s/employees", gw.employeeServiceURL))
	if err != nil {
		log.Printf("Error calling employee service: %v", err)
		http.Error(w, "Error communicating with employee service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (gw *APIGateway) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la petición
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Forward request to auth service
	resp, err := http.Post(
		fmt.Sprintf("%s/auth/login", gw.authServiceURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Printf("Error calling auth service: %v", err)
		http.Error(w, "Error communicating with auth service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Leer respuesta del auth service
	responseBody, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func main() {
	gateway := NewAPIGateway()

	router := mux.NewRouter()
	router.HandleFunc("/api/employees", gateway.CreateEmployeeHandler).Methods("POST")
	router.HandleFunc("/api/employees", gateway.GetEmployeesHandler).Methods("GET")
	router.HandleFunc("/api/auth/login", gateway.LoginHandler).Methods("POST")

	log.Println("API Gateway starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
