package domain

// LogEvent representa un evento de logging compatible con el logger-service
type LogEvent struct {
	EventType string   `json:"event_type"`
	Employee  Employee `json:"employee"`
	Timestamp string   `json:"timestamp"`
}
