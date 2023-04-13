package errs

import (
	"good/pkg/http"
	nethttp "net/http"
)

// BadRequest BadRequest
func BadRequest(msg string) *http.Error {
	return http.NewError(nethttp.StatusBadRequest, 4400, msg)
}

// ResourceNotFound ResourceNotFound
func ResourceNotFound(msg string) *http.Error {
	return http.NewError(nethttp.StatusNotFound, 4404, msg)
}

// AuthenticationFailed AuthenticationFailed
func AuthenticationFailed(msg string) *http.Error {
	return http.NewError(nethttp.StatusUnauthorized, 4401, msg)
}

// AuthorizationFailed AuthorizationFailed
func AuthorizationFailed(msg string) *http.Error {
	return http.NewError(nethttp.StatusForbidden, 4403, msg)
}

// Conflict Conflict
func Conflict(msg string) *http.Error {
	return http.NewError(nethttp.StatusMethodNotAllowed, 4405, msg)
}

// ValidationFailed ValidationFailed
func ValidationFailed(msg string) *http.Error {
	return http.NewError(nethttp.StatusUnprocessableEntity, 4422, msg)
}

//InternalError InternalError
func InternalError(msg string) *http.Error {
	return http.NewError(nethttp.StatusInternalServerError, 5500, msg)
}

//TooManyRequestsError TooManyRequestsError
func TooManyRequestsError(msg string) *http.Error {
	return http.NewError(nethttp.StatusTooManyRequests, 4429, msg)
}