package internalfacing

import (
	"context"

	"github.com/google/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

const (
	urlParamUserUUID        = "user_uuid"
	urlParamPaymentPlanUUID = "uuid"
	urlParamPaymentUUID     = "payment_uuid"
	urlParamInstallmentID   = "installments_id"
)

type PaymentPlanParam struct {
	UUID     *uuid.UUID `json:"uuid"`
	UserUUID *uuid.UUID `json:"user_uuid"`
}

func parsePaymentPlanParam(
	ctx context.Context, paramsGetter handlerwrap.NamedURLParamsGetter,
) (*PaymentPlanParam, *handlerwrap.ErrorResponse) {
	// user uuid
	userUUID, err := parseUUIDFormatParam(ctx, paramsGetter, urlParamUserUUID)
	if err != nil {
		return nil, err
	}

	// payment plan uuid
	paymentPlanUUID, err := parseUUIDFormatParam(ctx, paramsGetter, urlParamPaymentPlanUUID)
	if err != nil {
		return nil, err
	}

	return &PaymentPlanParam{
		UUID:     paymentPlanUUID,
		UserUUID: userUUID,
	}, nil
}

func parseUUIDFormatParam(
	ctx context.Context, paramsGetter handlerwrap.NamedURLParamsGetter, name string,
) (*uuid.UUID, *handlerwrap.ErrorResponse) {
	val, err := paramsGetter(ctx, name)
	if err != nil {
		return nil, err
	}

	uuidVal, parseErr := uuid.Parse(val)
	if parseErr != nil {
		return nil, handlerwrap.ParsingParamError{
			Name:  name,
			Value: val,
		}.ToErrorResponse()
	}

	return &uuidVal, nil
}
