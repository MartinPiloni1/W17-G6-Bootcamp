package httperrors

import (
	"errors"
)

type BadRequestError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type ConflictError struct {
	Message string
}

type UnprocessableEntityError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e ConflictError) Error() string {
	return e.Message
}

func (e UnprocessableEntityError) Error() string {
	return e.Message
}

func (e InternalServerError) Error() string {
	return e.Message
}

// GetErrorData inspects the given error and maps it to an HTTP status code
// and a client‐safe error message. It recognizes the following error types:
//
// • BadRequestError           → 400 Bad Request
// • NotFoundError             → 404 Not Found
// • ConflictError             → 409 Conflict
// • UnprocessableEntityError  → 422 Unprocessable Entity
// • InternalServerError       → 500 Internal Server Error
//
// If the error is none of the above, it returns 500 with a generic message.
func GetErrorData(err error) (int, string) {
	switch {
	case errors.As(err, &BadRequestError{}):
		return 400, err.Error()
	case errors.As(err, &NotFoundError{}):
		return 404, err.Error()
	case errors.As(err, &ConflictError{}):
		return 409, err.Error()
	case errors.As(err, &UnprocessableEntityError{}):
		return 422, err.Error()
	case errors.As(err, &InternalServerError{}):
		return 500, err.Error()
	default:
		return 500, "Internal Server Error"
	}
}