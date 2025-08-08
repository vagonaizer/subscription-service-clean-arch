package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/models"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/ports/repository"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/utils"
)

/*
subscriptionService — слой бизнес-логики для работы с подписками.
Отвечает за валидацию входных данных, вызов методов репозитория
и запись логов.
*/
type subscriptionService struct {
	repo repository.SubscriptionRepository
	log  *logger.Logger
}

/** Конструктор сервиса, принимает репозиторий и логгер. */
func NewSubscriptionService(repo repository.SubscriptionRepository, log *logger.Logger) *subscriptionService {
	return &subscriptionService{
		repo: repo,
		log:  log.Named("subscription-service"),
	}
}

/*
CreateSubscription — создаёт новую подписку.
- Валидирует входные данные.
- Парсит даты начала/окончания.
- Проверяет корректность диапазона.
- Сохраняет подписку через репозиторий.
*/
func (s *subscriptionService) CreateSubscription(ctx context.Context, serviceName string, price int, userID uuid.UUID, startDate string, endDate *string) (*models.Subscription, error) {
	s.log.Debug("creating subscription",
		zap.String("service_name", serviceName),
		zap.Int("price", price),
		zap.String("user_id", userID.String()))

	if err := s.validateCreateInput(serviceName, price, userID); err != nil {
		return nil, err
	}

	startTime, err := utils.ParseMonthYear(startDate)
	if err != nil {
		return nil, err
	}
	startTime = utils.StartOfMonth(startTime)

	subscription := models.NewSubscription(
		utils.NormalizeString(serviceName),
		price,
		userID,
		startTime,
	)

	if endDate != nil && *endDate != "" {
		endTime, err := utils.ParseMonthYear(*endDate)
		if err != nil {
			return nil, err
		}
		endTime = utils.EndOfMonth(endTime)

		if err := utils.ValidateDateRange(&startTime, &endTime); err != nil {
			return nil, err
		}

		subscription.SetEndDate(&endTime)
	}

	if err := subscription.Validate(); err != nil {
		return nil, apperror.InvalidSubscriptionData("subscription", err.Error())
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		s.log.Error("failed to create subscription", zap.Error(err))
		return nil, err
	}

	s.log.Info("subscription created successfully",
		zap.String("subscription_id", subscription.ID().String()),
		zap.String("service_name", serviceName))

	return subscription, nil
}

/** Получает подписку по ID, возвращает ошибку если не найдена. */
func (s *subscriptionService) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	s.log.Debug("getting subscription by id", zap.String("subscription_id", id.String()))

	if id == uuid.Nil {
		return nil, apperror.InvalidInput("id", "cannot be empty")
	}

	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return nil, apperror.SubscriptionNotFound(id.String())
	}

	return subscription, nil
}

/** Получает подписки по ID пользователя с пагинацией. */
func (s *subscriptionService) GetSubscriptionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error) {
	s.log.Debug("getting subscriptions by user",
		zap.String("user_id", userID.String()),
		zap.Int("limit", limit),
		zap.Int("offset", offset))

	if userID == uuid.Nil {
		return nil, apperror.InvalidUserID(userID.String())
	}

	limit, offset, err := utils.ValidatePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	subscriptions, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	s.log.Debug("retrieved subscriptions by user",
		zap.String("user_id", userID.String()),
		zap.Int("count", len(subscriptions)))

	return subscriptions, nil
}

