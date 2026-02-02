package domain

// EmployeeEvent representa un evento relacionado con un empleado
type EmployeeEvent struct {
	EventType string    `json:"event_type"`
	Employee  *Employee `json:"employee"`
	Timestamp string    `json:"timestamp"`
}
