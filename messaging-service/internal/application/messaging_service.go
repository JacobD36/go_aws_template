package application

import (
	"context"
	"log"
	"messaging-service/internal/domain"
	"messaging-service/internal/ports"
	"time"
)

// MessagingService implementa la lógica de negocio para el enví de mensajes
type MessagingService struct {
	repository ports.MessageRepository
	sender     ports.MessageSender
	publisher  ports.EventPublisher
}

// NewMessagingService crea una nueva instancia del servicio de mensajería
func NewMessagingService(
	repository ports.MessageRepository,
	sender ports.MessageSender,
	publisher ports.EventPublisher,
) *MessagingService {
	return &MessagingService{
		repository: repository,
		sender:     sender,
		publisher:  publisher,
	}
}

// ProcessEmployeeCreatedEvent procesa un evento de empleado creado
func (s *MessagingService) ProcessEmployeeCreatedEvent(ctx context.Context, event *domain.EmployeeEvent) error {
	log.Printf("Processing employee created event for: %s (%s)", event.Employee.Name, event.Employee.Email)

	// Verificar el tipo de evento
	if event.EventType != "employee.created" {
		log.Printf("Ignoring event type: %s", event.EventType)
		return nil
	}

	// Crear mensaje de bienvenida
	message := domain.NewWelcomeEmail(event.Employee.ID, event.Employee.Name, event.Employee.Email)

	// Enviar el mensaje según su tipo
	var err error
	switch message.Type {
	case domain.MessageTypeEmail:
		err = s.sender.SendEmail(ctx, message)
	case domain.MessageTypeSMS:
		err = s.sender.SendSMS(ctx, message)
	default:
		log.Printf("Unsupported message type: %s", message.Type)
		return domain.ErrInvalidMessage
	}

	if err != nil {		log.Printf("Error sending message: %v", err)
		message.Status = "failed"
		return domain.ErrMessageSendFailed
	}

	// Guardar el mensaje en el repositorio
	if err := s.repository.Save(ctx, message); err != nil {
		log.Printf("Error saving message to repository: %v", err)
		// No retornamos error aquí porque el mensaje ya fue enviado
	}

	// Publicar evento de log compatible con logger-service
	logEvent := &domain.LogEvent{
		EventType: "message.sent",
		Employee: domain.Employee{
			ID:        message.ID,
			Name:      "Messaging Service - Welcome Email",
			Email:     message.To,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := s.publisher.PublishLogEvent(ctx, logEvent); err != nil {
		log.Printf("Error publishing log event: %v", err)
		// No retornamos error aquí porque el mensaje ya fue enviado
	}

	log.Printf("Welcome message sent successfully to %s", message.To)
	return nil
}

// HandleEmployeeEvent es el handler para procesar eventos de empleado
func (s *MessagingService) HandleEmployeeEvent(event *domain.EmployeeEvent) error {
	ctx := context.Background()
	return s.ProcessEmployeeCreatedEvent(ctx, event)
}