/** Получает все подписки с фильтром и пагинацией. */
func (s *subscriptionService) GetAllSubscriptions(ctx context.Context, filter *models.SubscriptionFilter, limit, offset int) ([]*models.Subscription, error) {
	s.log.Debug("getting filtered subscriptions",
		zap.Int("limit", limit),
		zap.Int("offset", offset))

	if filter == nil {
		filter = models.NewSubscriptionFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, apperror.InvalidFilterParams("filter", err.Error())
	}

	limit, offset, err := utils.ValidatePagination(limit, offset)
	if err != nil {
		return nil, err
	}

	subscriptions, err := s.repo.GetAll(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	s.log.Debug("retrieved filtered subscriptions",
		zap.Int("count", len(subscriptions)))

	return subscriptions, nil
}

/*
UpdateSubscription — обновляет существующую подписку.
Обновляет только те поля, которые переданы и изменились.
*/
func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, serviceName *string, price *int, startDate *string, endDate *string) (*models.Subscription, error) {
	s.log.Debug("updating subscription", zap.String("subscription_id", id.String()))

	subscription, err := s.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.validateUpdateInput(serviceName, price); err != nil {
		return nil, err
	}

	hasChanges := false

	if serviceName != nil && *serviceName != "" {
		normalized := utils.NormalizeString(*serviceName)
		if normalized != subscription.ServiceName() {
			subscription.SetServiceName(normalized)
			hasChanges = true
		}
	}

	if price != nil && *price != subscription.Price() {
		subscription.SetPrice(*price)
		hasChanges = true
	}

	if startDate != nil && *startDate != "" {
		newStartDate, err := utils.ParseMonthYear(*startDate)
		if err != nil {
			return nil, err
		}
		newStartDate = utils.StartOfMonth(newStartDate)
		subscription.SetStartDate(newStartDate)
		hasChanges = true
	}

	if endDate != nil {
		if *endDate == "" {
			subscription.SetEndDate(nil)
			hasChanges = true
		} else {
			newEndDate, err := utils.ParseMonthYear(*endDate)
			if err != nil {
				return nil, err
			}
			newEndDate = utils.EndOfMonth(newEndDate)
			subscription.SetEndDate(&newEndDate)
			hasChanges = true
		}
	}

	if !hasChanges {
		return subscription, nil
	}

	if err := subscription.Validate(); err != nil {
		return nil, apperror.InvalidSubscriptionData("subscription", err.Error())
	}

	if err := s.repo.Update(ctx, subscription); err != nil {
		s.log.Error("failed to update subscription", zap.Error(err))
		return nil, err
	}

	s.log.Info("subscription updated successfully",
		zap.String("subscription_id", id.String()))

	return subscription, nil
}

/** Удаляет подписку по ID, проверяя её существование. */
func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	s.log.Debug("deleting subscription", zap.String("subscription_id", id.String()))

	if id == uuid.Nil {
		return apperror.InvalidInput("id", "cannot be empty")
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.SubscriptionNotFound(id.String())
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete subscription", zap.Error(err))
		return err
	}

	s.log.Info("subscription deleted successfully",
		zap.String("subscription_id", id.String()))

	return nil
}

/*
CalculateTotalCost — считает общую стоимость подписок за период.
Можно фильтровать по userID и имени сервиса.
*/
func (s *subscriptionService) CalculateTotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate string) (*models.CostSummary, error) {
	s.log.Debug("calculating total cost",
		zap.String("start_date", startDate),
		zap.String("end_date", endDate))

	startTime, endTime, err := utils.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if startTime == nil || endTime == nil {
		return nil, apperror.InvalidInput("date_range", "both start_date and end_date are required")
	}

	period := models.NewDatePeriod(*startTime, *endTime)
	if err := period.Validate(); err != nil {
		return nil, apperror.InvalidDateRange(startDate, endDate)
	}

	filter := models.NewSubscriptionFilter()
	if userID != nil {
		filter.SetUserID(userID)
	}
	if serviceName != nil && *serviceName != "" {
		normalized := utils.NormalizeString(*serviceName)
		filter.SetServiceName(&normalized)
	}

	totalCost, err := s.repo.GetTotalCostForPeriod(ctx, filter, period)
	if err != nil {
		return nil, err
	}

	summary := models.NewCostSummary(*period)
	summary.SetTotalCost(totalCost)

	s.log.Info("calculated total cost",
		zap.Int("total_cost", totalCost),
		zap.String("period", startDate+" to "+endDate))

	return summary, nil
}

/** Возвращает количество подписок (с фильтром по userID, если задан). */
func (s *subscriptionService) GetSubscriptionStats(ctx context.Context, userID *uuid.UUID) (int, error) {
	s.log.Debug("getting subscription stats")

	filter := models.NewSubscriptionFilter()
	if userID != nil {
		filter.SetUserID(userID)
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

/** Валидация входных данных для создания подписки. */
func (s *subscriptionService) validateCreateInput(serviceName string, price int, userID uuid.UUID) error {
	if err := utils.ValidateServiceName(serviceName); err != nil {
		return err
	}

	if err := utils.ValidatePrice(price); err != nil {
		return err
	}

	if userID == uuid.Nil {
		return apperror.InvalidUserID(userID.String())
	}

	return nil
}

/** Валидация входных данных для обновления подписки. */
func (s *subscriptionService) validateUpdateInput(serviceName *string, price *int) error {
	if serviceName != nil && *serviceName != "" {
		if err := utils.ValidateServiceName(*serviceName); err != nil {
			return err
		}
	}

	if price != nil {
		if err := utils.ValidatePrice(*price); err != nil {
			return err
		}
	}

	return nil
}
