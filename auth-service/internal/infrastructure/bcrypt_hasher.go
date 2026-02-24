package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implementa el puerto PasswordHasher usando bcrypt
type BcryptPasswordHasher struct{}

// NewBcryptPasswordHasher crea una nueva instancia del hasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

// Compare verifica si un password coincide con su hash
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
