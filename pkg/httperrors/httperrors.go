package httperrors

import (
	"errors"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

type NotFoundError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type ConflictError struct {
	Message string
}

type UnprocessableEntityError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e ConflictError) Error() string {
	return e.Message
}

func (e UnprocessableEntityError) Error() string {
	return e.Message
}

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
	default:
		return 500, "Internal Server Error"
	}
}
