package domain

// EmployeeEvent representa un evento relacionado con un empleado
type EmployeeEvent struct {
	EventType string             `json:"event_type"`
	Employee  *EmployeeEventData `json:"employee"`
	Timestamp string             `json:"timestamp"`
}

// EmployeeEventData representa los datos del empleado en el evento (sin información sensible)
type EmployeeEventData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
