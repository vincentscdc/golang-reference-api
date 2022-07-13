package internalfacing

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
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

	uuidVal, parseErr := uuid.FromString(val)
	if parseErr != nil {
		return nil, handlerwrap.ParsingParamError{
			Name:  name,
			Value: val,
		}.ToErrorResponse()
	}

	return &uuidVal, nil
}
