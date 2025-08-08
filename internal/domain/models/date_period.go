package models

import (
	"errors"
	"time"
)

/*
DatePeriod — структура для хранения временного интервала.
Хранит дату начала (from) и дату окончания (to).
Используется для проверки попадания дат в диапазон,
нахождения пересечений и расчёта длительности.
*/
type DatePeriod struct {
	from time.Time
	to   time.Time
}

/** Конструктор для создания нового периода с заданными датами. */
func NewDatePeriod(from, to time.Time) *DatePeriod {
	return &DatePeriod{
		from: from,
		to:   to,
	}
}

/** Геттер/сеттер для даты начала. */
func (dp *DatePeriod) From() time.Time {
	return dp.from
}

func (dp *DatePeriod) SetFrom(from time.Time) {
	dp.from = from
}

/** Геттер/сеттер для даты окончания. */
func (dp *DatePeriod) To() time.Time {
	return dp.to
}

func (dp *DatePeriod) SetTo(to time.Time) {
	dp.to = to
}

/** Проверяет, входит ли переданная дата в диапазон (включительно). */
func (dp *DatePeriod) Contains(date time.Time) bool {
	return !date.Before(dp.from) && !date.After(dp.to)
}

/** Проверяет, пересекается ли этот период с другим. */
func (dp *DatePeriod) Overlaps(other DatePeriod) bool {
	return dp.from.Before(other.to) && dp.to.After(other.from)
}

/** Возвращает длительность периода как time.Duration. */
func (dp *DatePeriod) Duration() time.Duration {
	return dp.to.Sub(dp.from)
}

/** Проверяет, что дата окончания не раньше даты начала. */
func (dp *DatePeriod) Validate() error {
	if dp.to.Before(dp.from) {
		return errors.New("end date cannot be before start date")
	}
	return nil
}
