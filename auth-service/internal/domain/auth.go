package domain

// LoginCredentials representa las credenciales de login
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate valida las credenciales de login
func (c *LoginCredentials) Validate() error {
	if c.Email == "" {
		return ErrInvalidEmail
	}
	if c.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}

// AuthToken representa el token de autenticación generado
type AuthToken struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}
