package sqlc

import (
	"context"
	"flag"
	"fmt"
	"io"
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

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monacohq/golang-common/database/pginit"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"go.uber.org/goleak"
)

// nolint: gochecknoglobals, nolintlint
var (
	testRefDockertestPool     *dockertest.Pool
	testRefDockertestResource *dockertest.Resource
	testRefPoolConn           *pgxpool.Pool
	testRefRepo               *Repo
	testQuerier               *db.Queries
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

	testRefDatabaseURL := fmt.Sprintf("postgres://postgres:%s@%s/datawarehouse?sslmode=disable", "postgres", getHostPort(testRefDockertestResource, "5432/tcp"))

	logger := zerolog.New(io.Discard)

	pgi, err := pginit.New(
		&pginit.Config{
			Host:         "localhost",
			Port:         strings.Split(getHostPort(testRefDockertestResource, "5432/tcp"), ":")[1],
			User:         "postgres",
			Password:     "postgres",
			Database:     "datawarehouse",
			MaxConns:     10,
			MaxIdleConns: 10,
			MaxLifeTime:  1 * time.Minute,
		},
		pginit.WithLogLevel(zerolog.WarnLevel),
		pginit.WithLogger(&logger, "request-id"),
		pginit.WithDecimalType(),
		pginit.WithUUIDType(),
	)
	if err != nil {
		log.Fatalf("Could not init pginit: %s", err)
	}

	if err = testRefDockertestPool.Retry(func() error {
		var poolErr error
		testRefPoolConn, poolErr = pgi.ConnPool(context.Background())

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

	testQuerier = db.New(testRefPoolConn)
	testRefRepo = NewSQLCRepository(testQuerier)

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

	userUUID, _ := uuid.NewV4()

	testcases := []struct {
		testName  string
		paramArg  *payments.CreatePlanParams
		expectErr bool
		expectRow *payments.Plan
	}{
		{
			testName: "happy",
			paramArg: &payments.CreatePlanParams{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(1098, 2),
				Status:   "pending",
			},
			expectErr: false,
			expectRow: &payments.Plan{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(1098, 2),
				Status:   "pending",
			},
		},
		{
			testName: "currency field nil",
			paramArg: &payments.CreatePlanParams{
				UserID: userUUID,
				Amount: *decimal.New(1098, 2),
				Status: "pending",
			},
			expectErr: true,
		},
		{
			testName: "negative amount",
			paramArg: &payments.CreatePlanParams{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(-1099, 2),
				Status:   "pending",
			},
			expectErr: true,
		},
		{
			testName: "decimal wrong precision",
			paramArg: &payments.CreatePlanParams{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(31485937839476927, 16),
				Status:   "pending",
			},
			expectErr: false,
			expectRow: &payments.Plan{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(31485937839476927, 16),
				Status:   "pending",
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			pp, err := testRefRepo.CreatePaymentPlan(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
			if err != nil {
				if !testcase.expectErr {
					t.Errorf("expect no err but err returned")
				}

				return
			}

			if pp.ID == uuid.Nil {
				t.Errorf("expect uuid but nil returned")
			}

			if pp.UserID != testcase.expectRow.UserID {
				t.Errorf("wrong expected user id: got %v, want %v", testcase.expectRow.UserID, pp.UserID)
			}

			if pp.Currency != testcase.expectRow.Currency {
				t.Errorf("wrong expected currency: got %v, want %v", testcase.expectRow.Currency, pp.Currency)
			}

			if pp.Amount.Cmp(&testcase.expectRow.Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", testcase.expectRow.Amount, pp.Amount)
			}

			if pp.Status != testcase.expectRow.Status {
				t.Errorf("wrong expected status: got %v, want %v", testcase.expectRow.Status, pp.Status)
			}
		})
	}
}

func TestSQLCRepo_CreatePaymentPlan_ExistingID(t *testing.T) {
	t.Parallel()

	userUUID, _ := uuid.NewV4()

	existingPlan := createRandomPaymentPlan(t, uuid.Must(uuid.NewV4()))

	testcases := []struct {
		testName  string
		paramArg  *db.CreatePaymentPlanParams
		expectErr bool
		expectRow *db.CreatePaymentPlanParams
	}{
		{
			testName: "plan id already exists",
			paramArg: &db.CreatePaymentPlanParams{
				ID:     existingPlan.ID,
				UserID: userUUID,
				Amount: *decimal.New(1098, 2),
				Status: "pending",
			},
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			_, err := testQuerier.CreatePaymentPlan(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_ListPaymentPlansByUserID_ListOne(t *testing.T) {
	t.Parallel()

	existingPlan := createRandomPaymentPlan(t, uuid.Must(uuid.NewV4()))

	testcases := []struct {
		testName          string
		paramUserID       uuid.UUID
		expectEmptyResult bool
	}{
		{
			testName:          "happy",
			paramUserID:       existingPlan.UserID,
			expectEmptyResult: false,
		},
		{
			testName:          "not found",
			paramUserID:       uuid.Must(uuid.NewV4()),
			expectEmptyResult: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			plans, err := testRefRepo.ListPaymentPlansByUserID(context.Background(), testcase.paramUserID)
			if err != nil {
				t.Fatalf("list payment plans err: %v", err)
			}

			if len(plans) < 1 {
				if !testcase.expectEmptyResult {
					t.Errorf("expect results but empty results returned: %v", plans)
				}

				return
			}

			if existingPlan.ID != plans[0].ID {
				t.Errorf("wrong expected id: got %v, want %v", plans[0].ID, existingPlan.ID)
			}

			if existingPlan.UserID != plans[0].UserID {
				t.Errorf("wrong expected user id: got %v, want %v", plans[0].UserID, existingPlan.UserID)
			}

			if existingPlan.Currency != plans[0].Currency {
				t.Errorf("wrong expected currency: got %v, want %v", plans[0].Currency, existingPlan.Currency)
			}

			if existingPlan.Amount.Cmp(&plans[0].Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", plans[0].Amount, existingPlan.Amount)
			}

			if existingPlan.Status != plans[0].Status {
				t.Errorf("wrong expected status: got %v, want %v", plans[0].Status, existingPlan.Status)
			}
		})
	}
}

func TestSQLCRepo_ListPaymentPlansByUserID_ListMany(t *testing.T) {
	t.Parallel()

	n := 10
	userID := uuid.Must(uuid.NewV4())

	for i := 0; i < n; i++ {
		createRandomPaymentPlan(t, userID)
	}

	testcases := []struct {
		testName        string
		paramUserID     uuid.UUID
		expectResultLen int
	}{
		{
			testName:        "happy",
			paramUserID:     userID,
			expectResultLen: n,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			plans, err := testRefRepo.ListPaymentPlansByUserID(context.Background(), testcase.paramUserID)
			if err != nil {
				t.Fatalf("list payment plans err: %v", err)
			}

			if len(plans) != testcase.expectResultLen {
				t.Errorf("expect %v results but %v results returned", n, len(plans))
			}
		})
	}
}

func TestSQLCRepo_CreatePaymentInstallment(t *testing.T) {
	t.Parallel()

	// create dependent payment entity
	createdPlan := createRandomPaymentPlan(t, uuid.Must(uuid.NewV4()))

	testcases := []struct {
		testName  string
		paramArg  *payments.CreateInstallmentParams
		expectErr bool
	}{
		{
			testName: "happy",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
				Status:        "pending",
			},
			expectErr: false,
		},
		{
			testName: "payment plan does not exist",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: uuid.Must(uuid.NewV4()),
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
				Status:        "pending",
			},
			expectErr: true,
		},
		{
			testName: "amount negative",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        *decimal.New(-1098, 2),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
				Status:        "pending",
			},
			expectErr: true,
		},
		{
			testName: "status field nil",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
			},
			expectErr: true,
		},
		{
			testName: "decimal wrong precision",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        *decimal.New(31485937839476927, 16),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
				Status:        "pending",
			},
			expectErr: false,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			ppi, err := testRefRepo.CreatePaymentInstallment(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
			if err != nil {
				if !testcase.expectErr {
					t.Errorf("expect no err but err returned")
				}

				return
			}
			if ppi.ID == uuid.Nil {
				t.Errorf("expect uuid but nil returned")
			}

			if ppi.Currency != testcase.paramArg.Currency {
				t.Errorf("wrong expected currency: got %v, want %v", testcase.paramArg.Currency, ppi.Currency)
			}

			if ppi.Amount.Cmp(&testcase.paramArg.Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", testcase.paramArg.Amount, ppi.Amount)
			}

			if !ppi.DueAt.Equal(testcase.paramArg.DueAt) {
				t.Errorf("wrong expected due at: got %v, want %v", testcase.paramArg.DueAt, ppi.DueAt)
			}

			if ppi.Status != testcase.paramArg.Status {
				t.Errorf("wrong expected status: got %v, want %v", testcase.paramArg.Status, ppi.Status)
			}
		})
	}
}

