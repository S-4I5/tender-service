package integrational

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity"
	"tender-service/internal/model/entity/bid"
	"tender-service/internal/model/entity/decision"
	"tender-service/internal/model/entity/tender"
	"tender-service/test"
)

func (s *ApiTestSuite) TestCreateBid() {
	orgId := s.createOrganization()
	empId := s.createEmployeeInOrg("test", orgId)

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	given := dto.CreateBidDto{
		Name:        "1",
		Description: "1",
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    empId,
	}

	actual, err := http.Post(s.host+"/bids/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestCreateTender")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn401WhenCreateBidAndEmployeeDontExists() {
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	given := dto.CreateBidDto{
		Name:        "1",
		Description: "1",
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    id,
	}

	actual, err := http.Post(s.host+"/bids/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn401WhenCreateBidAndEmployeeDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestReturn404WhenCreateBidAndTenderDontExists() {
	empId := s.createEmployee("test")
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	given := dto.CreateBidDto{
		Name:        "1",
		Description: "1",
		TenderId:    id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    empId,
	}

	actual, err := http.Post(s.host+"/bids/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn404WhenCreateBidAndTenderDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 404)
}

func (s *ApiTestSuite) TestReturn403WhenCreateBidByOrgWhenEmployeeNotInOrg() {
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	freeEmpId := s.createEmployee("test2")

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	given := dto.CreateBidDto{
		Name:        "1",
		Description: "1",
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorOrganization,
		AuthorId:    freeEmpId,
	}

	actual, err := http.Post(s.host+"/bids/new", typeJson, test.ToBuffer(given))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn403WhenCreateBidByOrgWhenEmployeeNotInOrg")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestGetUserBids() {
	ctx := context.Background()
	orgId := s.createOrganization()
	empId := s.createEmployeeInOrg("test", orgId)
	emplId2 := s.createEmployee("test2")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    empId,
	})

	s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    emplId2,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/my?username=%s", "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestGetUserBids")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn401WhenGetUserBidsAndEmployeeDontExists() {
	actual, err := http.Get(s.host + fmt.Sprintf("/bids/my?username=%s", "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn401WhenGetUserBidsAndEmployeeDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestGetTenderBids() {
	testCases := []struct {
		name     string
		username string
		status   int
	}{
		{name: "WhenOwner", username: "test", status: 200},
		{name: "WhenEmployeeInOrg", username: "test", status: 200},
		{name: "WhenEmployeeDontExists", username: "something", status: 401},
		{name: "WhenEmployeeNotInOrg", username: "other", status: 403},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			orgId := s.createOrganization()
			s.createEmployeeInOrg("test", orgId)
			s.createEmployeeInOrg("test2", orgId)
			bidCreatorId := s.createEmployee("creator")
			s.createEmployee("other")

			tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "1",
				Description:     "2",
				Status:          "Created",
				ServiceType:     "Delivery",
				OrganizationId:  orgId,
				CreatorUsername: "test",
			})

			s.bidRepository.SaveBid(ctx, bid.Bid{
				Name:        "3",
				Description: "3",
				Status:      bid.Created,
				TenderId:    tend.Id,
				AuthorType:  bid.AuthorUser,
				AuthorId:    bidCreatorId,
			})

			actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/list?username=%s", tend.Id.String(), tc.username))
			if err != nil {
				s.T().Fatalf("Failed to send request: %v", err)
			}
			defer actual.Body.Close()

			expected := test.ReadJson("/bid/response/TestGetTenderBids/" + tc.name)
			test.ValidateJsonResponse(s.T(), actual, expected, tc.status)
		})
	}
}

func (s *ApiTestSuite) TestReturn404WhenGetTenderBidsAndTenderDontExists() {
	s.createEmployee("test")
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/list?username=%s", id, "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn404WhenGetTenderBidsAndTenderDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 404)
}

