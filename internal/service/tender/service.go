package tender

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tender-service/internal/mapper"
	"tender-service/internal/model"
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/repository"
	service2 "tender-service/internal/service"
	"tender-service/internal/util"
)

type service struct {
	tenderRepository    repository.TenderRepository
	employeeService     service2.EmployeeService
	organizationService service2.OrganizationService
}

var errTenderVersionDoesNotExists = fmt.Errorf("given tender version dont exists")

func NewTenderService(
	tenderRepository repository.TenderRepository,
	employeeService service2.EmployeeService,
	organizationService service2.OrganizationService,
) *service {
	return &service{
		tenderRepository:    tenderRepository,
		employeeService:     employeeService,
		organizationService: organizationService,
	}
}

func (s *service) ValidateTenderExists(ctx context.Context, tenderId uuid.UUID) error {
	_, err := s.tenderRepository.GetTenderById(ctx, tenderId)
	return err
}

func (s *service) GetTenderById(ctx context.Context, tenderId uuid.UUID) (tender.Tender, error) {
	return s.tenderRepository.GetTenderById(ctx, tenderId)
}

func (s *service) GetTenders(ctx context.Context, page util.Page, serviceTypes []tender.ServiceType) ([]dto.TenderDto, error) {
	tenders, err := s.tenderRepository.GetTenderList(ctx, page, serviceTypes, "", true)
	if err != nil {
		return nil, err
	}

	return mapper.TenderListToTenderDtoList(tenders), nil
}

func (s *service) CreateNewTender(ctx context.Context, tenderDto dto.CreateTenderDto) (dto.TenderDto, error) {
	if err := s.employeeService.ValidateEmployeeExists(ctx, tenderDto.CreatorUsername); err != nil {
		return dto.TenderDto{}, err
	}

	err := s.organizationService.ValidateEmployeeBelongsToOrganization(ctx, tenderDto.OrganizationId, tenderDto.CreatorUsername)
	if err != nil {
		return dto.TenderDto{}, err
	}

	entity := mapper.CreateTenderDtoToTender(tenderDto)

	saved, err := s.tenderRepository.SaveTender(ctx, entity)
	if err != nil {
		fmt.Println(err.Error())
		return dto.TenderDto{}, err
	}

	return mapper.TenderToTenderDto(saved), nil
}

func (s *service) GetUserTenders(ctx context.Context, page util.Page, username string) ([]dto.TenderDto, error) {
	if err := s.employeeService.ValidateEmployeeExists(ctx, username); err != nil {
		return nil, err
	}

	tenders, err := s.tenderRepository.GetTenderList(ctx, page, nil, username, false)
	if err != nil {
		return nil, err
	}

	return mapper.TenderListToTenderDtoList(tenders), nil
}

func (s *service) GetTenderStatus(ctx context.Context, tenderId uuid.UUID, username string) (tender.Status, error) {
	entity, err := s.tenderRepository.GetTenderById(ctx, tenderId)
	if err != nil {
		return "", err
	}

	if entity.Status == tender.Published {
		return tender.Published, nil
	}

	err = s.ValidateEmployeeRightsOnTender(ctx, tenderId, username)
	if err != nil {
		return "", err
	}

	return entity.Status, nil
}

func (s *service) UpdateTenderStatus(ctx context.Context, tenderId uuid.UUID, username string, status tender.Status) (dto.TenderDto, error) {
	err := s.ValidateEmployeeRightsOnTender(ctx, tenderId, username)
	if err != nil {
		fmt.Println("e2")
		return dto.TenderDto{}, err
	}

	updated, err := s.tenderRepository.UpdateTenderStatus(ctx, tenderId, status)
	if err != nil {
		fmt.Println("ee", err)
		return dto.TenderDto{}, err
	}

	return mapper.TenderToTenderDto(updated), nil
}

func (s *service) EditTender(ctx context.Context, tenderDto dto.UpdateTenderDto, tenderId uuid.UUID, username string) (dto.TenderDto, error) {
	err := s.ValidateEmployeeRightsOnTender(ctx, tenderId, username)
	if err != nil {
		fmt.Println("XD")
		return dto.TenderDto{}, err
	}

	updated, err := s.tenderRepository.UpdateTender(ctx, tenderId, tenderDto.Name, tenderDto.Description, tenderDto.ServiceType)
	if err != nil {
		fmt.Println("XD1", err)
		return dto.TenderDto{}, err
	}

	return mapper.TenderToTenderDto(updated), nil
}

func (s *service) RollbackTender(ctx context.Context, tenderId uuid.UUID, username string, version int) (dto.TenderDto, error) {
	op := "tender_service.rollback_tender"

	tend, err := s.tenderRepository.GetTenderById(ctx, tenderId)
	if err != nil {
		return dto.TenderDto{}, err
	}

	if tend.Version < version {
		return dto.TenderDto{}, model.NewBadRequestError(op, errTenderVersionDoesNotExists)
	}

	err = s.ValidateEmployeeRightsOnTender(ctx, tenderId, username)
	if err != nil {
		return dto.TenderDto{}, err
	}

	updated, err := s.tenderRepository.RollbackTender(ctx, tenderId, version)
	if err != nil {
		return dto.TenderDto{}, err
	}
	return mapper.TenderToTenderDto(updated), err
}

func (s *service) ValidateEmployeeRightsOnTender(ctx context.Context, tenderId uuid.UUID, username string) error {
	curTender, err := s.tenderRepository.GetTenderById(ctx, tenderId)
	if err != nil {
		return err
	}

	if err = s.employeeService.ValidateEmployeeExists(ctx, username); err != nil {
		return err
	}

	err = s.organizationService.ValidateEmployeeBelongsToOrganization(ctx, curTender.OrganizationId, username)
	if err != nil {
		return err
	}
	return nil
}
