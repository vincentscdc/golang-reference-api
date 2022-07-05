package internalfacing

import (
	"context"

	"github.com/google/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

const (
	urlParamUserUUID    = "user_uuid"
	urlParamPaymentUUID = "payment_uuid"
)

type PaymentPlanParam struct {
	UUID     *uuid.UUID `json:"uuid"`
	UserUUID *uuid.UUID `json:"user_uuid"`
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