func (s *ApiTestSuite) TestGetBidStatusWhenBidByGroup() {
	testCases := []struct {
		name     string
		username string
		expected bid.Status
	}{
		{name: "WhenOwner", username: "test", expected: bid.Created},
		{name: "WhenEmployeeInOrg", username: "test2", expected: bid.Created},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			orgId := s.createOrganization()
			bidCreatorId := s.createEmployeeInOrg("test", orgId)
			s.createEmployeeInOrg("test2", orgId)

			tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "1",
				Description:     "2",
				Status:          "Created",
				ServiceType:     "Delivery",
				OrganizationId:  orgId,
				CreatorUsername: "test",
			})

			b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
				Name:        "3",
				Description: "3",
				Status:      bid.Created,
				TenderId:    tend.Id,
				AuthorType:  bid.AuthorOrganization,
				AuthorId:    bidCreatorId,
			})

			actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), tc.username))
			if err != nil {
				s.T().Fatalf("Failed to send request: %v", err)
			}
			defer actual.Body.Close()

			test.ValidateJsonStringResponse(s.T(), actual, string(tc.expected), 200)
		})
	}
}

func (s *ApiTestSuite) TestGetBidStatusWhenBidByUser() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), "creator"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	test.ValidateJsonStringResponse(s.T(), actual, "Created", 200)
}

func (s *ApiTestSuite) TestReturn404WhenGetBidStatusAndBidDontExists() {
	s.createEmployee("test")
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", id, "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn404WhenGetBidStatusAndBidDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 404)
}

func (s *ApiTestSuite) TestReturn403WhenGetBidStatusAndEmployeeNotInOrg() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	s.createEmployee("test2")
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorOrganization,
		AuthorId:    bidCreatorId,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn403WhenGetBidStatusAndEmployeeNotInOrg")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestReturn403WhenGetBidStatusAndEmployeeNotAuthor() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	s.createEmployee("test2")
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn403WhenGetBidStatusAndEmployeeNotAuthor")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestReturn403WhenGetBidStatusAndEmployeeInSameOrgButBidAuthorUser() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	s.createEmployeeInOrg("test2", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn403WhenGetBidStatusAndEmployeeInSameOrgButBidAuthorUser")
	test.ValidateJsonResponse(s.T(), actual, expected, 403)
}

func (s *ApiTestSuite) TestReturn401WhenGetBidStatusAndEmployeeDontExists() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/status?username=%s", b.Id.String(), "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn401WhenGetBidStatusAndEmployeeDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestUpdateBidStatus() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/status?username=%s&status=Published", b.Id.String(), "creator"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestUpdateBidStatus")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn400WhenUpdateBidStatusAndChooseUnselectableStatus() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/status?username=%s&status=Accepted", b.Id.String(), "creator"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn400WhenUpdateBidStatusAndChooseUnselectableStatus")
	test.ValidateJsonResponse(s.T(), actual, expected, 400)
}

func (s *ApiTestSuite) TestEditBid() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, err := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})
	if err != nil {
		log.Println(err.Error())
	}

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})
	if err != nil {
		log.Println(err.Error())
	}

	given := dto.UpdateBidDto{
		Name:        "new",
		Description: "new",
	}

	actual, err := test.HttpPatch(s.host+fmt.Sprintf("/bids/%s/edit?username=%s", b.Id.String(), "creator"), given)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestEditBid")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestRollbackBid() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	s.bidRepository.UpdateBid(ctx, b.Id, "upd", "upd")

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/rollback/1?username=%s", b.Id.String(), "creator"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestRollbackBid")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
}

func (s *ApiTestSuite) TestReturn400WhenRollbackBidAndVersionDontExists() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/rollback/3?username=%s", b.Id.String(), "creator"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn400WhenRollbackBidAndVersionDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 400)
}

