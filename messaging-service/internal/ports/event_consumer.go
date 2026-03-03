package ports

import (
	"context"
	"messaging-service/internal/domain"
)

// EventConsumer define el puerto para consumir eventos
type EventConsumer interface {
	ConsumeEvents(ctx context.Context, handler func(*domain.EmployeeEvent) error) error
}
