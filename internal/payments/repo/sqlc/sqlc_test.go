package sqlc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"golangreferenceapi/database"
	"golangreferenceapi/internal/db"
	"golangreferenceapi/internal/payments"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/shopspring/decimal"
	"go.uber.org/goleak"
)

// nolint: gochecknoglobals, nolintlint
var (
	testRefHost, testRefPort, testRefDatabaseURL string
	testRefDockertestPool                        *dockertest.Pool
	testRefDockertestResource                    *dockertest.Resource
	testRefPoolConn                              *pgxpool.Pool
	testRefRepo                                  *Repo
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	var err error

	testRefDockertestPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testRefDockertestPool.MaxWait = 120 * time.Second

	testRefDockertestResource, err = testRefDockertestPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=datawarehouse",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start testRefDockertestResource: %s", err)
	}

	_ = testRefDockertestResource.Expire(180)

	testRefDatabaseURL = fmt.Sprintf("postgres://postgres:%s@%s/datawarehouse?sslmode=disable", "postgres", getHostPort(testRefDockertestResource, "5432/tcp"))

	if err = testRefDockertestPool.Retry(func() error {
		var poolErr error
		testRefPoolConn, poolErr = pgxpool.Connect(context.Background(), testRefDatabaseURL)
		if poolErr != nil {
			return poolErr
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer func() {
		testRefPoolConn.Close()
	}()

	err = runMigrations(testRefDatabaseURL)
	if err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}

	querier := db.New(testRefPoolConn)
	testRefRepo = NewSQLCRepository(querier)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := testRefDockertestPool.Purge(testRefDockertestResource); err != nil {
		log.Fatalf("Could not purge testRefDockertestResource: %s", err)
	}

	os.Exit(code)
}

func TestNewSQLCRepository(t *testing.T) {
	t.Parallel()

	querier := db.New(&pgx.Conn{})
	repo := NewSQLCRepository(querier)

	if reflect.TypeOf(repo) != reflect.TypeOf(&Repo{}) {
		t.Errorf("returned testRefRepo is not of Repo")
	}
}

func TestSQLCRepo_CreatePaymentPlan(t *testing.T) {
	t.Parallel()

	// mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	makeMockQuerier := func(result *db.CreatePaymentPlanRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().CreatePaymentPlan(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	// actual db setup
	dbConn, dbURL, err := createNewDatabaseAndConn(testRefDockertestResource, testRefPoolConn, "create_payment_plan_db")
	if err != nil {
		t.Fatalf("Failed to create new database: %s", err)
	}

	// run migrations
	err = runMigrations(dbURL)
	if err != nil {
		t.Fatalf("Failed to run migrations: %s", err)
	}

	querier := db.New(dbConn)
	amount, _ := decimal.NewFromString("10.98")

	testcases := []struct {
		testName      string
		paramArg      *payments.CreatePlanParams
		mockedQuerier db.Querier
		expectErr     bool
	}{
		{
			testName: "happy",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   amount,
				Status:   "pending",
			},
			mockedQuerier: nil,
			expectErr:     false,
		},
		{
			testName: "error - CreatePaymentPlan db error",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   amount,
				Status:   "pending",
			},
			mockedQuerier: makeMockQuerier(nil, errors.New("some db err")),
			expectErr:     true,
		},
		{
			testName: "error - CreatePaymentPlan db error",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   amount,
				Status:   "pending",
			},
			mockedQuerier: makeMockQuerier(nil, errors.New("some err")),
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			var repo *Repo

			if testcase.mockedQuerier == nil {
				repo = NewSQLCRepository(querier)
			} else {
				repo = NewSQLCRepository(testcase.mockedQuerier)
			}

			_, err := repo.CreatePaymentPlan(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_ListPaymentPlansByUserID(t *testing.T) {
	t.Parallel()

	// mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// actual db setup
	dbConn, dbURL, err := createNewDatabaseAndConn(testRefDockertestResource, testRefPoolConn, "list_payment_plans_by_userid_db")
	if err != nil {
		t.Fatalf("Failed to create new database: %s", err)
	}

	// run migrations
	err = runMigrations(dbURL)
	if err != nil {
		t.Fatalf("Failed to run migrations: %s", err)
	}

	userUUID := uuid.Must(uuid.NewV4())

	querier := db.New(dbConn)
	makeMockQuerier := func(result []*db.ListPaymentPlansByUserIDRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().ListPaymentPlansByUserID(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName      string
		paramUserID   uuid.UUID
		mockedQuerier db.Querier
		expectErr     bool
	}{
		{
			testName:      "happy",
			paramUserID:   userUUID,
			mockedQuerier: nil,
			expectErr:     false,
		},
		{
			testName:      "error - ListPaymentPlansByUserID failed",
			paramUserID:   userUUID,
			mockedQuerier: makeMockQuerier(nil, errors.New("some db err")),
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			var repo *Repo

			if testcase.mockedQuerier == nil {
				repo = NewSQLCRepository(querier)
			} else {
				repo = NewSQLCRepository(testcase.mockedQuerier)
			}

			_, err := repo.ListPaymentPlansByUserID(context.Background(), testcase.paramUserID)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_CreatePaymentInstallment(t *testing.T) {
	t.Parallel()

	// mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// actual db setup
	dbConn, dbURL, err := createNewDatabaseAndConn(testRefDockertestResource, testRefPoolConn, "create_payment_installments_db")
	if err != nil {
		t.Fatalf("Failed to create new database: %s", err)
	}

	// run migrations
	err = runMigrations(dbURL)
	if err != nil {
		t.Fatalf("Failed to run migrations: %s", err)
	}

	querier := db.New(dbConn)
	amount, _ := decimal.NewFromString("10.98")

	// create dependent payment entity
	createdPlan, _ := querier.CreatePaymentPlan(context.Background(), &db.CreatePaymentPlanParams{
		ID:       uuid.UUID{},
		UserID:   uuid.UUID{},
		Currency: "usdc",
		Amount:   amount,
		Status:   "pending",
	})

	makeMockQuerier := func(result *db.CreatePaymentInstallmentsRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().CreatePaymentInstallments(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName      string
		paramArg      *payments.CreateInstallmentParams
		mockedQuerier db.Querier
		expectErr     bool
	}{
		{
			testName: "happy",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        amount,
				DueAt:         time.Time{},
				Status:        "pending",
			},
			mockedQuerier: nil,
			expectErr:     false,
		},
		{
			testName: "error - CreatePaymentInstallments failed",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: uuid.UUID{},
				Currency:      "usdc",
				Amount:        amount,
				DueAt:         time.Time{},
				Status:        "pending",
			},
			mockedQuerier: makeMockQuerier(nil, errors.New("some err")),
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			var repo *Repo

			if testcase.mockedQuerier == nil {
				repo = NewSQLCRepository(querier)
			} else {
				repo = NewSQLCRepository(testcase.mockedQuerier)
			}

			_, err := repo.CreatePaymentInstallment(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_ListPaymentInstallmentsByPlanID(t *testing.T) {
	t.Parallel()

	// mock setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// actual db setup
	dbConn, dbURL, err := createNewDatabaseAndConn(testRefDockertestResource, testRefPoolConn, "list_payment_installments_by_planid_db")
	if err != nil {
		t.Fatalf("Failed to create new database: %s", err)
	}

	// run migrations
	err = runMigrations(dbURL)
	if err != nil {
		t.Fatalf("Failed to run migrations: %s", err)
	}

	querier := db.New(dbConn)
	amount, _ := decimal.NewFromString("10.98")

	// create dependent payment entity
	createdPlan, _ := querier.CreatePaymentPlan(context.Background(), &db.CreatePaymentPlanParams{
		ID:       uuid.UUID{},
		UserID:   uuid.UUID{},
		Currency: "usdc",
		Amount:   amount,
		Status:   "pending",
	})

	makeMockQuerier := func(result []*db.ListPaymentInstallmentsByPlanIDRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().ListPaymentInstallmentsByPlanID(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName      string
		paramPlanID   uuid.UUID
		mockedQuerier db.Querier
		expectErr     bool
	}{
		{
			testName:      "happy",
			paramPlanID:   createdPlan.ID,
			mockedQuerier: nil,
			expectErr:     false,
		},
		{
			testName:      "error - ListPaymentInstallmentsByPlanID failed",
			paramPlanID:   createdPlan.ID,
			mockedQuerier: makeMockQuerier(nil, errors.New("some db err")),
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			var repo *Repo

			if testcase.mockedQuerier == nil {
				repo = NewSQLCRepository(querier)
			} else {
				repo = NewSQLCRepository(testcase.mockedQuerier)
			}

			_, err := repo.ListPaymentInstallmentsByPlanID(context.Background(), testcase.paramPlanID)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_newPlanFromDBEntity(t *testing.T) {
	t.Parallel()

	querier := db.New(&pgx.Conn{})
	repo := NewSQLCRepository(querier)

	type unsupportedStruct struct{}

	testcases := []struct {
		testName      string
		paramDBEntity interface{}
		expectErr     bool
	}{
		{
			testName:      "happy - CreatePaymentPlanRow",
			paramDBEntity: &db.CreatePaymentPlanRow{},
			expectErr:     false,
		},
		{
			testName:      "happy - ListPaymentPlansByUserIDRow",
			paramDBEntity: &db.ListPaymentPlansByUserIDRow{},
			expectErr:     false,
		},
		{
			testName:      "happy - PaymentPlan",
			paramDBEntity: &db.PaymentPlan{},
			expectErr:     false,
		},
		{
			testName:      "failed - unsupported DB entity type",
			paramDBEntity: &unsupportedStruct{},
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			plan, err := repo.newPlanFromDBEntity(testcase.paramDBEntity)

			if testcase.expectErr {
				if err == nil {
					t.Errorf("expected err but nil returned")
				}
			} else {
				if reflect.TypeOf(plan) != reflect.TypeOf(&payments.Plan{}) {
					t.Errorf("returned entity is not of *payments.Plan")
				}
			}
		})
	}
}

func TestSQLCRepo_newInstallmentFromDBEntity(t *testing.T) {
	t.Parallel()

	querier := db.New(&pgx.Conn{})
	repo := NewSQLCRepository(querier)

	testcases := []struct {
		testName      string
		paramDBEntity interface{}
		expectErr     bool
	}{
		{
			testName:      "happy - CreatePaymentInstallmentsRow",
			paramDBEntity: &db.CreatePaymentInstallmentsRow{},
			expectErr:     false,
		},
		{
			testName:      "happy - ListPaymentInstallmentsByPlanIDRow",
			paramDBEntity: &db.ListPaymentInstallmentsByPlanIDRow{},
			expectErr:     false,
		},
		{
			testName:      "happy - PaymentInstallment",
			paramDBEntity: &db.PaymentInstallment{},
			expectErr:     false,
		},
		{
			testName:      "failed - json err",
			paramDBEntity: "some invalid input",
			expectErr:     true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			plan, err := repo.newInstallmentFromDBEntity(testcase.paramDBEntity)

			if testcase.expectErr {
				if err == nil {
					t.Errorf("expected err but nil returned")
				}
			} else {
				if reflect.TypeOf(plan) != reflect.TypeOf(&payments.Installment{}) {
					t.Errorf("returned entity is not of *payments.Plan")
				}
			}
		})
	}
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		hostAndPort := resource.GetHostPort("5432/tcp")
		hp := strings.Split(hostAndPort, ":")
		testRefHost = hp[0]
		testRefPort = hp[1]

		return testRefHost + ":" + testRefPort
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}

	testRefHost = u.Hostname()
	testRefPort = resource.GetPort(id)

	return u.Hostname() + ":" + resource.GetPort(id)
}

func createNewDatabaseAndConn(resource *dockertest.Resource, dbConn *pgxpool.Pool, dbName string) (*pgxpool.Pool, string, error) {
	_, err := dbConn.Exec(context.Background(), "CREATE DATABASE "+dbName+";")
	if err != nil {
		return nil, "", err
	}

	dbURL := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", "postgres", getHostPort(resource, "5432/tcp"), dbName)

	newConn, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, "", err
	}

	return newConn, dbURL, nil
}

func runMigrations(dbURL string) error {
	d, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
