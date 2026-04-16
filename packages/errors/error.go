package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     any
}

func (e *APIError) Error() string {
	return fmt.Errorf("API Error %d:%s", e.Code, e.Message).Error()
}

func NotFoundError() APIError {
	return APIError{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
}

func InternalServerError() APIError {
	return APIError{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
}

func NotAcceptableError() APIError {
	return APIError{Code: http.StatusNotAcceptable, Message: http.StatusText(http.StatusNotAcceptable)}
}
