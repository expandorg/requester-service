package svc

import (
	"net/http"
)

// ServiceError business logic error
type ServiceError struct {
	Status int
	Err    error
}

func (e *ServiceError) Error() string {
	return e.Err.Error()
}

// NotFound resource error
func NotFound(err error) *ServiceError {
	return &ServiceError{http.StatusNotFound, err}
}

// ArgumentsErr Arguments resource error
func ArgumentsErr(err error) *ServiceError {
	return &ServiceError{http.StatusBadRequest, err}
}

// Forbidden Access Forbidden error
func Forbidden(err error) *ServiceError {
	return &ServiceError{http.StatusForbidden, err}
}

// General business logic error
func ApplicationError(err error) *ServiceError {
	return &ServiceError{http.StatusBadRequest, err}
}
