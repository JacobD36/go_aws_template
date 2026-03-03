package ports

import (
	"context"
	"messaging-service/internal/domain"
)

// MessageRepository define el puerto para persistir mensajes
type MessageRepository interface {
	Save(ctx context.Context, message *domain.Message) error
	FindByID(ctx context.Context, id string) (*domain.Message, error)
}
