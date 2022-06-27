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

func TestNewSQLCRepository(t *testing.T) {
	t.Parallel()

	querier := db.New(&pgx.Conn{})
	repo := NewSQLCRepository(querier)

	if reflect.TypeOf(repo) != reflect.TypeOf(&SQLCRepo{}) {
		t.Errorf("returned repo is not of SQLCRepo")
	}
}

func TestSQLCRepo_CreatePaymentPlan(t *testing.T) {
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
				Amount:   amount,
				Status:   "pending",
			},
			mockedQuerier:   makeMockQuerier(createPaymentPlanDBEntity, nil),
			expectedErrType: reflect.TypeOf(nil),
		},
		{
			testName: "error - CreatePaymentPlan failed",
			paramArg: &payments.CreatePlanParams{
				UserID:   uuid.UUID{},
				Currency: "usdc",
				Amount:   amount,
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

func TestSQLCRepo_ListPaymentPlansByUserID(t *testing.T) {
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
		paramUserID     uuid.UUID
		mockedQuerier   db.Querier
		expectedErrType reflect.Type
	}{
		{
			testName:        "happy",
			paramUserID:     userUUID,
			mockedQuerier:   makeMockQuerier(paymentPlanDBEntities, nil),
			expectedErrType: reflect.TypeOf(nil),
		},
		{
			testName:        "error - ListPaymentPlansByUserID failed",
			paramUserID:     userUUID,
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

func TestSQLCRepo_CreatePaymentInstallment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amount, _ := decimal.NewFromString("10.98")
	createPaymentInstallmentDBEntity := &db.CreatePaymentInstallmentsRow{
		ID:            uuid.UUID{},
		PaymentPlanID: uuid.UUID{},
		Currency:      "usdc",
		Amount:        amount,
		DueAt:         time.Time{},
		Status:        "pending",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	makeMockQuerier := func(result *db.CreatePaymentInstallmentsRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().CreatePaymentInstallments(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName        string
		paramArg        *payments.CreateInstallmentParams
		mockedQuerier   db.Querier
		expectedErrType reflect.Type
	}{
		{
			testName: "happy",
			paramArg: &payments.CreateInstallmentParams{
				PaymentPlanID: uuid.UUID{},
				Currency:      "usdc",
				Amount:        amount,
				DueAt:         time.Time{},
				Status:        "pending",
			},
			mockedQuerier:   makeMockQuerier(createPaymentInstallmentDBEntity, nil),
			expectedErrType: reflect.TypeOf(nil),
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
			mockedQuerier:   makeMockQuerier(nil, errors.New("some err")),
			expectedErrType: reflect.TypeOf(pgxDBQueryRunError{}),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			repo := NewSQLCRepository(testcase.mockedQuerier)
			_, err := repo.CreatePaymentInstallment(context.Background(), testcase.paramArg)
			errType := reflect.TypeOf(err)
			if errType != testcase.expectedErrType {
				t.Errorf("expects err type %v but %v returned", testcase.expectedErrType.Name(), errType.Name())
			}
		})
	}
}

func TestSQLCRepo_ListPaymentInstallmentsByPlanID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amount, _ := decimal.NewFromString("10.98")
	paymentInstallmentsDBEntities := []*db.ListPaymentInstallmentsByPlanIDRow{
		{
			ID:            uuid.UUID{},
			PaymentPlanID: uuid.UUID{},
			Currency:      "usdc",
			Amount:        amount,
			DueAt:         time.Time{},
			Status:        "pending",
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
		},
	}

	userUUID := uuid.New()
	makeMockQuerier := func(result []*db.ListPaymentInstallmentsByPlanIDRow, err error) db.Querier {
		mockedQuerier := db.NewMockQuerier(ctrl)
		mockedQuerier.EXPECT().ListPaymentInstallmentsByPlanID(gomock.Any(), gomock.Any()).Return(result, err).AnyTimes()

		return mockedQuerier
	}

	testcases := []struct {
		testName        string
		paramPlanID     uuid.UUID
		mockedQuerier   db.Querier
		expectedErrType reflect.Type
	}{
		{
			testName:        "happy",
			paramPlanID:     userUUID,
			mockedQuerier:   makeMockQuerier(paymentInstallmentsDBEntities, nil),
			expectedErrType: reflect.TypeOf(nil),
		},
		{
			testName:        "error - ListPaymentInstallmentsByPlanID failed",
			paramPlanID:     userUUID,
			mockedQuerier:   makeMockQuerier(nil, errors.New("some db err")),
			expectedErrType: reflect.TypeOf(pgxDBQueryRunError{}),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.testName, func(t *testing.T) {
			t.Parallel()

			repo := NewSQLCRepository(testcase.mockedQuerier)
			_, err := repo.ListPaymentInstallmentsByPlanID(context.Background(), testcase.paramPlanID)
			errType := reflect.TypeOf(err)
			if errType != testcase.expectedErrType {
				t.Errorf("expects err type %v but %v returned", testcase.expectedErrType.Name(), errType.Name())
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