func (s *ApiTestSuite) TestPutBidFeedback() {
	ctx := context.Background()
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Created,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/feedback?username=%s&bidFeedback=amazing", b.Id.String(), "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	actualFromDb, _ := s.feedbackRepository.GetFeedbackListForGroup(ctx, tend.Id, bidCreatorId)
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestPutBidFeedback")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
	require.NotEmpty(s.T(), actualFromDb)
}

func (s *ApiTestSuite) TestReturn400WhenPutBidFeedbackAndNoFeedbackGiven() {
	id, _ := uuid.Parse("12d5ca77-d755-49c4-a5ab-1502966ccde0")

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/feedback?username=%s", id, "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn400WhenPutBidFeedbackAndNoFeedbackGiven")
	test.ValidateJsonResponse(s.T(), actual, expected, 400)
}

func (s *ApiTestSuite) TestGetReviews() {
	testCases := []struct {
		name     string
		username string
		status   int
	}{
		{name: "Found", username: "creator", status: 200},
		{name: "NotFound", username: "creator2", status: 404},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			orgId := s.createOrganization()
			s.createEmployeeInOrg("test", orgId)
			s.createEmployeeInOrg("test2", orgId)
			bidCreatorId := s.createEmployee("creator")
			s.createEmployee("creator2")

			tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
				Name:            "1",
				Description:     "2",
				Status:          "Created",
				ServiceType:     "Delivery",
				OrganizationId:  orgId,
				CreatorUsername: "test",
			})

			b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
				Name:        "3",
				Description: "3",
				Status:      bid.Created,
				TenderId:    tend.Id,
				AuthorType:  bid.AuthorUser,
				AuthorId:    bidCreatorId,
			})

			s.feedbackRepository.SaveFeedback(ctx, entity.Feedback{
				BidId:       b.Id,
				Description: "aaaaafa",
				Username:    "test2",
			})

			actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/reviews?authorUsername=%s&requesterUsername=%s", tend.Id.String(), tc.username, "test"))
			if err != nil {
				s.T().Fatalf("Failed to send request: %v", err)
			}
			defer actual.Body.Close()

			expected := test.ReadJson("/bid/response/TestGetReviews/" + tc.name)
			test.ValidateJsonResponse(s.T(), actual, expected, tc.status)
		})
	}
}

func (s *ApiTestSuite) TestReturn401WhenGetReviewsAndRequesterDontExists() {
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/reviews?authorUsername=%s&requesterUsername=%s", tend.Id.String(), "creator", "test2"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn401WhenGetReviewsAndRequesterDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestReturn401WhenGetReviewsAndAuthorDontExists() {
	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)

	tend, _ := s.tenderRepository.SaveTender(context.Background(), tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          "Created",
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	actual, err := http.Get(s.host + fmt.Sprintf("/bids/%s/reviews?authorUsername=%s&requesterUsername=%s", tend.Id.String(), "creator", "test"))
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestReturn401WhenGetReviewsAndAuthorDontExists")
	test.ValidateJsonResponse(s.T(), actual, expected, 401)
}

func (s *ApiTestSuite) TestSubmitDecisionApprove() {
	ctx := context.Background()

	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Published,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	s.bidRepository.UpdateBidStatus(ctx, b.Id, bid.Published)

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/submit_decision?username=%s&decision=Approved", b.Id.String(), "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	actualAmountOfDecisionFromDb, _ := s.decisionRepository.CountDecisionForBid(ctx, b.Id)
	actualTenderFromDb, _ := s.tenderRepository.GetTenderById(ctx, tend.Id)
	actualBidFromDb, _ := s.bidRepository.GetBidById(ctx, b.Id)
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestSubmitDecisionApprove")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
	require.Equal(s.T(), 1, actualAmountOfDecisionFromDb)
	require.Equal(s.T(), bid.Approved, actualBidFromDb.Decision)
	require.Equal(s.T(), tender.Closed, actualTenderFromDb.Status)
}

