package domain

import "time"

// Employee representa la entidad de dominio para un empleado
type Employee struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// NewEmployee crea una nueva instancia de Employee
func NewEmployee(name, email string) *Employee {
	return &Employee{
		Name:      name,
		Email:     email,
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
	return nil
}
