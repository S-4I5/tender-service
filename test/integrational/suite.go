package integrational

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"math/rand"
	"reflect"
	"tender-service/internal/app"
	"tender-service/internal/config"
	"tender-service/internal/repository"
	"unsafe"
)

const typeJson = "application/json"

type ApiTestSuite struct {
	suite.SetupAllSuite
	suite.Suite
	container        *postgres.PostgresContainer
	app              *app.App
	pool             *pgxpool.Pool
	host             string
	tenderRepository repository.TenderRepository
}

func (s *ApiTestSuite) SetupSuite() {
	fmt.Println("before")

	ctx := context.Background()

	psql, _ := postgres.RunContainer(
		ctx,
		postgres.WithDatabase("tender-service"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	s.container = psql

	host, _ := s.container.Host(ctx)

	conn := "postgres://postgres:postgres@" + host + "/tender-service"

	randomPort := rand.Intn(65535-1024+1) + 1024

	s.host = fmt.Sprintf("http://localhost:%d/api", randomPort)

	curApp, err := app.NewApp(context.Background(), config.Config{
		Server: config.ServerConfig{Address: fmt.Sprintf(":%d", randomPort)},
		Postgres: config.PostgresConfig{
			MigrationsDir: "../../migrations/test",
			Conn:          conn,
		},
	})
	if err != nil {
		fmt.Println(err)
		s.T().Fatalf("XD")
	}

	go func() {
		err = curApp.Run()
		if err != nil {
			log.Printf("failed to run app: %s\n", err.Error())
		}
	}()

	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		panic(err.Error())
	}

	s.pool = pool
	s.app = curApp

	v := reflect.ValueOf(s.app).Elem()
	providerField := v.FieldByName("provider")
	providerValue := reflect.NewAt(providerField.Type(), unsafe.Pointer(providerField.UnsafeAddr())).Elem()
	tenderRepoField := providerValue.Elem().FieldByName("tenderRepository")
	s.tenderRepository = reflect.NewAt(tenderRepoField.Type(), unsafe.Pointer(tenderRepoField.UnsafeAddr())).Elem().Interface().(repository.TenderRepository)

	//a, _ := s.tenderRepository.GetTenderList(ctx, util.NewPage(0, 1), nil, "", false)
	//fmt.Println(a)
}

func (s *ApiTestSuite) TearDownSuite() {
	_ = s.container.Terminate(context.Background())
	_ = s.app.Stop()
	_ = s.pool
}

func (s *ApiTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Println("clear")
	_, _ = s.pool.Exec(context.Background(),
		"TRUNCATE employee, organization, organization_responsible, tender, tender_version, bid, bid_version, decisions, feedback;")
}

func (s *ApiTestSuite) SetupSubTest() {
	fmt.Println("clear sub")
	_, _ = s.pool.Exec(context.Background(),
		"TRUNCATE employee, organization, organization_responsible, tender, tender_version, bid, bid_version, decisions, feedback;")
}

func (s *ApiTestSuite) createEmployeeInOrg(username string, orgId uuid.UUID) uuid.UUID {
	id := s.createEmployee(username)
	s.bindEmployeeToOrg(id, orgId)
	return id
}

func (s *ApiTestSuite) createEmployee(username string) uuid.UUID {
	builder := squirrel.Insert("employee").PlaceholderFormat(squirrel.Dollar).
		Columns("username").Values(username).
		Suffix("RETURNING id")

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("Failed to builder: %s", err)
	}

	rows, err := s.pool.Query(context.Background(), sql, args...)
	if err != nil {
		log.Fatalf("Failed to exec: %s", err)
	}

	rows.Next()

	var id uuid.UUID
	if err = rows.Scan(&id); err != nil {
		log.Fatalf("Failed to scan: %s", err)
	}

	rows.Close()

	return id
}

func (s *ApiTestSuite) createOrganization() uuid.UUID {
	builder := squirrel.Insert("organization").PlaceholderFormat(squirrel.Dollar).
		Columns("name").Values("test").
		Suffix("RETURNING id")

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("2Failed to builder: %s", err)
	}

	rows, err := s.pool.Query(context.Background(), sql, args...)
	if err != nil {
		log.Fatalf("2Failed to exec: %s", err)
	}

	rows.Next()

	var id uuid.UUID
	if err = rows.Scan(&id); err != nil {
		log.Fatalf("2Failed to scan: %s", err)
	}

	rows.Close()

	return id
}

func (s *ApiTestSuite) bindEmployeeToOrg(emplId, orgId uuid.UUID) {
	builder := squirrel.Insert("organization_responsible").PlaceholderFormat(squirrel.Dollar).
		Columns("organization_id", "user_id").Values(orgId.String(), emplId.String())

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("3Failed to builder: %s", err)
	}

	rows, err := s.pool.Query(context.Background(), sql, args...)
	if err != nil {
		log.Fatalf("3Failed to exec: %s", err)
	}

	rows.Close()
}
