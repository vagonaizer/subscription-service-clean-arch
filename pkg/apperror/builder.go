package apperror

import "fmt"

func NotFound(resource string) *AppError {
	message := ErrorMessages[CodeNotFound]
	if resource != "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	return New(CodeNotFound, message).WithDetail("resource", resource)
}

func InvalidInput(field, reason string) *AppError {
	return New(CodeInvalidInput, ErrorMessages[CodeInvalidInput]).
		WithDetail("field", field).
		WithDetail("reason", reason)
}

func ValidationFailed(field, reason string) *AppError {
	return New(CodeValidationFailed, ErrorMessages[CodeValidationFailed]).
		WithDetail("field", field).
		WithDetail("reason", reason)
}

func DatabaseError(operation string, cause error) *AppError {
	return Wrap(cause, CodeDatabaseError, ErrorMessages[CodeDatabaseError]).
		WithDetail("operation", operation)
}

func InternalError(message string, cause error) *AppError {
	if message == "" {
		message = ErrorMessages[CodeInternalError]
	}
	return Wrap(cause, CodeInternalError, message)
}

func SubscriptionNotFound(subscriptionID string) *AppError {
	return New(CodeSubscriptionNotFound, ErrorMessages[CodeSubscriptionNotFound]).
		WithDetail("subscription_id", subscriptionID)
}

func InvalidSubscriptionData(field, reason string) *AppError {
	return New(CodeInvalidSubscriptionData, ErrorMessages[CodeInvalidSubscriptionData]).
		WithDetail("field", field).
		WithDetail("reason", reason)
}

func InvalidDateFormat(value string) *AppError {
	return New(CodeInvalidDateFormat, ErrorMessages[CodeInvalidDateFormat]).
		WithDetail("value", value).
		WithDetail("expected_format", "MM-YYYY")
}

func InvalidDateRange(startDate, endDate string) *AppError {
	return New(CodeInvalidDateRange, ErrorMessages[CodeInvalidDateRange]).
		WithDetail("start_date", startDate).
		WithDetail("end_date", endDate)
}

func InvalidUserID(userID string) *AppError {
	return New(CodeInvalidUserID, ErrorMessages[CodeInvalidUserID]).
		WithDetail("user_id", userID)
}

func InvalidPrice(price interface{}) *AppError {
	return New(CodeInvalidPrice, ErrorMessages[CodeInvalidPrice]).
		WithDetail("price", fmt.Sprintf("%v", price))
}

func InvalidServiceName() *AppError {
	return New(CodeInvalidServiceName, ErrorMessages[CodeInvalidServiceName])
}

func InvalidPaginationParams(limit, offset int) *AppError {
	return New(CodeInvalidPaginationParams, ErrorMessages[CodeInvalidPaginationParams]).
		WithDetail("limit", fmt.Sprintf("%d", limit)).
		WithDetail("offset", fmt.Sprintf("%d", offset))
}

func ServiceUnavailable(service string, cause error) *AppError {
	return Wrap(cause, CodeServiceUnavailable, ErrorMessages[CodeServiceUnavailable]).
		WithDetail("service", service)
}

func Conflict(resource, reason string) *AppError {
	return New(CodeConflict, ErrorMessages[CodeConflict]).
		WithDetail("resource", resource).
		WithDetail("reason", reason)
}

func InvalidFilterParams(field, reason string) *AppError {
	return New(CodeInvalidFilterParams, ErrorMessages[CodeInvalidFilterParams]).
		WithDetail("field", field).
		WithDetail("reason", reason)
}
