package domain

import "time"

// MessageType representa el tipo de mensaje a enviar
type MessageType string

const (
	MessageTypeEmail MessageType = "EMAIL"
	MessageTypeSMS   MessageType = "SMS"
)

// Message representa un mensaje a enviar
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	To        string      `json:"to"`
	Subject   string      `json:"subject,omitempty"`
	Body      string      `json:"body"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	SentAt    *time.Time  `json:"sent_at,omitempty"`
}

// NewWelcomeEmail crea un mensaje de bienvenida por email
func NewWelcomeEmail(userID, name, email string) *Message {
	return &Message{
		ID:        generateMessageID(userID),
		Type:      MessageTypeEmail,
		To:        email,
		Subject:   "¡Bienvenido a nuestro sistema!",
		Body:      buildWelcomeEmailBody(name),
		Status:    "pending",
		CreatedAt: time.Now(),
	}
}

// generateMessageID genera un ID único para el mensaje
func generateMessageID(userID string) string {
	return "msg-" + userID + "-" + time.Now().Format("20060102150405")
}

// buildWelcomeEmailBody construye el cuerpo del email de bienvenida
func buildWelcomeEmailBody(name string) string {
	return `Hola ` + name + `,

¡Bienvenido a nuestro sistema! Estamos encantados de tenerte con nosotros.

Tu cuenta ha sido creada exitosamente y ya puedes comenzar a utilizar todos nuestros servicios.

Si tienes alguna pregunta, no dudes en contactarnos.

Saludos cordiales,
El equipo`
}