func TestSQLCRepo_CreatePaymentInstallment_ExistingID(t *testing.T) {
	t.Parallel()

	// create dependent payment entity
	createdPlan := createRandomPaymentPlan(t, uuid.Must(uuid.NewV4()))
	createdInstallment := createRandomPaymentPlanInstallment(t, createdPlan.ID)

	testcases := []struct {
		testName  string
		paramArg  *db.CreatePaymentInstallmentsParams
		expectErr bool
	}{
		{
			testName: "payment plan id already exists",
			paramArg: &db.CreatePaymentInstallmentsParams{
				ID:            createdInstallment.ID,
				PaymentPlanID: createdPlan.ID,
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         time.Now().UTC().Truncate(time.Microsecond),
				Status:        "pending",
			},
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			_, err := testQuerier.CreatePaymentInstallments(context.Background(), testcase.paramArg)
			if testcase.expectErr && err == nil {
				t.Errorf("expects err but nil returned")
			}
		})
	}
}

func TestSQLCRepo_ListPaymentInstallmentsByPlanID_ListOne(t *testing.T) {
	t.Parallel()

	createdPlan := createRandomPaymentPlan(t, uuid.Must(uuid.NewV4()))
	createdInstallment := createRandomPaymentPlanInstallment(t, createdPlan.ID)

	testcases := []struct {
		testName          string
		paramPlanID       uuid.UUID
		expectEmptyResult bool
	}{
		{
			testName:          "happy",
			paramPlanID:       createdPlan.ID,
			expectEmptyResult: false,
		},
		{
			testName:          "not found",
			paramPlanID:       uuid.Must(uuid.NewV4()),
			expectEmptyResult: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			installments, err := testRefRepo.ListPaymentInstallmentsByPlanID(context.Background(), testcase.paramPlanID)
			if err != nil {
				t.Fatalf("list payment plans err: %v", err)
			}

			if len(installments) < 1 {
				if !testcase.expectEmptyResult {
					t.Errorf("expect results but empty results returned: %v", installments)
				}

				return
			}

			if createdInstallment.Currency != installments[0].Currency {
				t.Errorf("wrong expected currency: got %v, want %v", installments[0].Currency, createdInstallment.Currency)
			}

			if createdInstallment.Amount.Cmp(&installments[0].Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", installments[0].Amount, createdInstallment.Amount)
			}

			if !createdInstallment.DueAt.Equal(installments[0].DueAt) {
				t.Errorf("wrong expected user id: got %v, want %v", installments[0].DueAt, createdInstallment.DueAt)
			}

			if createdInstallment.Status != installments[0].Status {
				t.Errorf("wrong expected status: got %v, want %v", installments[0].Status, createdInstallment.Status)
			}
		})
	}
}

