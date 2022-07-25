package rest

import (
	"errors"
	"net/http"

	"golangreferenceapi/internal/payments/service"

	"github.com/monacohq/golang-common/transport/http/handlerwrap/v3"
)

func ServiceErrorToErrorResp(err error) *handlerwrap.ErrorResponse {
	switch {
	case errors.As(err, &service.CreatePaymentPlanError{}):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusInternalServerError,
			"create_payment_plan_failed",
			"create payment plan failed",
		)
	case errors.As(err, &service.ListPaymentPlansByUserIDError{}):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusInternalServerError,
			"list_payment_plan_by_userid_failed",
			"list payment plan by userid failed",
		)
	case errors.As(err, &service.CreatePaymentInstallmentError{}):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusInternalServerError,
			"create_payment_installment_failed",
			"create payment installment failed",
		)
	case errors.As(err, &service.ListPaymentInstallmentsByPlanIDError{}):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusInternalServerError,
			"list_payment_installments_by_planid_failed",
			"list payment installments by planid failed",
		)
	case errors.As(err, &service.PaymentRecordNotFoundError{}):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusNotFound,
			"payment_record_not_found",
			"payment record not found",
		)
	default:
		return handlerwrap.InternalServerError{Err: err}.ToErrorResponse()
	}
}
