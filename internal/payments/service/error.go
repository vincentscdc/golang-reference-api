package service

import "errors"

var (
	ErrGenerateUUID   = errors.New("failed to generate uuid")
	ErrRecordNotFound = errors.New("record not fund")
)