func TestSQLCRepo_ListPaymentInstallmentsByPlanID_ListMany(t *testing.T) {
	t.Parallel()

	n := 10
	userID := uuid.Must(uuid.NewV4())
	plan := createRandomPaymentPlan(t, userID)

	for i := 0; i < n; i++ {
		createRandomPaymentPlanInstallment(t, plan.ID)
	}

	testcases := []struct {
		testName        string
		paramPlanID     uuid.UUID
		expectResultLen int
	}{
		{
			testName:        "happy",
			paramPlanID:     plan.ID,
			expectResultLen: n,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			installments, err := testRefRepo.ListPaymentInstallmentsByPlanID(context.Background(), testcase.paramPlanID)
			if err != nil {
				t.Fatalf("list payment plans err: %v", err)
			}

			if len(installments) != testcase.expectResultLen {
				t.Errorf("expect %v results but %v results returned", n, len(installments))
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

func createRandomPaymentPlan(t *testing.T, userID uuid.UUID) *payments.Plan {
	t.Helper()

	plan, err := testRefRepo.CreatePaymentPlan(context.Background(), &payments.CreatePlanParams{
		UserID:   userID,
		Currency: "usdc",
		Amount:   *decimal.New(1098, 2),
		Status:   "pending",
	})
	if err != nil {
		t.Fatalf("fail to create payment plan: %v", err)
	}

	return plan
}

func createRandomPaymentPlanInstallment(t *testing.T, id uuid.UUID) *payments.Installment {
	t.Helper()

	plan, err := testRefRepo.CreatePaymentInstallment(context.Background(), &payments.CreateInstallmentParams{
		PaymentPlanID: id,
		Currency:      "usdc",
		Amount:        *decimal.New(1098, 2),
		DueAt:         time.Now().UTC().Truncate(time.Microsecond),
		Status:        "pending",
	})
	if err != nil {
		t.Fatalf("fail to create installment: %v", err)
	}

	return plan
}

func getHostPort(resource *dockertest.Resource, id string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		hostAndPort := resource.GetHostPort("5432/tcp")
		hp := strings.Split(hostAndPort, ":")
		testRefHost := hp[0]
		testRefPort := hp[1]

		return testRefHost + ":" + testRefPort
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}

	return u.Hostname() + ":" + resource.GetPort(id)
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
