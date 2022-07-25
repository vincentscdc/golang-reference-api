package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/payments"
	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/mock/repomock"
	"golangreferenceapi/internal/payments/repo/sqlc"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
)

func TestPaymentServiceImp_GetPaymentPlanByUserID(t *testing.T) {
	t.Parallel()

	var (
		decimalAmount = *decimal.New(1000, 0)
		userID        = uuid.Must(uuid.NewV4())
		planID        = uuid.Must(uuid.NewV4())
		installmentID = uuid.Must(uuid.NewV4())
		ctx           = context.Background()
		createdAt, _  = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		updatedAt, _  = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		dueAt, _      = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		currency      = "usdc"
		status        = "pending"
		paymentPlans  = []*payments.Plan{
			{
				ID:        planID,
				UserID:    userID,
				Currency:  currency,
				Amount:    decimalAmount,
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
		paymentInstallments = []*payments.Installment{
			{
				ID:            installmentID,
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt,
				Status:        status,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		}
		paymentPlanResponse = []PaymentPlans{
			{
				ID:          planID.String(),
				UserID:      userID.String(),
				Currency:    currency,
				TotalAmount: decimalAmount.String(),
				Status:      status,
				CreatedAt:   createdAt.Format(common.TimeFormat),
				Installments: []PaymentPlanInstallment{
					{
						ID:       installmentID.String(),
						Amount:   decimalAmount.String(),
						Currency: currency,
						DueAt:    dueAt.Format(common.TimeFormat),
						Status:   status,
					},
				},
			},
		}
	)

	type args struct {
		userID uuid.UUID
	}

	tests := []struct {
		name    string
		prepare func(rm *repomock.MockRepository)
		args    args
		want    []PaymentPlans
		wantErr bool
	}{
		{
			name: "happy path",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(paymentPlans, nil),
					rm.EXPECT().ListPaymentInstallmentsByPlanID(ctx, gomock.Eq(planID)).Return(paymentInstallments, nil),
				)
			},
			args: args{
				userID: userID,
			},
			want:    paymentPlanResponse,
			wantErr: false,
		},
		{
			name: "ListPaymentPlansByUserID error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				userID: userID,
			},
			wantErr: true,
		},
		{
			name: "ListPaymentInstallmentsByPlanID error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(paymentPlans, nil),
					rm.EXPECT().ListPaymentInstallmentsByPlanID(ctx, gomock.Eq(planID)).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				userID: userID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repomock.NewMockRepository(ctrl)

			if tt.prepare != nil {
				tt.prepare(repo)
			}

			p := &PaymentServiceImp{repository: repo}

			got, err := p.GetPaymentPlanByUserID(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentServiceImp.GetPaymentPlanByUserID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaymentServiceImp.GetPaymentPlanByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentServiceImp_CreatePendingPaymentPlan(t *testing.T) {
	t.Parallel()

	var (
		decimalAmount  = *decimal.New(1000, 0)
		userID         = uuid.Must(uuid.NewV4())
		planID         = uuid.Must(uuid.NewV4())
		installmentID  = uuid.Must(uuid.NewV4())
		installmentID2 = uuid.Must(uuid.NewV4())
		ctx            = context.Background()
		createdAt, _   = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		updatedAt, _   = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		dueAt, _       = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		currency       = "usdc"
		status         = "pending"

		paymentPlanParamMock = &payments.CreatePlanParams{
			UserID:   userID,
			Currency: currency,
			Amount:   decimalAmount,
			Status:   status,
		}

		paymentPlanMock = &payments.Plan{
			ID:        planID,
			UserID:    userID,
			Currency:  currency,
			Amount:    decimalAmount,
			Status:    status,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		paymentInstallmentParamMock = []*payments.CreateInstallmentParams{
			{
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt,
				Status:        status,
			},
			{
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt.Add(1 * time.Hour),
				Status:        status,
			},
		}

		paymentInstallmentMock = []*payments.Installment{
			{
				ID:            installmentID,
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt,
				Status:        status,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
			{
				ID:            installmentID2,
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt.Add(1 * time.Hour),
				Status:        status,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		}

		paymentPlanParams = &CreatePaymentPlanParams{
			UserID:      userID,
			Currency:    "usdc",
			TotalAmount: "1000",
			Installments: []PaymentPlanInstallmentParams{
				{
					Currency: "usdc",
					Amount:   "1000",
					DueAt:    dueAt,
				},
				{
					Currency: "usdc",
					Amount:   "1000",
					DueAt:    dueAt.Add(1 * time.Hour),
				},
			},
		}

		paymentPlanResponse = &PaymentPlans{
			ID:          planID.String(),
			UserID:      userID.String(),
			Currency:    currency,
			TotalAmount: decimalAmount.String(),
			Status:      status,
			CreatedAt:   createdAt.Format(common.TimeFormat),
			Installments: []PaymentPlanInstallment{
				{
					ID:       installmentID.String(),
					Amount:   decimalAmount.String(),
					Currency: currency,
					DueAt:    dueAt.Format(common.TimeFormat),
					Status:   status,
				},
				{
					ID:       installmentID2.String(),
					Amount:   decimalAmount.String(),
					Currency: currency,
					DueAt:    dueAt.Add(1 * time.Hour).Format(common.TimeFormat),
					Status:   status,
				},
			},
		}
	)

	type args struct {
		createPaymentPlanParams *CreatePaymentPlanParams
	}

	tests := []struct {
		name    string
		prepare func(rm *repomock.MockRepository)
		args    args
		want    *PaymentPlans
		wantErr bool
	}{
		{
			name: "happy path",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().CreatePaymentPlan(ctx, gomock.Eq(paymentPlanParamMock)).Return(paymentPlanMock, nil),
					rm.EXPECT().CreatePaymentInstallment(ctx, gomock.Eq(paymentInstallmentParamMock[0])).Return(paymentInstallmentMock[0], nil),
					rm.EXPECT().CreatePaymentInstallment(ctx, gomock.Eq(paymentInstallmentParamMock[1])).Return(paymentInstallmentMock[1], nil),
				)
			},
			args: args{
				createPaymentPlanParams: paymentPlanParams,
			},
			want:    paymentPlanResponse,
			wantErr: false,
		},
		{
			name: "CreatePaymentPlan error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().CreatePaymentPlan(ctx, gomock.Eq(paymentPlanParamMock)).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				createPaymentPlanParams: paymentPlanParams,
			},
			wantErr: true,
		},
		{
			name: "CreatePaymentInstallment error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().CreatePaymentPlan(ctx, gomock.Eq(paymentPlanParamMock)).Return(paymentPlanMock, nil),
					rm.EXPECT().CreatePaymentInstallment(ctx, gomock.Eq(paymentInstallmentParamMock[0])).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				createPaymentPlanParams: paymentPlanParams,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repomock.NewMockRepository(ctrl)

			if tt.prepare != nil {
				tt.prepare(repo)
			}

			p := &PaymentServiceImp{repository: repo}

			got, err := p.CreatePendingPaymentPlan(ctx, tt.args.createPaymentPlanParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentServiceImp.CreatePendingPaymentPlan() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaymentServiceImp.CreatePendingPaymentPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentServiceImp_CompletePaymentPlanCreation(t *testing.T) {
	t.Parallel()

	var (
		decimalAmount = *decimal.New(1000, 0)
		userID        = uuid.Must(uuid.NewV4())
		planID        = uuid.Must(uuid.NewV4())
		installmentID = uuid.Must(uuid.NewV4())
		ctx           = context.Background()
		createdAt, _  = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		updatedAt, _  = time.Parse(common.TimeFormat, "2021-10-10T23:00:00Z")
		dueAt, _      = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		currency      = "usdc"
		status        = "pending"
		paymentPlans  = []*payments.Plan{
			{
				ID:        uuid.Must(uuid.NewV4()),
				UserID:    userID,
				Currency:  currency,
				Amount:    decimalAmount,
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			{
				ID:        planID,
				UserID:    userID,
				Currency:  currency,
				Amount:    decimalAmount,
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
		paymentInstallments = []*payments.Installment{
			{
				ID:            installmentID,
				PaymentPlanID: planID,
				Currency:      currency,
				Amount:        decimalAmount,
				DueAt:         dueAt,
				Status:        status,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		}

		paymentPlanParams = &CompletePaymentPlanParams{
			UserID: userID,
		}

		paymentPlanResponse = &PaymentPlans{
			ID:          planID.String(),
			UserID:      userID.String(),
			Currency:    currency,
			TotalAmount: decimalAmount.String(),
			Status:      status,
			CreatedAt:   createdAt.Format(common.TimeFormat),
			Installments: []PaymentPlanInstallment{
				{
					ID:       installmentID.String(),
					Amount:   decimalAmount.String(),
					Currency: currency,
					DueAt:    dueAt.Format(common.TimeFormat),
					Status:   status,
				},
			},
		}
	)

	type args struct {
		planID                    uuid.UUID
		completePaymentPlanParams *CompletePaymentPlanParams
	}

	tests := []struct {
		name    string
		prepare func(rm *repomock.MockRepository)
		args    args
		want    *PaymentPlans
		wantErr bool
	}{
		{
			name: "happy path",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(paymentPlans, nil),
					rm.EXPECT().ListPaymentInstallmentsByPlanID(ctx, gomock.Eq(planID)).Return(paymentInstallments, nil),
				)
			},
			args: args{
				planID:                    planID,
				completePaymentPlanParams: paymentPlanParams,
			},
			want:    paymentPlanResponse,
			wantErr: false,
		},
		{
			name: "ListPaymentPlansByUserID error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				planID:                    planID,
				completePaymentPlanParams: paymentPlanParams,
			},
			wantErr: true,
		},
		{
			name: "ListPaymentInstallmentsByPlanID error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return(paymentPlans, nil),
					rm.EXPECT().ListPaymentInstallmentsByPlanID(ctx, gomock.Eq(planID)).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args: args{
				planID:                    planID,
				completePaymentPlanParams: paymentPlanParams,
			},
			wantErr: true,
		},
		{
			name: "PaymentRecordNotFound error",
			prepare: func(rm *repomock.MockRepository) {
				gomock.InOrder(
					rm.EXPECT().ListPaymentPlansByUserID(ctx, gomock.Eq(userID)).Return([]*payments.Plan{}, nil),
				)
			},
			args: args{
				planID:                    planID,
				completePaymentPlanParams: paymentPlanParams,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repomock.NewMockRepository(ctrl)

			if tt.prepare != nil {
				tt.prepare(repo)
			}

			p := &PaymentServiceImp{repository: repo}

			got, err := p.CompletePaymentPlanCreation(ctx, tt.args.planID, tt.args.completePaymentPlanParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentServiceImp.CompletePaymentPlanCreation() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaymentServiceImp.CompletePaymentPlanCreation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentServiceImp_UseRepo(t *testing.T) {
	t.Parallel()

	service := NewPaymentPlanService()
	service.UseRepo(&sqlc.Repo{})

	if reflect.TypeOf(service.repository) != reflect.TypeOf(&sqlc.Repo{}) {
		t.Errorf("returned repository is not of Repo")
	}
}
