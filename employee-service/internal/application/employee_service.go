package application

import (
	"context"
	"employee-service/internal/domain"
	"employee-service/internal/ports"
	"time"

	"github.com/google/uuid"
)

// EmployeeService implementa la lógica de negocio para empleados
type EmployeeService struct {
	repository ports.EmployeeRepository
	publisher  ports.EventPublisher
}

// NewEmployeeService crea una nueva instancia del servicio
func NewEmployeeService(repo ports.EmployeeRepository, pub ports.EventPublisher) *EmployeeService {
	return &EmployeeService{
		repository: repo,
		publisher:  pub,
	}
}

// CreateEmployee crea un nuevo empleado
func (s *EmployeeService) CreateEmployee(ctx context.Context, name, email string) (*domain.Employee, error) {
	employee := domain.NewEmployee(name, email)

	if err := employee.Validate(); err != nil {
		return nil, err
	}

	employee.ID = uuid.New().String()

	if err := s.repository.Save(ctx, employee); err != nil {
		return nil, err
	}

	// Publicar evento
	event := &domain.EmployeeEvent{
		EventType: "employee.created",
		Employee:  employee,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := s.publisher.PublishEmployeeCreated(ctx, event); err != nil {
		return nil, err
	}

	return employee, nil
}

// GetAllEmployees obtiene todos los empleados
func (s *EmployeeService) GetAllEmployees(ctx context.Context) ([]*domain.Employee, error) {
	return s.repository.FindAll(ctx)
}

// GetEmployeeByID obtiene un empleado por su ID
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	return s.repository.FindByID(ctx, id)
}
