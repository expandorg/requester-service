package api

import (
	"net/http"

	"github.com/expandorg/requester-service/pkg/svc-kit/svc"
)

// Error Api struct
type Error struct {
	Private bool
	Status  int
	Err     error
	Body    interface{}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// GetPublicMessage for api
func (e *Error) GetPublicMessage() string {
	if e.Private == true {
		return "Internal server error"
	}
	return e.Err.Error()
}

// WrapErr return http error
func WrapErr(err error) *Error {
	if err, ok := err.(*Error); ok {
		return err
	}
	if err, ok := err.(*svc.ServiceError); ok {
		return &Error{false, err.Status, err, nil}
	}
	return &Error{true, http.StatusInternalServerError, err, nil}
}

// ErrWithBody error with body
func ErrWithBody(err error, body interface{}) *Error {
	if err, ok := err.(*Error); ok {
		return &Error{err.Private, err.Status, err.Err, body}
	}
	if err, ok := err.(*svc.ServiceError); ok {
		return &Error{false, err.Status, err.Err, body}
	}
	return &Error{true, http.StatusInternalServerError, err, body}
}

// NotFound 404
func NotFound(err error) *Error {
	return &Error{false, http.StatusNotFound, err, nil}
}

// BadRequest 400
func BadRequest(err error) *Error {
	return &Error{false, http.StatusBadRequest, err, nil}
}

// Forbidden 403
func Forbidden(err error) *Error {
	return &Error{false, http.StatusForbidden, err, nil}
}

// Unauthorized 401
func Unauthorized(err error) *Error {
	return &Error{false, http.StatusUnauthorized, err, nil}
}
