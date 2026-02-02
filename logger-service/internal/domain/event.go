package domain

// EmployeeEvent representa un evento relacionado con un empleado
type EmployeeEvent struct {
	EventType string   `json:"event_type"`
	Employee  Employee `json:"employee"`
	Timestamp string   `json:"timestamp"`
}

// Employee representa los datos básicos de un empleado en el evento
type Employee struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
