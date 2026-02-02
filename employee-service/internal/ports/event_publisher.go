package ports

import (
	"context"
	"employee-service/internal/domain"
)

// EventPublisher define el puerto para publicar eventos
type EventPublisher interface {
	PublishEmployeeCreated(ctx context.Context, event *domain.EmployeeEvent) error
}
