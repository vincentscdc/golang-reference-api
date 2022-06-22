package repo

import "fmt"

type decimalParseError struct {
	value      interface{}
	WrappedErr error
}

func (e decimalParseError) Error() string {
	return fmt.Sprintf("failed to set value %v as decimal value", e.value)
}

type pgxDBQueryRunError struct {
	WrappedErr error
}

func (e pgxDBQueryRunError) Error() string {
	return "failed to run query"
}

type uuidParseError struct {
	WrappedErr error
	value      interface{}
}

func (e uuidParseError) Error() string {
	return fmt.Sprintf("failed to parse value %v as UUID", e.value)
}

type unsupportedDBEntityError struct {
	WrappedErr error
}

func (e unsupportedDBEntityError) Error() string {
	return "DB entity interface does not match any supported struct"
}
