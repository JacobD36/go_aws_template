package infrastructure

import (
	"context"
	"log"
	"messaging-service/internal/domain"
	"time"
)

// SimulatedMessageSender implementa el envío simulado de mensajes
type SimulatedMessageSender struct{}

// NewSimulatedMessageSender crea una nueva instancia del sender simulado
func NewSimulatedMessageSender() *SimulatedMessageSender {
	return &SimulatedMessageSender{}
}

// SendEmail simula el envío de un email
func (s *SimulatedMessageSender) SendEmail(ctx context.Context, message *domain.Message) error {
	log.Printf("=== SIMULACIÓN DE ENVÍO DE EMAIL ===")
	log.Printf("Para: %s", message.To)
	log.Printf("Asunto: %s", message.Subject)
	log.Printf("Cuerpo:\n%s", message.Body)
	log.Printf("===================================")

	// Simular un pequeño delay como si estuviera enviando realmente
	time.Sleep(100 * time.Millisecond)

	log.Printf("✓ Email enviado exitosamente a %s", message.To)
	return nil
}

// SendSMS simula el envío de un SMS
func (s *SimulatedMessageSender) SendSMS(ctx context.Context, message *domain.Message) error {
	log.Printf("=== SIMULACIÓN DE ENVÍO DE SMS ===")
	log.Printf("Para: %s", message.To)
	log.Printf("Mensaje: %s", message.Body)
	log.Printf("==================================")

	// Simular un pequeño delay como si estuviera enviando realmente
	time.Sleep(100 * time.Millisecond)

	log.Printf("✓ SMS enviado exitosamente a %s", message.To)
	return nil
}
