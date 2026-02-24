package ports

import (
	"auth-service/internal/domain"
	"context"
)

// UserRepository define el puerto para el repositorio de usuarios
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
