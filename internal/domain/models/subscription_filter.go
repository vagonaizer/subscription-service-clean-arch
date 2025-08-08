package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/*
*
SubscriptionFilter — вспомогательная структура для фильтрации подписок
по разным критериям. Все поля указатели, чтобы можно было легко отличить
"не задано" от "задано пустым значением".
*/
type SubscriptionFilter struct {
	userID      *uuid.UUID
	serviceName *string
	startDate   *time.Time
	endDate     *time.Time
	isActive    *bool
}

/** Создаёт пустой фильтр без условий. */
func NewSubscriptionFilter() *SubscriptionFilter {
	return &SubscriptionFilter{}
}

/** Геттер/сеттер для фильтра по userID. */
func (f *SubscriptionFilter) UserID() *uuid.UUID {
	return f.userID
}

func (f *SubscriptionFilter) SetUserID(userID *uuid.UUID) {
	f.userID = userID
}

/** Геттер/сеттер для фильтра по названию сервиса. */
func (f *SubscriptionFilter) ServiceName() *string {
	return f.serviceName
}

func (f *SubscriptionFilter) SetServiceName(serviceName *string) {
	f.serviceName = serviceName
}

/** Геттер/сеттер для фильтра по дате начала. */
func (f *SubscriptionFilter) StartDate() *time.Time {
	return f.startDate
}

func (f *SubscriptionFilter) SetStartDate(startDate *time.Time) {
	f.startDate = startDate
}

/** Геттер/сеттер для фильтра по дате окончания. */
func (f *SubscriptionFilter) EndDate() *time.Time {
	return f.endDate
}

func (f *SubscriptionFilter) SetEndDate(endDate *time.Time) {
	f.endDate = endDate
}

/** Геттер/сеттер для фильтра по статусу активности подписки. */
func (f *SubscriptionFilter) IsActive() *bool {
	return f.isActive
}

func (f *SubscriptionFilter) SetIsActive(isActive *bool) {
	f.isActive = isActive
}

/** Проверки, задано ли конкретное поле в фильтре. */
func (f *SubscriptionFilter) HasUserID() bool {
	return f.userID != nil
}

func (f *SubscriptionFilter) HasServiceName() bool {
	return f.serviceName != nil && *f.serviceName != ""
}

func (f *SubscriptionFilter) HasDateRange() bool {
	return f.startDate != nil || f.endDate != nil
}

/*
*
Validate — проверяет, что диапазон дат корректный.
Например, дата окончания не может быть раньше даты начала.
*/
func (f *SubscriptionFilter) Validate() error {
	if f.startDate != nil && f.endDate != nil && f.endDate.Before(*f.startDate) {
		return errors.New("end date cannot be before start date")
	}
	return nil
}
