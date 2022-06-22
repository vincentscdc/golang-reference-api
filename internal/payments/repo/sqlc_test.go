package repo

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/db"
	"golangreferenceapi/internal/payments"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
)

// type anything struct{}

func TestNewPGXRepository(t *testing.T) {
	t.Parallel()

	querier := db.New(&pgx.Conn{})
	repo := NewSQLCRepository(querier)

	if reflect.TypeOf(repo) != reflect.TypeOf(&SQLCRepo{}) {
		t.Errorf("returned repo is not of SQLCRepo")
	}
}

func TestPGXRepo_CreatePaymentPlan(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amount, _ := decimal.NewFromString("10.98")
	createPaymentPlanDBEntity := &db.CreatePaymentPlanRow{
		ID:        uuid.UUID{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Currency:  "usdc",
		UserID:    uuid.UUID{},
		Amount:    amount,
		Status:    "pending",
	}

	makeMockQuerier := func(result *db.CreatePaymentPlanRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().CreatePaymentPlan(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName        string
		paramArg        *payments.CreatePlanParams
		mockedQuerier   db.Querier
		expectedErrType reflect.Type
	}{
		{
			testName: "happy",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   "10.98",
				Status:   "pending",
			},
			mockedQuerier:   makeMockQuerier(createPaymentPlanDBEntity, nil),
			expectedErrType: reflect.TypeOf(nil),
		},
		{
			testName: "error - parse string to decimal failed",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   "tendotninty8",
				Status:   "pending",
			},
			mockedQuerier:   makeMockQuerier(createPaymentPlanDBEntity, nil),
			expectedErrType: reflect.TypeOf(decimalParseError{}),
		},
		{
			testName: "error - CreatePaymentPlan failed",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   "10.98",
				Status:   "pending",
			},
			mockedQuerier:   makeMockQuerier(nil, errors.New("some db err")),
			expectedErrType: reflect.TypeOf(pgxDBQueryRunError{}),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			repo := NewSQLCRepository(testcase.mockedQuerier)
			_, err := repo.CreatePaymentPlan(context.Background(), testcase.paramArg)
			errType := reflect.TypeOf(err)
			if errType != testcase.expectedErrType {
				t.Errorf("expects err type %v but %v returned", testcase.expectedErrType.Name(), errType.Name())
			}
		})
	}
}

func TestPGXRepo_ListPaymentPlansByUserID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amount, _ := decimal.NewFromString("10.98")
	paymentPlanDBEntities := []*db.ListPaymentPlansByUserIDRow{
		{
			ID:        uuid.UUID{},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			Currency:  "usdc",
			UserID:    uuid.UUID{},
			Amount:    amount,
			Status:    "pending",
		},
	}

	userUUID := uuid.New()
	makeMockQuerier := func(result []*db.ListPaymentPlansByUserIDRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().ListPaymentPlansByUserID(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName        string
		paramUserID     string
		mockedQuerier   db.Querier
		expectedErrType reflect.Type
	}{
		{
			testName:        "happy",
			paramUserID:     userUUID.String(),
			mockedQuerier:   makeMockQuerier(paymentPlanDBEntities, nil),
			expectedErrType: reflect.TypeOf(nil),
		},
		{
			testName:        "error - UUID parse failed",
			paramUserID:     "userUUID1",
			mockedQuerier:   makeMockQuerier(paymentPlanDBEntities, nil),
			expectedErrType: reflect.TypeOf(uuidParseError{}),
		},
		{
			testName:        "error - ListPaymentPlansByUserID failed",
			paramUserID:     userUUID.String(),
			mockedQuerier:   makeMockQuerier(nil, errors.New("some db err")),
			expectedErrType: reflect.TypeOf(pgxDBQueryRunError{}),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			repo := NewSQLCRepository(testcase.mockedQuerier)
			_, err := repo.ListPaymentPlansByUserID(context.Background(), testcase.paramUserID)
			errType := reflect.TypeOf(err)
			if errType != testcase.expectedErrType {
				t.Errorf("expects err type %v but %v returned", testcase.expectedErrType.Name(), errType.Name())
			}
		})
	}
}

func TestPGXRepo_newPlanFromDBEntity(t *testing.T) {
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
