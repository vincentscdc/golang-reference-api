package repo

type pgxDBQueryRunError struct {
	WrappedErr error
}

func (e pgxDBQueryRunError) Error() string {
	return "failed to run query"
}

type unsupportedDBEntityError struct {
	WrappedErr error
}

func (e unsupportedDBEntityError) Error() string {
	return "DB entity interface does not match any supported struct"
}
