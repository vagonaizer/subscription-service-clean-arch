package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
)

func ValidateUUID(id string, fieldName string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperror.InvalidInput(fieldName, "cannot be empty")
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.InvalidUserID(id)
	}

	return parsedUUID, nil
}

func ValidateServiceName(serviceName string) error {
	if strings.TrimSpace(serviceName) == "" {
		return apperror.InvalidServiceName()
	}
	if len(serviceName) > 255 {
		return apperror.InvalidInput("service_name", "must not exceed 255 characters")
	}
	return nil
}

func ValidatePrice(price int) error {
	if price <= 0 {
		return apperror.InvalidPrice(price)
	}
	if price > 1000000 {
		return apperror.InvalidInput("price", "must not exceed 1,000,000")
	}
	return nil
}

func ValidatePagination(limit, offset int) (int, int, error) {
	if limit < 0 {
		return 0, 0, apperror.InvalidPaginationParams(limit, offset).
			WithDetail("limit_error", "must be non-negative")
	}
	if offset < 0 {
		return 0, 0, apperror.InvalidPaginationParams(limit, offset).
			WithDetail("offset_error", "must be non-negative")
	}

	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return limit, offset, nil
}

func NormalizeString(s string) string {
	return strings.TrimSpace(s)
}

func IsEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}
