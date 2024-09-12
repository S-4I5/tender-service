package integrational

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity/tender"
	"tender-service/test"
	"testing"
)

func (s *ApiTestSuite) TestCreateTender() {
	orgId := s.createOrganization()
	_ = s.createEmployeeInOrg("test", orgId)

	given := dto.CreateTenderDto{
		Name:            "1",
		Description:     "1",
		Status:          "Created",
		ServiceType:     tender.Construction,
		OrganizationId:  orgId,
		CreatorUsername: "test",
	}

	actual, err := http.Post(s.host+"/tenders/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestCreateTender")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn404WhenCreateTenderAndGroupDontExists() {
	_ = s.createEmployee("test")
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	given := dto.CreateTenderDto{
		Name:            "1",
		Description:     "1",
		Status:          "Created",
		ServiceType:     tender.Construction,
		OrganizationId:  id,
		CreatorUsername: "test",
	}

	actual, err := http.Post(s.host+"/tenders/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn400WhenCreateTenderAndGroupDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 404)
}

func (s *ApiTestSuite) TestReturn401WhenCreateTenderAndEmployeeDontExists() {
	id := s.createOrganization()

	given := dto.CreateTenderDto{
		Name:            "1",
		Description:     "1",
		Status:          "Created",
		ServiceType:     tender.Construction,
		OrganizationId:  id,
		CreatorUsername: "test",
	}

	actual, err := http.Post(s.host+"/tenders/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn401WhenCreateTenderAndEmployeeDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestReturn403WhenCreateTenderAndEmployeeNotInOrganization() {
	orgId := s.createOrganization()
	_ = s.createEmployee("test")

	given := dto.CreateTenderDto{
		Name:            "1",
		Description:     "1",
		Status:          "Created",
		ServiceType:     tender.Construction,
		OrganizationId:  orgId,
		CreatorUsername: "test",
	}

	actual, err := http.Post(s.host+"/tenders/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn403WhenCreateTenderAndEmployeeNotInOrganization")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestGetTendersList() {
	ctx := context.Background()
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  id,
		CreatorUsername: "aboba",
	})

	s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "2",
		Description:     "3",
		Status:          "Published",
		ServiceType:     "Delivery",
		OrganizationId:  id,
		CreatorUsername: "aboba",
	})

	s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "2",
		Description:     "3",
		Status:          "Closed",
		ServiceType:     "Delivery",
		OrganizationId:  id,
		CreatorUsername: "aboba",
	})

	actual, err := http.Get(s.host + "/tenders")
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestGetTendersList")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestGetTendersListWithFilters() {
	testCases := []struct {
		name    string
		filters []tender.ServiceType
	}{
		{name: "Construction", filters: []tender.ServiceType{tender.Construction}},
		{name: "Delivery", filters: []tender.ServiceType{tender.Delivery}},
		{name: "Manufacture", filters: []tender.ServiceType{tender.Manufacture}},
		{name: "ManufactureConstruction", filters: []tender.ServiceType{tender.Manufacture, tender.Construction}},
		{name: "ManufactureDelivery", filters: []tender.ServiceType{tender.Manufacture, tender.Delivery}},
		{name: "All", filters: []tender.ServiceType{tender.Manufacture, tender.Delivery, tender.Construction}},
		{name: "No", filters: []tender.ServiceType{tender.Manufacture, tender.Delivery, tender.Construction}},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()

			id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "1",
				Description:     "2",
				Status:          "Published",
				ServiceType:     tender.Delivery,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})
			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "2",
				Description:     "3",
				Status:          "Created",
				ServiceType:     tender.Delivery,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})

			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "2",
				Description:     "3",
				Status:          "Created",
				ServiceType:     tender.Construction,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})
			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "2",
				Description:     "3",
				Status:          "Published",
				ServiceType:     tender.Construction,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})

			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "2",
				Description:     "3",
				Status:          "Closed",
				ServiceType:     tender.Manufacture,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})
			s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "2",
				Description:     "3",
				Status:          "Published",
				ServiceType:     tender.Manufacture,
				OrganizationId:  id,
				CreatorUsername: "aboba",
			})

			filtersAsString := ""
			for i := 0; i < len(tc.filters); i++ {
				filtersAsString += string(tc.filters[i]) + ","
				fmt.Println("XD")
			}
			filtersAsString = filtersAsString[:len(filtersAsString)-1]
			fmt.Println(filtersAsString)

			actual, err := http.Get(s.host + fmt.Sprintf("/tenders?service_type=%s", filtersAsString))
			if err != nil {
				s.T().Fatalf("Failed to send request: %v", err)
			}
			defer actual.Body.Close()

			expected := test.ReadJson("/tender/response/TestGetTendersListWithFilters/" + tc.name)
			test.ValidateJsonResponse(s.T(), actual, expected, 200)
		})
	}
}

func (s *ApiTestSuite) TestGetTenderStatusSubTest() {
	testCases := []struct {
		name     string
		userName string
		expected tender.Status
	}{
		{name: "WhenOwner", userName: "test", expected: tender.Created},
		{name: "WhenInGroup", userName: "test2", expected: tender.Created},
		{name: "WhenNotInGroupAndTenderPublished", userName: "other", expected: tender.Published},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			orgId := s.createOrganization()
			_ = s.createEmployeeInOrg("test", orgId)
			_ = s.createEmployeeInOrg("test2", orgId)
			s.createEmployee("other")

			tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
				Name:            "1",
				Description:     "2",
				Status:          tc.expected,
				ServiceType:     "Delivery",
				OrganizationId:  orgId,
				CreatorUsername: "test",
			})

			actual, err := http.Get(s.host + fmt.Sprintf("/tenders/%s/status?username=%s", tend.Id.String(), tc.userName))
			if err != nil {
				s.T().Fatalf("Failed to send request: %v", err)
			}
			defer actual.Body.Close()

			test.ValidateResponse(s.T(), actual, string(tc.expected), 200)
		})
	}
}

func (s *ApiTestSuite) TestReturn401WhenGetTenderStatusAndEmployeeDontExists() {
	orgId := s.createOrganization()
	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/tenders/%s/status?username=%s", tend.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn401WhenGetTenderStatusAndEmployeeDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestReturn403WhenGetTenderStatusAndEmployeeNotInGroup() {
	orgId := s.createOrganization()
	_ = s.createEmployeeInOrg("test", orgId)
	_ = s.createEmployee("test2")

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/tenders/%s/status?username=%s", tend.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn403WhenGetTenderStatusAndEmployeeNotInGroup")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestReturn404WhenGetTenderStatusAndTenderDoesNotExists() {
	_ = s.createEmployee("test2")
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	actual, err := http.Get(s.host + fmt.Sprintf("/tenders/%s/status?username=%s", id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn404WhenGetTenderStatusAndTenderDoesNotExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 404)
}

func (s *ApiTestSuite) TestGetMyTenders() {
	orgId := s.createOrganization()
	_ = s.createEmployeeInOrg("test", orgId)

	s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})
	s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Published",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})
	s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Closed",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/tenders/my?username=%s", "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestGetMyTenders")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn401WhenGetMyTendersAndEmployeeDoesNotExists() {
	actual, err := http.Get(s.host + fmt.Sprintf("/tenders/my?username=%s", "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/tender/response/TestReturn401WhenGetMyTendersAndEmployeeDoesNotExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func TestTenderController(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
