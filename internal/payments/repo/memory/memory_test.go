package memory

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/payments"
	"golangreferenceapi/internal/payments/common"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestNewInMemRepository(t *testing.T) {
	t.Parallel()

	if repo := NewInMemRepository(); reflect.TypeOf(repo) != reflect.TypeOf(&InMemRepo{}) {
		t.Error("returned repo is not of Repo")
	}
}

func TestMemoryError_ErrGenerateUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ErrGenerateUUID,
			expectedString: "failed to generate uuid",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestMemoryError_ErrRecordNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ErrRecordNotFound,
			expectedString: "no records found",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestMemoryError_ErrMapTypeAssertion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ErrMapTypeAssertion,
			expectedString: "type assertion failed when load map",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestInMemRepository_CreatePaymentPlan(t *testing.T) {
	t.Parallel()

	var (
		userUUID, _    = uuid.NewV4()
		newUserUUID, _ = uuid.NewV4()
		repo           = NewInMemRepository()
	)

	_, err := repo.CreatePaymentPlan(context.Background(), &payments.CreatePlanParams{
		UserID:   userUUID,
		Currency: "usdc",
		Amount:   *decimal.New(1098, 2),
		Status:   "pending",
	})
	if err != nil {
		t.Fatalf("fail to create payment plan: %v", err)
	}

	type args struct {
		arg *payments.CreatePlanParams
	}

	tests := []struct {
		name    string
		args    args
		want    *payments.Plan
		wantErr bool
	}{
		{
			name: "happy path - user has no existing plans",
			args: args{
				arg: &payments.CreatePlanParams{
					UserID:   newUserUUID,
					Currency: "usdc",
					Amount:   *decimal.New(1098, 2),
					Status:   "pending",
				},
			},
			want: &payments.Plan{
				UserID:   newUserUUID,
				Currency: "usdc",
				Amount:   *decimal.New(1098, 2),
				Status:   "pending",
			},
			wantErr: false,
		},
		{
			name: "happy path - user has existing plans",
			args: args{
				arg: &payments.CreatePlanParams{
					UserID:   userUUID,
					Currency: "usdc",
					Amount:   *decimal.New(1098, 2),
					Status:   "pending",
				},
			},
			want: &payments.Plan{
				UserID:   userUUID,
				Currency: "usdc",
				Amount:   *decimal.New(1098, 2),
				Status:   "pending",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := repo.CreatePaymentPlan(context.Background(), tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got.ID == uuid.Nil {
				t.Errorf("expect uuid but nil returned")
			}

			if got.UserID != tt.want.UserID {
				t.Errorf("wrong expected user id: got %v, want %v", tt.want.UserID, got.UserID)
			}

			if got.Currency != tt.want.Currency {
				t.Errorf("wrong expected currency: got %v, want %v", tt.want.Currency, got.Currency)
			}

			if got.Amount.Cmp(&tt.want.Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", tt.want.Amount, got.Amount)
			}

			if got.Status != tt.want.Status {
				t.Errorf("wrong expected status: got %v, want %v", tt.want.Status, got.Status)
			}
		})
	}
}

func TestInMemRepository_CreatePaymentPlan_Concurrency(t *testing.T) {
	t.Parallel()

	var (
		repo        = NewInMemRepository()
		n           = 20
		userUUID, _ = uuid.NewV4()
		params      = &payments.CreatePlanParams{
			UserID:   userUUID,
			Currency: "usdc",
			Amount:   *decimal.New(1098, 2),
			Status:   "pending",
		}
		want = &payments.Plan{
			UserID:   userUUID,
			Currency: "usdc",
			Amount:   *decimal.New(1098, 2),
			Status:   "pending",
		}
	)

	errs := make(chan error)
	results := make(chan *payments.Plan)

	for i := 0; i < n; i++ {
		go func() {
			res, err := repo.CreatePaymentPlan(context.Background(), params)

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Errorf("got error = %v", err)
		}

		res := <-results

		if res.ID == uuid.Nil {
			t.Errorf("expect uuid but nil returned")
		}

		if res.UserID != want.UserID {
			t.Errorf("wrong expected user id: got %v, want %v", want.UserID, res.UserID)
		}

		if res.Currency != want.Currency {
			t.Errorf("wrong expected currency: got %v, want %v", want.Currency, res.Currency)
		}

		if res.Amount.Cmp(&want.Amount) != 0 {
			t.Errorf("wrong expected amount: got %v, want %v", want.Amount, res.Amount)
		}

		if res.Status != want.Status {
			t.Errorf("wrong expected status: got %v, want %v", want.Status, res.Status)
		}
	}

	plans, err := repo.ListPaymentPlansByUserID(context.Background(), userUUID)
	if err != nil {
		t.Errorf("got error = %v", err)
	}

	if len(plans) != n {
		t.Errorf("got len(plans) = %v, want %v", len(plans), n)
	}
}

func TestInMemRepository_ListPaymentPlansByUserID(t *testing.T) {
	t.Parallel()

	var (
		repo        = NewInMemRepository()
		userUUID, _ = uuid.NewV4()
	)

	existingPlan, err := repo.CreatePaymentPlan(context.Background(), &payments.CreatePlanParams{
		UserID:   userUUID,
		Currency: "usdc",
		Amount:   *decimal.New(1098, 2),
		Status:   "pending",
	})
	if err != nil {
		t.Fatalf("fail to create payment plan: %v", err)
	}

	type args struct {
		userID uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		want    []*payments.Plan
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				userID: userUUID,
			},
			want:    []*payments.Plan{existingPlan},
			wantErr: false,
		},
		{
			name: "no records",
			args: args{
				userID: uuid.Must(uuid.NewV4()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := repo.ListPaymentPlansByUserID(context.Background(), tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemRepository_CreatePaymentInstallment(t *testing.T) {
	t.Parallel()

	var (
		repo         = NewInMemRepository()
		dueAt, _     = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		planID, _    = uuid.NewV4()
		newPlanID, _ = uuid.NewV4()
	)

	_, err := repo.CreatePaymentInstallment(context.Background(), &payments.CreateInstallmentParams{
		PaymentPlanID: planID,
		Currency:      "usdc",
		Amount:        *decimal.New(1098, 2),
		DueAt:         dueAt,
		Status:        "pending",
	})
	if err != nil {
		t.Fatalf("fail to create installment: %v", err)
	}

	type args struct {
		arg *payments.CreateInstallmentParams
	}

	tests := []struct {
		name    string
		args    args
		want    *payments.Installment
		wantErr bool
	}{
		{
			name: "happy path - plan has no existing installments",
			args: args{
				arg: &payments.CreateInstallmentParams{
					PaymentPlanID: newPlanID,
					Currency:      "usdc",
					Amount:        *decimal.New(1098, 2),
					DueAt:         dueAt,
					Status:        "pending",
				},
			},
			want: &payments.Installment{
				PaymentPlanID: newPlanID,
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         dueAt,
				Status:        "pending",
			},
			wantErr: false,
		},
		{
			name: "happy path - plan has existing installments",
			args: args{
				arg: &payments.CreateInstallmentParams{
					PaymentPlanID: planID,
					Currency:      "usdc",
					Amount:        *decimal.New(1098, 2),
					DueAt:         dueAt,
					Status:        "pending",
				},
			},
			want: &payments.Installment{
				PaymentPlanID: planID,
				Currency:      "usdc",
				Amount:        *decimal.New(1098, 2),
				DueAt:         dueAt,
				Status:        "pending",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := repo.CreatePaymentInstallment(context.Background(), tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got.ID == uuid.Nil {
				t.Errorf("expect uuid but nil returned")
			}

			if got.PaymentPlanID != tt.want.PaymentPlanID {
				t.Errorf("wrong expected user id: got %v, want %v", tt.want.PaymentPlanID, got.PaymentPlanID)
			}

			if got.Currency != tt.want.Currency {
				t.Errorf("wrong expected currency: got %v, want %v", tt.want.Currency, got.Currency)
			}

			if got.Amount.Cmp(&tt.want.Amount) != 0 {
				t.Errorf("wrong expected amount: got %v, want %v", tt.want.Amount, got.Amount)
			}

			if !got.DueAt.Equal(tt.want.DueAt) {
				t.Errorf("wrong expected due at: got %v, want %v", tt.want.DueAt, got.DueAt)
			}

			if got.Status != tt.want.Status {
				t.Errorf("wrong expected status: got %v, want %v", tt.want.Status, got.Status)
			}
		})
	}
}

func TestInMemRepository_CreatePaymentInstallment_Concurrency(t *testing.T) {
	t.Parallel()

	var (
		repo      = NewInMemRepository()
		n         = 20
		ctx       = context.Background()
		planID, _ = uuid.NewV4()
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		params    = &payments.CreateInstallmentParams{
			PaymentPlanID: planID,
			Currency:      "usdc",
			Amount:        *decimal.New(1098, 2),
			DueAt:         dueAt,
			Status:        "pending",
		}
		want = &payments.Installment{
			PaymentPlanID: planID,
			Currency:      "usdc",
			Amount:        *decimal.New(1098, 2),
			DueAt:         dueAt,
			Status:        "pending",
		}
	)

	errs := make(chan error)
	results := make(chan *payments.Installment)

	for i := 0; i < n; i++ {
		go func() {
			res, err := repo.CreatePaymentInstallment(ctx, params)

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Errorf("got error = %v", err)
		}

		res := <-results

		if res.ID == uuid.Nil {
			t.Errorf("expect uuid but nil returned")
		}

		if res.PaymentPlanID != want.PaymentPlanID {
			t.Errorf("wrong expected user id: got %v, want %v", want.PaymentPlanID, res.PaymentPlanID)
		}

		if res.Currency != want.Currency {
			t.Errorf("wrong expected currency: got %v, want %v", want.Currency, res.Currency)
		}

		if res.Amount.Cmp(&want.Amount) != 0 {
			t.Errorf("wrong expected amount: got %v, want %v", want.Amount, res.Amount)
		}

		if !res.DueAt.Equal(want.DueAt) {
			t.Errorf("wrong expected due at: got %v, want %v", want.DueAt, res.DueAt)
		}

		if res.Status != want.Status {
			t.Errorf("wrong expected status: got %v, want %v", want.Status, res.Status)
		}
	}

	installments, err := repo.ListPaymentInstallmentsByPlanID(ctx, planID)
	if err != nil {
		t.Errorf("got error = %v", err)
	}

	if len(installments) != n {
		t.Errorf("got len(plans) = %v, want %v", len(installments), n)
	}
}

func TestInMemRepository_ListPaymentInstallmentsByPlanID(t *testing.T) {
	t.Parallel()

	var (
		repo      = NewInMemRepository()
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
		planID, _ = uuid.NewV4()
	)

	installment, err := repo.CreatePaymentInstallment(context.Background(), &payments.CreateInstallmentParams{
		PaymentPlanID: planID,
		Currency:      "usdc",
		Amount:        *decimal.New(1098, 2),
		DueAt:         dueAt,
		Status:        "pending",
	})
	if err != nil {
		t.Fatalf("fail to create installment: %v", err)
	}

	type args struct {
		planID uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		want    []*payments.Installment
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				planID: planID,
			},
			want: []*payments.Installment{
				installment,
			},
			wantErr: false,
		},
		{
			name: "no records",
			args: args{
				planID: uuid.Must(uuid.NewV4()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := repo.ListPaymentInstallmentsByPlanID(context.Background(), tt.args.planID)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCreatePaymentPlan(b *testing.B) {
	repo := NewInMemRepository()
	param := &payments.CreatePlanParams{
		UserID:   uuid.Must(uuid.NewV4()),
		Currency: "usdc",
		Amount:   *decimal.New(1098, 2),
		Status:   "pending",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.CreatePaymentPlan(context.Background(), param)
		if err != nil {
			b.Fatalf("err creating payment plan")
		}
	}
}

func BenchmarkListPaymentPlansByUserID(b *testing.B) {
	repo := NewInMemRepository()
	userID, _ := uuid.NewV4()

	repo.CreatePaymentPlan(context.Background(), &payments.CreatePlanParams{
		UserID:   userID,
		Currency: "usdc",
		Amount:   *decimal.New(1098, 2),
		Status:   "pending",
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ListPaymentPlansByUserID(context.Background(), userID)
		if err != nil {
			b.Fatalf("err listing payment plans by user id")
		}
	}
}

func BenchmarkCreatePaymentInstallment(b *testing.B) {
	repo := NewInMemRepository()
	params := &payments.CreateInstallmentParams{
		PaymentPlanID: uuid.Must(uuid.NewV4()),
		Currency:      "usdc",
		Amount:        *decimal.New(1098, 2),
		DueAt:         time.Time{},
		Status:        "pending",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.CreatePaymentInstallment(context.Background(), params)
		if err != nil {
			b.Fatalf("err creating payment installment")
		}
	}
}

func BenchmarkListPaymentInstallmentsByPlanID(b *testing.B) {
	repo := NewInMemRepository()
	planID, _ := uuid.NewV4()

	repo.CreatePaymentInstallment(context.Background(), &payments.CreateInstallmentParams{
		PaymentPlanID: planID,
		Currency:      "usdc",
		Amount:        *decimal.New(1098, 2),
		DueAt:         time.Time{},
		Status:        "pending",
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.ListPaymentInstallmentsByPlanID(context.Background(), planID)
		if err != nil {
			b.Fatalf("err listing payment installments by plan id")
		}
	}
}
