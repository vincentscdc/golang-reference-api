package userfacing

import (
	"net/http"
	"time"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

const (
	paymentPlansDefaultLimit = 10
)

// PaymentPlanResponse represents a specific payment plan
type PaymentPlanResponse struct {
	Description                 string            `json:"description"`
	Currency                    string            `json:"currency"`
	TotalAmount                 string            `json:"total_amount"`
	PayableAmount               string            `json:"payable_amount"`
	TotalPaidAmount             string            `json:"total_paid_amount"`
	TotalRefundAmount           string            `json:"total_refund_amount"`
	TotalLateChargeAmount       string            `json:"total_late_charge_amount"`
	OutstandingLateChargeAmount string            `json:"outstanding_late_charge_amount"`
	Status                      string            `json:"status"`
	RepaymentStatusForDisplay   string            `json:"repayment_status_for_display"`
	IsLiquidated                bool              `json:"is_liquidated"`
	FromCurrency                string            `json:"from_currency"`
	NextRepaymentID             string            `json:"next_repayment_id"`
	CreatedAt                   time.Time         `json:"created_at"`
	Payments                    []PaymentResponse `json:"payments"`
}

// PaymentResponse represents a specific payment
type PaymentResponse struct {
	ID                string    `json:"id"`
	FromAmount        string    `json:"from_amount"`
	FromCurrency      string    `json:"from_currency"`
	Amount            string    `json:"amount"`
	LateChargeAmount  string    `json:"late_charge_amount"`
	RefundAmount      string    `json:"refund_amount"`
	OutstandingAmount string    `json:"outstanding_amount"`
	Currency          string    `json:"currency"`
	DueAt             time.Time `json:"due_at"`
	SettledAt         time.Time `json:"settled_at"`
	Status            string    `json:"status"`
}

func listPaymentPlansHandler() handlerwrap.TypedHandler {
	return func(r *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		// user uuid
		_ = getUserUUID(r.Context())
		// pagination params
		_, err := parsePaginationURLQuery(r.URL, paymentPlansDefaultLimit, paymentPlansCreatedAtOrderDESC)
		if err != nil {
			return nil, err
		}

		var plans []PaymentPlanResponse

		resp := &handlerwrap.Response{Body: plans, HTTPStatusCode: http.StatusOK}

		return resp, nil
	}
}
