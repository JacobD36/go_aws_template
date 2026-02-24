package domain

import "time"

// User representa un usuario en el sistema de autenticación
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hash del password (nunca se serializa)
	CreatedAt time.Time `json:"created_at"`
}

// Validate valida los datos básicos del usuario
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}
