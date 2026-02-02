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
}

func NewAPIGateway() *APIGateway {
	serviceURL := os.Getenv("EMPLOYEE_SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:8081"
	}
	return &APIGateway{
		employeeServiceURL: serviceURL,
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

func main() {
	gateway := NewAPIGateway()

	router := mux.NewRouter()
	router.HandleFunc("/api/employees", gateway.CreateEmployeeHandler).Methods("POST")
	router.HandleFunc("/api/employees", gateway.GetEmployeesHandler).Methods("GET")

	log.Println("API Gateway starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
