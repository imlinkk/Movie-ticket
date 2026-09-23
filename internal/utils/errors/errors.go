package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden resource")
	ErrBadRequest        = errors.New("bad request")
	ErrConflict          = errors.New("resource already exists")
	ErrInternalServer    = errors.New("internal server error")
	ErrSeatAlreadyBooked = errors.New("seat has already been booked for this showtime")
	ErrNoAvailableSeats  = errors.New("no available seats left for this showtime")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NewAppError(statusCode int, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

func Conflict(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, err)
}

func Internal(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}
