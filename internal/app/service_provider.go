package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tender-service/internal/config"
	"tender-service/internal/controller"
	bid3 "tender-service/internal/controller/bid"
	"tender-service/internal/controller/ping"
	tender3 "tender-service/internal/controller/tender"
	"tender-service/internal/httperr"
	"tender-service/internal/repository"
	"tender-service/internal/repository/bid"
	"tender-service/internal/repository/decision"
	"tender-service/internal/repository/employee"
	"tender-service/internal/repository/feedback"
	"tender-service/internal/repository/organization"
	"tender-service/internal/repository/responsible"
	"tender-service/internal/repository/tender"
	"tender-service/internal/service"
	bid2 "tender-service/internal/service/bid"
	employee2 "tender-service/internal/service/employee"
	organization2 "tender-service/internal/service/organization"
	tender2 "tender-service/internal/service/tender"
)

type serviceProvider struct {
	pool                              *pgxpool.Pool
	config                            config.Config
	pingController                    controller.PingController
	bidController                     controller.BidController
	tenderController                  controller.TenderController
	bidRepository                     repository.BidRepository
	employeeRepository                repository.EmployeeRepository
	decisionRepository                repository.DecisionRepository
	tenderRepository                  repository.TenderRepository
	organizationResponsibleRepository repository.OrganizationResponsibleRepository
	feedbackRepository                repository.FeedbackRepository
	organizationRepository            repository.OrganizationRepository
	tenderService                     service.TenderService
	bidService                        service.BidService
	employeeService                   service.EmployeeService
	organizationService               service.OrganizationService
	handler                           httperr.ApiErrorHandler
}

func newServiceProvider(cfg config.Config) *serviceProvider {
	return &serviceProvider{
		config: cfg,
	}
}

func (s *serviceProvider) Handler() httperr.ApiErrorHandler {
	if s.handler == nil {
		s.handler = httperr.NewApiErrorHandler()
	}
	return s.handler
}

func (s *serviceProvider) PingController() controller.PingController {
	if s.pingController == nil {
		s.pingController = ping.NewPingController()
	}
	return s.pingController
}

func (s *serviceProvider) BidController() controller.BidController {
	if s.bidController == nil {
		s.bidController = bid3.NewBidController(s.BidService(), s.Handler())
	}
	return s.bidController
}

func (s *serviceProvider) TenderController() controller.TenderController {
	if s.tenderController == nil {
		s.tenderController = tender3.NewTenderController(s.TenderService(), s.Handler())
	}
	return s.tenderController
}

func (s *serviceProvider) TenderService() service.TenderService {
	if s.tenderService == nil {
		s.tenderService = tender2.NewTenderService(s.TenderRepository(), s.EmployeeService(), s.OrganizationService())
	}
	return s.tenderService
}

func (s *serviceProvider) BidService() service.BidService {
	if s.bidService == nil {
		s.bidService = bid2.NewBidService(s.EmployeeService(), s.OrganizationService(), s.BidRepository(), s.TenderService(), s.FeedbackRepository(), s.DecisionRepository())
	}
	return s.bidService
}

func (s *serviceProvider) EmployeeService() service.EmployeeService {
	if s.employeeService == nil {
		s.employeeService = employee2.NewEmployeeService(s.EmployeeRepository())
	}
	return s.employeeService
}

func (s *serviceProvider) OrganizationService() service.OrganizationService {
	if s.organizationService == nil {
		s.organizationService = organization2.NewOrganizationService(s.OrganizationRepository(), s.OrganizationResponsibleRepository())
	}
	return s.organizationService
}

func (s *serviceProvider) BidRepository() repository.BidRepository {
	if s.bidRepository == nil {
		s.bidRepository = bid.NewBidRepository(s.Pool())
	}
	return s.bidRepository
}

func (s *serviceProvider) EmployeeRepository() repository.EmployeeRepository {
	if s.employeeRepository == nil {
		s.employeeRepository = employee.NewEmployeeRepository(s.Pool())
	}
	return s.employeeRepository
}

func (s *serviceProvider) DecisionRepository() repository.DecisionRepository {
	if s.decisionRepository == nil {
		s.decisionRepository = decision.NewDecisionRepository(s.Pool())
	}
	return s.decisionRepository
}

func (s *serviceProvider) TenderRepository() repository.TenderRepository {
	if s.tenderRepository == nil {
		s.tenderRepository = tender.NewTenderRepository(s.Pool())
	}
	return s.tenderRepository
}

func (s *serviceProvider) OrganizationResponsibleRepository() repository.OrganizationResponsibleRepository {
	if s.organizationResponsibleRepository == nil {
		s.organizationResponsibleRepository = responsible.NewOrganizationResponsibleRepository(s.Pool())
	}
	return s.organizationResponsibleRepository
}

func (s *serviceProvider) FeedbackRepository() repository.FeedbackRepository {
	if s.feedbackRepository == nil {
		s.feedbackRepository = feedback.NewFeedbackRepository(s.Pool())
	}
	return s.feedbackRepository
}

func (s *serviceProvider) OrganizationRepository() repository.OrganizationRepository {
	if s.organizationRepository == nil {
		s.organizationRepository = organization.NewOrganizationRepository(s.Pool())
	}
	return s.organizationRepository
}

func (s *serviceProvider) Pool() *pgxpool.Pool {
	if s.pool == nil {
		ctx := context.TODO()
		pool, err := pgxpool.New(ctx, s.config.Postgres.Conn)
		if err != nil {
			panic(err.Error())
		}
		s.pool = pool
	}
	return s.pool
}
