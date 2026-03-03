package ports

import (
	"context"
	"messaging-service/internal/domain"
)

// MessageSender define el puerto para enviar mensajes
type MessageSender interface {
	SendEmail(ctx context.Context, message *domain.Message) error
	SendSMS(ctx context.Context, message *domain.Message) error
}
