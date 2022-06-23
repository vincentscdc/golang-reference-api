package internalfacing

import (
	"fmt"
	"time"

	"golangreferenceapi/internal/payments/common"

	"golangreferenceapi/internal/payments/service"

	"github.com/google/uuid"
)

func NewCreatePaymentPlanParams(request *CreatePendingPaymentPlanRequest) (*service.CreatePaymentPlanParams, error) {
	paymentID, err := uuid.Parse(request.PendingPayment.ID)
	if err != nil {
		return nil, inputParseError{fieldName: "id", err: err}
	}

	params := &service.CreatePaymentPlanParams{
		ID:           paymentID,
		Currency:     request.PendingPayment.Currency,
		TotalAmount:  request.PendingPayment.TotalAmount,
		Installments: make([]service.PaymentPlanInstallmentParams, 0, len(request.PendingPayment.Installments)),
	}

	for _, inst := range request.PendingPayment.Installments {
		instID, err := uuid.Parse(inst.ID)
		if err != nil {
			return nil, inputParseError{fieldName: "installments", err: err}
		}

		dueAt, err := time.Parse(common.TimeFormat, inst.DueAt)
		if err != nil {
			return nil, inputParseError{fieldName: "due_at", err: err}
		}

		instParam := service.PaymentPlanInstallmentParams{
			ID:       instID,
			Amount:   inst.Amount,
			Currency: inst.Currency,
			DueAt:    dueAt,
		}

		params.Installments = append(params.Installments, instParam)
	}

	return params, nil
}

type inputParseError struct {
	fieldName string
	err       error
}

func (i inputParseError) Error() string {
	return fmt.Sprintf("%v, err: %v", i.fieldName, i.err.Error())
}