func (s *ApiTestSuite) TestSubmitDecisionApproveWhenNotEnoughApproves() {
	ctx := context.Background()

	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	s.createEmployeeInOrg("test2", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Published,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	s.bidRepository.UpdateBidStatus(ctx, b.Id, bid.Published)

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/submit_decision?username=%s&decision=Approved", b.Id.String(), "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	actualAmountOfDecisionFromDb, _ := s.decisionRepository.CountDecisionForBid(ctx, b.Id)
	actualTenderFromDb, _ := s.tenderRepository.GetTenderById(ctx, tend.Id)
	actualBidFromDb, _ := s.bidRepository.GetBidById(ctx, b.Id)
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestSubmitDecisionApproveWhenNotEnoughApproves")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
	require.Equal(s.T(), 1, actualAmountOfDecisionFromDb)
	require.Equal(s.T(), bid.None, actualBidFromDb.Decision)
	require.Equal(s.T(), tender.Published, actualTenderFromDb.Status)
}

func (s *ApiTestSuite) TestSubmitDecisionApproveWhenEnoughApproves() {
	ctx := context.Background()

	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	s.createEmployeeInOrg("test2", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Published,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	s.bidRepository.UpdateBidStatus(ctx, b.Id, bid.Published)

	s.decisionRepository.SaveDecision(ctx, decision.Decision{
		Verdict:  decision.Approved,
		Username: "test2",
		BidId:    b.Id,
	})

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/submit_decision?username=%s&decision=Approved", b.Id.String(), "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	actualAmountOfDecisionFromDb, _ := s.decisionRepository.CountDecisionForBid(ctx, b.Id)
	actualTenderFromDb, _ := s.tenderRepository.GetTenderById(ctx, tend.Id)
	actualBidFromDb, _ := s.bidRepository.GetBidById(ctx, b.Id)
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestSubmitDecisionApproveWhenEnoughApproves")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
	require.Equal(s.T(), 2, actualAmountOfDecisionFromDb)
	require.Equal(s.T(), bid.Approved, actualBidFromDb.Decision)
	require.Equal(s.T(), tender.Closed, actualTenderFromDb.Status)
}

func (s *ApiTestSuite) TestSubmitDecisionRejected() {
	ctx := context.Background()

	orgId := s.createOrganization()
	s.createEmployeeInOrg("test", orgId)
	bidCreatorId := s.createEmployee("creator")

	tend, _ := s.tenderRepository.SaveTender(ctx, tender.Tender{
		Name:            "1",
		Description:     "2",
		Status:          tender.Published,
		ServiceType:     "Delivery",
		OrganizationId:  orgId,
		CreatorUsername: "test",
	})

	b, _ := s.bidRepository.SaveBid(ctx, bid.Bid{
		Name:        "3",
		Description: "3",
		Status:      bid.Published,
		TenderId:    tend.Id,
		AuthorType:  bid.AuthorUser,
		AuthorId:    bidCreatorId,
	})

	s.bidRepository.UpdateBidStatus(ctx, b.Id, bid.Published)

	actual, err := test.HttpPut(s.host+fmt.Sprintf("/bids/%s/submit_decision?username=%s&decision=Rejected", b.Id.String(), "test"), nil)
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	actualAmountOfDecisionFromDb, _ := s.decisionRepository.CountDecisionForBid(ctx, b.Id)
	actualTenderFromDb, _ := s.tenderRepository.GetTenderById(ctx, tend.Id)
	actualBidFromDb, _ := s.bidRepository.GetBidById(ctx, b.Id)
	defer actual.Body.Close()

	expected := test.ReadJson("/bid/response/TestSubmitDecisionRejected")
	test.ValidateJsonResponse(s.T(), actual, expected, 200)
	require.Equal(s.T(), 0, actualAmountOfDecisionFromDb)
	require.Equal(s.T(), bid.Rejected, actualBidFromDb.Decision)
	require.Equal(s.T(), tender.Published, actualTenderFromDb.Status)
}
