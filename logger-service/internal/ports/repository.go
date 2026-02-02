package ports

import (
	"context"
	"logger-service/internal/domain"
)

// LogRepository define el puerto para el repositorio de logs
type LogRepository interface {
	Save(ctx context.Context, entry *domain.LogEntry) error
	FindAll(ctx context.Context) ([]*domain.LogEntry, error)
}
