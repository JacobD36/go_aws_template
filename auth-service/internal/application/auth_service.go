package application

import (
	"auth-service/internal/domain"
	"auth-service/internal/ports"
	"context"
	"log"
)

// AuthService implementa la lógica de negocio para autenticación
type AuthService struct {
	repository     ports.UserRepository
	passwordHasher ports.PasswordHasher
	tokenGenerator ports.TokenGenerator
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokenGen ports.TokenGenerator,
) *AuthService {
	return &AuthService{
		repository:     repo,
		passwordHasher: hasher,
		tokenGenerator: tokenGen,
	}
}

// Login autentica un usuario y genera un token JWT
func (s *AuthService) Login(ctx context.Context, credentials *domain.LoginCredentials) (*domain.AuthToken, error) {
	// Validar credenciales
	if err := credentials.Validate(); err != nil {
		return nil, err
	}

	// Buscar usuario por email
	user, err := s.repository.FindByEmail(ctx, credentials.Email)
	if err != nil {
		log.Printf("User not found: %s", credentials.Email)
		return nil, domain.ErrInvalidCredentials
	}

	// Comparar password (usando el puerto PasswordHasher)
	err = s.passwordHasher.Compare(user.Password, credentials.Password)
	if err != nil {
		log.Printf("Invalid password for user: %s", credentials.Email)
		return nil, domain.ErrInvalidCredentials
	}

	// Generar token JWT (usando el puerto TokenGenerator)
	token, err := s.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, domain.ErrTokenGeneration
	}

	log.Printf("User authenticated successfully: %s (ID: %s)", user.Email, user.ID)
	return token, nil
}

// ValidateToken valida un token JWT y retorna el ID del usuario
func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	return s.tokenGenerator.ValidateToken(token)
}
