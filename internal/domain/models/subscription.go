package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/*
Subscription описывает подписку пользователя на какой-то сервис.
Я специально сделал поля неэкспортируемыми, чтобы управлять ими
только через методы (инкапсуляция и контроль изменений).
*/
type Subscription struct {
	id          uuid.UUID
	serviceName string
	price       int
	userID      uuid.UUID
	startDate   time.Time
	endDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

/*
*
NewSubscription создаёт новую подписку с текущим временем как createdAt/updatedAt.
ID генерируется автоматически, чтобы не зависеть от внешнего кода.
*/
func NewSubscription(serviceName string, price int, userID uuid.UUID, startDate time.Time) *Subscription {
	now := time.Now()
	return &Subscription{
		id:          uuid.New(),
		serviceName: serviceName,
		price:       price,
		userID:      userID,
		startDate:   startDate,
		createdAt:   now,
		updatedAt:   now,
	}
}

/** Геттер и сеттер для ID. Сеттер есть, чтобы можно было восстановить данные из БД. */
func (s *Subscription) ID() uuid.UUID {
	return s.id
}

func (s *Subscription) SetID(id uuid.UUID) {
	s.id = id
}

/** Аналогично — геттер/сеттер для имени сервиса. Сеттер обновляет updatedAt. */
func (s *Subscription) ServiceName() string {
	return s.serviceName
}

func (s *Subscription) SetServiceName(serviceName string) {
	s.serviceName = serviceName
	s.updatedAt = time.Now()
}

/** Управление ценой подписки. */
func (s *Subscription) Price() int {
	return s.price
}

func (s *Subscription) SetPrice(price int) {
	s.price = price
	s.updatedAt = time.Now()
}

/** Привязка к конкретному пользователю. */
func (s *Subscription) UserID() uuid.UUID {
	return s.userID
}

func (s *Subscription) SetUserID(userID uuid.UUID) {
	s.userID = userID
	s.updatedAt = time.Now()
}

/** Даты начала и конца подписки. */
func (s *Subscription) StartDate() time.Time {
	return s.startDate
}

func (s *Subscription) SetStartDate(startDate time.Time) {
	s.startDate = startDate
	s.updatedAt = time.Now()
}

func (s *Subscription) EndDate() *time.Time {
	return s.endDate
}

func (s *Subscription) SetEndDate(endDate *time.Time) {
	s.endDate = endDate
	s.updatedAt = time.Now()
}

/** Метаданные о создании и обновлении. */
func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Subscription) SetCreatedAt(createdAt time.Time) {
	s.createdAt = createdAt
}

func (s *Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Subscription) SetUpdatedAt(updatedAt time.Time) {
	s.updatedAt = updatedAt
}

/** Проверяет, активна ли подписка на конкретную дату. */
func (s *Subscription) IsActive(date time.Time) bool {
	if date.Before(s.startDate) {
		return false
	}
	if s.endDate != nil && date.After(*s.endDate) {
		return false
	}
	return true
}

/** Проверяет, истекла ли подписка на указанную дату. */
func (s *Subscription) IsExpired(date time.Time) bool {
	if s.endDate == nil {
		return false
	}
	return date.After(*s.endDate)
}

/*
*
CalculateCostForPeriod считает стоимость подписки за определённый диапазон дат.
Рассчёт идёт по количеству месяцев, начиная от startDate и до endDate (если есть).
*/
func (s *Subscription) CalculateCostForPeriod(from, to time.Time) int {
	if !s.IsActive(from) && !s.IsActive(to) {
		return 0
	}

	start := s.startDate
	if from.After(start) {
		start = from
	}

	end := to
	if s.endDate != nil && s.endDate.Before(end) {
		end = *s.endDate
	}

	if start.After(end) {
		return 0
	}

	startMonth := start.Year()*12 + int(start.Month()) - 1
	endMonth := end.Year()*12 + int(end.Month()) - 1

	months := endMonth - startMonth + 1
	if months <= 0 {
		return 0
	}

	return s.price * months
}

/*
*
Validate проверяет, что обязательные поля заполнены корректно:
- название сервиса не пустое
- цена > 0
- userID задан
- дата окончания не раньше даты начала
*/
func (s *Subscription) Validate() error {
	if s.serviceName == "" {
		return errors.New("service name cannot be empty")
	}
	if s.price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if s.userID == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}
	if s.endDate != nil && s.endDate.Before(s.startDate) {
		return errors.New("end date cannot be before start date")
	}
	return nil
}
