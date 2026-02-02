package application

import (
	"context"
	"fmt"
	"log"
	"logger-service/internal/domain"
	"logger-service/internal/ports"
	"time"

	"github.com/google/uuid"
)

// LoggerService implementa la lógica de negocio para el logger
type LoggerService struct {
	repository ports.LogRepository
	consumer   ports.EventConsumer
}

// NewLoggerService crea una nueva instancia del servicio
func NewLoggerService(repo ports.LogRepository, consumer ports.EventConsumer) *LoggerService {
	return &LoggerService{
		repository: repo,
		consumer:   consumer,
	}
}

// ProcessEvent procesa un evento de empleado
func (s *LoggerService) ProcessEvent(ctx context.Context, event *domain.EmployeeEvent) error {
	// Parsear el timestamp del evento
	timestamp, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	// Crear entrada de log
	logEntry := domain.NewLogEntry(
		event.EventType,
		event.Employee.ID,
		event.Employee.Name,
		event.Employee.Email,
		timestamp,
	)
	logEntry.ID = uuid.New().String()

	// Guardar en la base de datos
	if err := s.repository.Save(ctx, logEntry); err != nil {
		return fmt.Errorf("error saving log entry: %w", err)
	}

	// Mostrar en consola
	s.displayEventInfo(logEntry)

	return nil
}

// displayEventInfo muestra la información del evento en consola
func (s *LoggerService) displayEventInfo(entry *domain.LogEntry) {
	log.Println("========================================")
	log.Printf("EVENTO RECIBIDO: %s", entry.EventType)
	log.Printf("ID Empleado: %s", entry.EmployeeID)
	log.Printf("Nombre: %s", entry.Name)
	log.Printf("Email: %s", entry.Email)
	log.Printf("Timestamp del evento: %s", entry.Timestamp.Format("2006-01-02 15:04:05"))
	log.Printf("Procesado el: %s", entry.ProcessedAt.Format("2006-01-02 15:04:05"))
	log.Println("========================================")
}

// StartConsuming inicia el consumo de eventos
func (s *LoggerService) StartConsuming(ctx context.Context) error {
	log.Println("Logger service started consuming events...")

	return s.consumer.ConsumeEvents(ctx, func(event *domain.EmployeeEvent) error {
		return s.ProcessEvent(ctx, event)
	})
}
