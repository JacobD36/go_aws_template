package domain

import "time"

// LogEntry representa una entrada de log en el sistema
type LogEntry struct {
	ID          string    `json:"id"`
	EventType   string    `json:"event_type"`
	EmployeeID  string    `json:"employee_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessedAt time.Time `json:"processed_at"`
}

// NewLogEntry crea una nueva entrada de log
func NewLogEntry(eventType, employeeID, name, email string, timestamp time.Time) *LogEntry {
	return &LogEntry{
		EventType:   eventType,
		EmployeeID:  employeeID,
		Name:        name,
		Email:       email,
		Timestamp:   timestamp,
		ProcessedAt: time.Now(),
	}
}
