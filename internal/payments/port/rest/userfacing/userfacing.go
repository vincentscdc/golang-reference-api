// Package userfacing contains all user HTTP APIs
package userfacing

type contextValKey string

const (
	contextValKeyUserUUID contextValKey = "authenticatedUserUUID"
)
