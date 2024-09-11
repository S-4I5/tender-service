package employee

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tender-service/internal/model"
	"tender-service/internal/model/entity"
	"tender-service/internal/repository"
)

var (
	errEmployeeDoesNotExists = fmt.Errorf("employee does not exists")
)

type service struct {
	employeeRepository repository.EmployeeRepository
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository) *service {
	return &service{employeeRepository: employeeRepository}
}

func (s *service) GetEmployeeByUsername(ctx context.Context, username string) (entity.Employee, error) {
	return s.employeeRepository.GetEmployeeByUsername(ctx, username)
}

func (s *service) ValidateEmployeeExists(ctx context.Context, username string) error {
	op := "employee_service.validate_employee_exists"

	exists, err := s.employeeRepository.EmployeeExistByUsername(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return model.NewNotAuthorizedError(op, errEmployeeDoesNotExists)
	}

	return nil
}

func (s *service) GetEmployeeByUsernameById(ctx context.Context, id uuid.UUID) (entity.Employee, error) {
	return s.employeeRepository.GetEmployeeById(ctx, id)
}

func (s *service) ValidateEmployeeExistsById(ctx context.Context, id uuid.UUID) error {
	exists, err := s.employeeRepository.EmployeeExistById(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errEmployeeDoesNotExists
	}

	return nil
}
