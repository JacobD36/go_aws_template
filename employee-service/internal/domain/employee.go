package domain

import (
	"regexp"
	"time"
)

// Employee representa la entidad de dominio para un empleado
type Employee struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hash del password (nunca se serializa en JSON)
	CreatedAt time.Time `json:"created_at"`
}

// EmployeePublic representa un empleado sin información sensible
type EmployeePublic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToPublic convierte un Employee a EmployeePublic (sin password)
func (e *Employee) ToPublic() *EmployeePublic {
	return &EmployeePublic{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		CreatedAt: e.CreatedAt,
	}
}

// NewEmployee crea una nueva instancia de Employee
func NewEmployee(name, email, password string) *Employee {
	return &Employee{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

// Validate valida los datos del empleado
func (e *Employee) Validate() error {
	if e.Name == "" {
		return ErrInvalidName
	}
	if e.Email == "" {
		return ErrInvalidEmail
	}
	if e.Password == "" {
		return ErrInvalidPassword
	}
	if err := ValidatePassword(e.Password); err != nil {
		return err
	}
	return nil
}

// ValidatePassword valida la complejidad del password
// Debe tener mínimo 8 caracteres, una letra mayúscula, un número y un caracter especial
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}

	// Verificar al menos una letra mayúscula
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return ErrInvalidPassword
	}

	// Verificar al menos un número
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return ErrInvalidPassword
	}

	// Verificar al menos un caracter especial
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~]`).MatchString(password)
	if !hasSpecial {
		return ErrInvalidPassword
	}

	return nil
}
