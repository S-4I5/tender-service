package organization

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tender-service/internal/model"
	"tender-service/internal/repository"
)

type service struct {
	organizationRepository            repository.OrganizationRepository
	organizationResponsibleRepository repository.OrganizationResponsibleRepository
}

var (
	errNotInOrganization    = fmt.Errorf("given user not in given organization")
	errOrganizationNotFound = fmt.Errorf("organization not found")
	errEmployeeNotInOrg     = fmt.Errorf("employee not in organization")
)

func NewOrganizationService(
	organizationRepository repository.OrganizationRepository,
	organizationResponsibleRepository repository.OrganizationResponsibleRepository,
) *service {
	return &service{
		organizationRepository:            organizationRepository,
		organizationResponsibleRepository: organizationResponsibleRepository,
	}
}

func (s *service) UsersHasSimilarOrganization(ctx context.Context, userId uuid.UUID, username string) (bool, error) {
	return s.organizationResponsibleRepository.UsersHasSimilarOrganization(ctx, userId, username)
}

func (s *service) ValidateOrganizationExists(ctx context.Context, id uuid.UUID) error {
	op := "organization_service.validate_organization_exists"
	exists, err := s.organizationRepository.OrganizationExistById(ctx, id)
	if !exists {
		return model.NewNotFoundError(op, errOrganizationNotFound)
	}
	return err
}

func (s *service) ValidateEmployeeInAnyOrganization(ctx context.Context, userId uuid.UUID) error {
	op := "organization_service.validate_employee_is_is_any_organization"

	result, err := s.organizationResponsibleRepository.IsEmployeeInAnyOrganization(ctx, userId)
	if !result {
		return model.NewForbiddenError(op, errEmployeeNotInOrg)
	}
	return err
}

func (s *service) ValidateEmployeeBelongsToOrganization(ctx context.Context, orgId uuid.UUID, username string) error {
	op := "organization_service.validate_employee_belongs_to_organization"

	if err := s.ValidateOrganizationExists(ctx, orgId); err != nil {
		return err
	}

	result, err := s.organizationResponsibleRepository.IsResponsibleInOrganization(ctx, username, orgId)
	if !result {
		return model.NewForbiddenError(op, errNotInOrganization)
	}
	return err
}

func (s *service) GetOrganizationEmployeeCount(ctx context.Context, id uuid.UUID) (int, error) {
	count, err := s.organizationResponsibleRepository.CountEmployeesInOrganization(ctx, id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
