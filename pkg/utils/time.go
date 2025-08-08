package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
)

const DateLayout = "01-2006"

func ParseMonthYear(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, apperror.InvalidDateFormat(dateStr)
	}

	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return time.Time{}, apperror.InvalidDateFormat(dateStr)
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, apperror.InvalidDateFormat(dateStr)
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil || year < 2000 || year > 2100 {
		return time.Time{}, apperror.InvalidDateFormat(dateStr)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

func FormatMonthYear(t time.Time) string {
	return t.Format(DateLayout)
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Nanosecond)
}

func ValidateDateRange(startDate, endDate *time.Time) error {
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return apperror.InvalidDateRange(
			FormatMonthYear(*startDate),
			FormatMonthYear(*endDate),
		)
	}
	return nil
}

func ParseDateRange(startDateStr, endDateStr string) (*time.Time, *time.Time, error) {
	var startDate, endDate *time.Time

	if startDateStr != "" {
		start, err := ParseMonthYear(startDateStr)
		if err != nil {
			return nil, nil, err
		}
		start = StartOfMonth(start)
		startDate = &start
	}

	if endDateStr != "" {
		end, err := ParseMonthYear(endDateStr)
		if err != nil {
			return nil, nil, err
		}
		end = EndOfMonth(end)
		endDate = &end
	}

	if err := ValidateDateRange(startDate, endDate); err != nil {
		return nil, nil, err
	}

	return startDate, endDate, nil
}

func MonthsDifference(start, end time.Time) int {
	startMonth := start.Year()*12 + int(start.Month()) - 1
	endMonth := end.Year()*12 + int(end.Month()) - 1
	diff := endMonth - startMonth + 1
	if diff < 0 {
		return 0
	}
	return diff
}
