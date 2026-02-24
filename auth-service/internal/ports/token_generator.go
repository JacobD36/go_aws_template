package ports

import "auth-service/internal/domain"

// TokenGenerator define el puerto para generar tokens JWT
// Aplica el patrón Strategy y el principio de Inversión de Dependencias
type TokenGenerator interface {
	// GenerateToken crea un token JWT con el ID del usuario
	GenerateToken(userID string) (*domain.AuthToken, error)

	// ValidateToken valida un token JWT y retorna el ID del usuario
	ValidateToken(token string) (string, error)
}
