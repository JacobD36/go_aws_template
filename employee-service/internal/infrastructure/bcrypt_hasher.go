package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implementa el puerto PasswordHasher usando bcrypt
// Este es un adaptador que encapsula la lógica de hashing con bcrypt
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher crea una nueva instancia del hasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: bcrypt.DefaultCost, // Cost factor de 10
	}
}

// Hash genera un hash bcrypt del password
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Compare verifica si un password coincide con su hash
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
