package sqlc

type UnsupportedDBEntityError struct{}

func (e UnsupportedDBEntityError) Error() string {
	return "DB entity interface does not match any supported struct"
}
