package apperror

const (
	CodeNotFound             = "NOT_FOUND"
	CodeInvalidInput         = "INVALID_INPUT"
	CodeValidationFailed     = "VALIDATION_FAILED"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeForbidden            = "FORBIDDEN"
	CodeConflict             = "CONFLICT"
	CodeTooManyRequests      = "TOO_MANY_REQUESTS"
	CodeInternalError        = "INTERNAL_ERROR"
	CodeDatabaseError        = "DATABASE_ERROR"
	CodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
)

const (
	CodeSubscriptionNotFound    = "SUBSCRIPTION_NOT_FOUND"
	CodeSubscriptionExists      = "SUBSCRIPTION_EXISTS"
	CodeInvalidSubscriptionData = "INVALID_SUBSCRIPTION_DATA"
	CodeInvalidDateFormat       = "INVALID_DATE_FORMAT"
	CodeInvalidDateRange        = "INVALID_DATE_RANGE"
	CodeInvalidUserID           = "INVALID_USER_ID"
	CodeInvalidPrice            = "INVALID_PRICE"
	CodeInvalidServiceName      = "INVALID_SERVICE_NAME"
	CodeInvalidPaginationParams = "INVALID_PAGINATION_PARAMS"
	CodeInvalidFilterParams     = "INVALID_FILTER_PARAMS"
)

var ErrorMessages = map[string]string{
	CodeNotFound:             "Resource not found",
	CodeInvalidInput:         "Invalid input provided",
	CodeValidationFailed:     "Validation failed",
	CodeUnauthorized:         "Unauthorized access",
	CodeForbidden:            "Access forbidden",
	CodeConflict:             "Resource conflict",
	CodeTooManyRequests:      "Too many requests",
	CodeInternalError:        "Internal server error",
	CodeDatabaseError:        "Database operation failed",
	CodeExternalServiceError: "External service error",
	CodeServiceUnavailable:   "Service temporarily unavailable",

	CodeSubscriptionNotFound:    "Subscription not found",
	CodeSubscriptionExists:      "Subscription already exists",
	CodeInvalidSubscriptionData: "Invalid subscription data",
	CodeInvalidDateFormat:       "Invalid date format, expected MM-YYYY",
	CodeInvalidDateRange:        "Invalid date range",
	CodeInvalidUserID:           "Invalid user ID format",
	CodeInvalidPrice:            "Price must be a positive integer",
	CodeInvalidServiceName:      "Service name cannot be empty",
	CodeInvalidPaginationParams: "Invalid pagination parameters",
	CodeInvalidFilterParams:     "Invalid filter parameters",
}
