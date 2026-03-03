package ports

import (
	"context"
	"messaging-service/internal/domain"
)

// EventPublisher define el puerto para publicar eventos
type EventPublisher interface {
	PublishLogEvent(ctx context.Context, event *domain.LogEvent) error
}
