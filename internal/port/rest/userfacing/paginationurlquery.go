package userfacing

import (
	"net/url"
	"strconv"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

const (
	offsetKey            = "offset"
	limitKey             = "limit"
	createAtOrderKey     = "created_at_order"
	paginationIntBase    = 10
	paginationIntBitSize = 64

	paymentPlansCreatedAtOrderASC  = "asc"
	paymentPlansCreatedAtOrderDESC = "desc"
)

type PaginationURLQuery struct {
	Offset         int64  `json:"offset"`
	Limit          int64  `json:"limit"`
	CreatedAtOrder string `json:"created_at_order"`
}

func parsePaginationURLQuery(
	u *url.URL, defaultLimit int64,
	defaultCreatedAtOrder string,
) (*PaginationURLQuery, *handlerwrap.ErrorResponse) {
	var (
		queryValues    = u.Query()
		offset         int64
		limit          int64
		createdAtOrder string
		err            *handlerwrap.ErrorResponse
	)

	offset, err = parsePaginationOffset(queryValues)
	if err != nil {
		return nil, err
	}

	limit, err = parsePaginationLimit(queryValues, defaultLimit)
	if err != nil {
		return nil, err
	}

	createdAtOrder, err = parsePaginationCreatedAtOrder(queryValues, defaultCreatedAtOrder)
	if err != nil {
		return nil, err
	}

	return &PaginationURLQuery{
		Offset:         offset,
		Limit:          limit,
		CreatedAtOrder: createdAtOrder,
	}, nil
}

func parsePaginationOffset(values url.Values) (int64, *handlerwrap.ErrorResponse) {
	var (
		offset int64
		err    error
	)

	offsetVal := values.Get(offsetKey)
	if offsetVal != "" {
		offset, err = strconv.ParseInt(offsetVal, paginationIntBase, paginationIntBitSize)
		if err != nil {
			return 0, handlerwrap.ParsingParamError{
				Name:  limitKey,
				Value: offsetVal,
			}.ToErrorResponse()
		}
	}

	if offset < 0 {
		offset = 0
	}

	return offset, nil
}

func parsePaginationLimit(values url.Values, defaultLimit int64) (int64, *handlerwrap.ErrorResponse) {
	var (
		limit int64
		err   error
	)

	limitVal := values.Get(limitKey)
	if limitVal != "" {
		limit, err = strconv.ParseInt(limitVal, paginationIntBase, paginationIntBitSize)
		if err != nil {
			return 0, handlerwrap.ParsingParamError{
				Name:  limitKey,
				Value: limitVal,
			}.ToErrorResponse()
		}
	}

	if limit <= 0 || limit > defaultLimit {
		limit = defaultLimit
	}

	return limit, nil
}

func parsePaginationCreatedAtOrder(
	values url.Values, defaultCreatedAtOrder string,
) (string, *handlerwrap.ErrorResponse) {
	createdAtOrderVal := values.Get(createAtOrderKey)
	if createdAtOrderVal == paymentPlansCreatedAtOrderASC || createdAtOrderVal == paymentPlansCreatedAtOrderDESC {
		return createdAtOrderVal, nil
	}

	if createdAtOrderVal != "" {
		return "", handlerwrap.ParsingParamError{
			Name:  createAtOrderKey,
			Value: createdAtOrderVal,
		}.ToErrorResponse()
	}

	return defaultCreatedAtOrder, nil
}
