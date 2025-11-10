package errors

import (
	"errors"
	"net/http"
)

const (
	ErrorCodeValidation          = "VALIDATION_ERROR"
	ErrorCodeUnauthorized        = "UNAUTHORIZED"
	ErrorCodeRequestTimeout      = "REQUEST_TIMEOUT"
	ErrorCodeDataNotFound        = "DATA_NOT_FOUND"
	ErrorCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrorCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	ErrorCodeNotImplemented      = "NOT_IMPLEMENTED"
	ErrorCodeDataDuplicate       = "DATA_DUPLICATE"
	ErrorCodeConflict            = "CONFLICT"
)

var (
	ErrUnauthorized = &Error{
		Code:       ErrorCodeUnauthorized,
		Message:    "Unauthorized access",
		StatusCode: http.StatusUnauthorized,
	}

	ErrDataNotFound = &Error{
		Code:       ErrorCodeDataNotFound,
		Message:    "Data not found",
		StatusCode: http.StatusNotFound,
	}

	ErrDuplicatedData = &Error{
		Code:       ErrorCodeDataDuplicate,
		Message:    "Duplicated data",
		StatusCode: http.StatusConflict,
	}

	ErrRequestTimeout = &Error{
		Code:       ErrorCodeRequestTimeout,
		Message:    "Request timeout",
		StatusCode: http.StatusRequestTimeout,
	}

	ErrInternalServerError = &Error{
		Code:       ErrorCodeInternalServerError,
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrServiceUnavailable = &Error{
		Code:       ErrorCodeServiceUnavailable,
		Message:    "Service is currently unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}

	ErrNotImplemented = &Error{
		Code:       ErrorCodeNotImplemented,
		Message:    "This functionality is not implemented yet",
		StatusCode: http.StatusNotImplemented,
	}

	ErrConflict = &Error{
		Code:       ErrorCodeConflict,
		Message:    "Resource conflict",
		StatusCode: http.StatusConflict,
	}
)

func NewValidationError(err error) *Error {
	return &Error{
		Code:       ErrorCodeValidation,
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(err error) *Error {
	return &Error{
		Code:       ErrorCodeDataNotFound,
		Message:    err.Error(),
		StatusCode: http.StatusNotFound,
	}
}

func NewDuplicatedDataError(err error) *Error {
	return &Error{
		Code:       ErrorCodeDataDuplicate,
		Message:    err.Error(),
		StatusCode: http.StatusConflict,
	}
}

func NewConflictError(err error) *Error {
	return &Error{
		Code:       ErrorCodeConflict,
		Message:    err.Error(),
		StatusCode: http.StatusConflict,
	}
}

func NewUnauthorizedError(err error) *Error {
	return &Error{
		Code:       ErrorCodeUnauthorized,
		Message:    err.Error(),
		StatusCode: http.StatusUnauthorized,
	}
}

func IsNotFound(err error) bool {
	return IsErrorCode(err, ErrorCodeDataNotFound)
}

func IsErrorCode(err error, code string) bool {
	var domainError *Error
	return errors.As(err, &domainError) && domainError.Code == code
}
