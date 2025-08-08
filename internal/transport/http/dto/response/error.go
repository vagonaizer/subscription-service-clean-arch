package response

import "time"

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string            `json:"code" example:"INVALID_INPUT"`
	Message   string            `json:"message" example:"Invalid input provided"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp" example:"2025-01-15T10:30:00Z"`
	RequestID string            `json:"request_id,omitempty" example:"20250115103000-abc123"`
}

type ValidationErrorResponse struct {
	Error            ErrorDetail       `json:"error"`
	ValidationErrors []ValidationError `json:"validation_errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"price"`
	Message string `json:"message" example:"must be greater than 0"`
	Value   string `json:"value,omitempty" example:"-100"`
}

func NewErrorResponse(code, message string, details map[string]string, requestID string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			Details:   details,
			Timestamp: time.Now(),
			RequestID: requestID,
		},
	}
}

func NewValidationErrorResponse(code, message string, validationErrors []ValidationError, requestID string) ValidationErrorResponse {
	return ValidationErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			Timestamp: time.Now(),
			RequestID: requestID,
		},
		ValidationErrors: validationErrors,
	}
}
