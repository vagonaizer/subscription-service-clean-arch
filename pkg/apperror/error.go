package apperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	code       string
	message    string
	details    map[string]string
	cause      error
	httpStatus int
}

func New(code, message string) *AppError {
	return &AppError{
		code:       code,
		message:    message,
		details:    make(map[string]string),
		httpStatus: getDefaultHTTPStatus(code),
	}
}

func Wrap(err error, code, message string) *AppError {
	return &AppError{
		code:       code,
		message:    message,
		details:    make(map[string]string),
		cause:      err,
		httpStatus: getDefaultHTTPStatus(code),
	}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *AppError) Code() string {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Details() map[string]string {
	return e.details
}

func (e *AppError) Cause() error {
	return e.cause
}

func (e *AppError) HTTPStatus() int {
	return e.httpStatus
}

func (e *AppError) WithDetail(key, value string) *AppError {
	e.details[key] = value
	return e
}

func (e *AppError) WithDetails(details map[string]string) *AppError {
	for k, v := range details {
		e.details[k] = v
	}
	return e
}

func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.httpStatus = status
	return e
}

func (e *AppError) WithCause(cause error) *AppError {
	e.cause = cause
	return e
}

func (e *AppError) Is(target error) bool {
	if ae, ok := target.(*AppError); ok {
		return e.code == ae.code
	}
	return false
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) ToJSON() []byte {
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.code,
			"message": e.message,
		},
	}

	if len(e.details) > 0 {
		response["error"].(map[string]interface{})["details"] = e.details
	}

	data, _ := json.Marshal(response)
	return data
}

func (e *AppError) Clone() *AppError {
	details := make(map[string]string)
	for k, v := range e.details {
		details[k] = v
	}

	return &AppError{
		code:       e.code,
		message:    e.message,
		details:    details,
		cause:      e.cause,
		httpStatus: e.httpStatus,
	}
}

func IsAppError(err error) (*AppError, bool) {
	if ae, ok := err.(*AppError); ok {
		return ae, true
	}
	return nil, false
}

func getDefaultHTTPStatus(code string) int {
	switch code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInvalidInput, CodeValidationFailed:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeConflict:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInternalError, CodeDatabaseError, CodeExternalServiceError:
		return http.StatusInternalServerError
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
