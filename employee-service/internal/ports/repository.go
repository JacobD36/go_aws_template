package ports

import (
	"context"
	"employee-service/internal/domain"
)

// EmployeeRepository define el puerto para el repositorio de empleados
type EmployeeRepository interface {
	Save(ctx context.Context, employee *domain.Employee) error
	FindByID(ctx context.Context, id string) (*domain.Employee, error)
	FindAll(ctx context.Context) ([]*domain.Employee, error)
}
