package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/models"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/infrastructure/database/postgres"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type subscriptionRepository struct {
	db  *postgres.DB
	log *logger.Logger
}

func NewSubscriptionRepository(db *postgres.DB, log *logger.Logger) *subscriptionRepository {
	return &subscriptionRepository{
		db:  db,
		log: log.Named("subscription-repository"),
	}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Pool().Exec(ctx, query,
		subscription.ID(),
		subscription.ServiceName(),
		subscription.Price(),
		subscription.UserID(),
		subscription.StartDate(),
		subscription.EndDate(),
		subscription.CreatedAt(),
		subscription.UpdatedAt(),
	)

	if err != nil {
		r.log.Error("failed to create subscription",
			zap.String("subscription_id", subscription.ID().String()),
			zap.Error(err))
		return apperror.DatabaseError("create subscription", err)
	}

	r.log.Debug("subscription created",
		zap.String("subscription_id", subscription.ID().String()),
		zap.String("service_name", subscription.ServiceName()))

	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions 
		WHERE id = $1`

	row := r.db.Pool().QueryRow(ctx, query, id)

	subscription, err := r.scanSubscription(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		r.log.Error("failed to get subscription by id",
			zap.String("subscription_id", id.String()),
			zap.Error(err))
		return nil, apperror.DatabaseError("get subscription by id", err)
	}

	return subscription, nil
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Pool().Query(ctx, query, userID, limit, offset)
	if err != nil {
		r.log.Error("failed to get subscriptions by user id",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, fmt.Errorf("get subscriptions by user id: %w", err)
	}
	defer rows.Close()

	return r.scanSubscriptions(rows)
}

func (r *subscriptionRepository) GetAll(ctx context.Context, filter *models.SubscriptionFilter, limit, offset int) ([]*models.Subscription, error) {
	query, args := r.buildFilterQuery(filter, limit, offset)

	rows, err := r.db.Pool().Query(ctx, query, args...)
	if err != nil {
		r.log.Error("failed to get filtered subscriptions", zap.Error(err))
		return nil, fmt.Errorf("get filtered subscriptions: %w", err)
	}
	defer rows.Close()

	return r.scanSubscriptions(rows)
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.db.Pool().Exec(ctx, query,
		subscription.ID(),
		subscription.ServiceName(),
		subscription.Price(),
		subscription.UserID(),
		subscription.StartDate(),
		subscription.EndDate(),
		subscription.UpdatedAt(),
	)

	if err != nil {
		r.log.Error("failed to update subscription",
			zap.String("subscription_id", subscription.ID().String()),
			zap.Error(err))
		return fmt.Errorf("update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	r.log.Debug("subscription updated",
		zap.String("subscription_id", subscription.ID().String()))

	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Pool().Exec(ctx, query, id)
	if err != nil {
		r.log.Error("failed to delete subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err))
		return fmt.Errorf("delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	r.log.Debug("subscription deleted",
		zap.String("subscription_id", id.String()))

	return nil
}

func (r *subscriptionRepository) GetTotalCostForPeriod(ctx context.Context, filter *models.SubscriptionFilter, period *models.DatePeriod) (int, error) {
	baseQuery := `
		SELECT COALESCE(SUM(price), 0) as total_cost
		FROM subscriptions
		WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`

	args := []interface{}{period.To(), period.From()}
	conditions := []string{}
	argIndex := 3

	if filter.HasUserID() {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID())
		argIndex++
	}

	if filter.HasServiceName() {
		conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argIndex))
		args = append(args, "%"+*filter.ServiceName()+"%")
		argIndex++
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var totalCost int
	err := r.db.Pool().QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		r.log.Error("failed to get total cost for period", zap.Error(err))
		return 0, fmt.Errorf("get total cost for period: %w", err)
	}

	return totalCost, nil
}

func (r *subscriptionRepository) Count(ctx context.Context, filter *models.SubscriptionFilter) (int, error) {
	query, args := r.buildCountQuery(filter)

	var count int
	err := r.db.Pool().QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		r.log.Error("failed to count subscriptions", zap.Error(err))
		return 0, fmt.Errorf("count subscriptions: %w", err)
	}

	return count, nil
}

func (r *subscriptionRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM subscriptions WHERE id = $1)`

	var exists bool
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		r.log.Error("failed to check subscription existence",
			zap.String("subscription_id", id.String()),
			zap.Error(err))
		return false, fmt.Errorf("check subscription existence: %w", err)
	}

	return exists, nil
}

func (r *subscriptionRepository) scanSubscription(row pgx.Row) (*models.Subscription, error) {
	var (
		id          uuid.UUID
		serviceName string
		price       int
		userID      uuid.UUID
		startDate   time.Time
		endDate     *time.Time
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := row.Scan(&id, &serviceName, &price, &userID, &startDate, &endDate, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	subscription := models.NewSubscription(serviceName, price, userID, startDate)
	subscription.SetID(id)
	subscription.SetEndDate(endDate)
	subscription.SetCreatedAt(createdAt)
	subscription.SetUpdatedAt(updatedAt)

	return subscription, nil
}

func (r *subscriptionRepository) scanSubscriptions(rows pgx.Rows) ([]*models.Subscription, error) {
	subscriptions := make([]*models.Subscription, 0)

	for rows.Next() {
		subscription, err := r.scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) buildFilterQuery(filter *models.SubscriptionFilter, limit, offset int) (string, []interface{}) {
	baseQuery := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.HasUserID() {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID())
		argIndex++
	}

	if filter.HasServiceName() {
		conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argIndex))
		args = append(args, "%"+*filter.ServiceName()+"%")
		argIndex++
	}

	if filter.HasDateRange() {
		if filter.StartDate() != nil {
			conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
			args = append(args, *filter.StartDate())
			argIndex++
		}
		if filter.EndDate() != nil {
			conditions = append(conditions, fmt.Sprintf("(end_date IS NULL OR end_date <= $%d)", argIndex))
			args = append(args, *filter.EndDate())
			argIndex++
		}
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	return query, args
}

func (r *subscriptionRepository) buildCountQuery(filter *models.SubscriptionFilter) (string, []interface{}) {
	baseQuery := `SELECT COUNT(*) FROM subscriptions`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.HasUserID() {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID())
		argIndex++
	}

	if filter.HasServiceName() {
		conditions = append(conditions, fmt.Sprintf("service_name ILIKE $%d", argIndex))
		args = append(args, "%"+*filter.ServiceName()+"%")
		argIndex++
	}

	if filter.HasDateRange() {
		if filter.StartDate() != nil {
			conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
			args = append(args, *filter.StartDate())
			argIndex++
		}
		if filter.EndDate() != nil {
			conditions = append(conditions, fmt.Sprintf("(end_date IS NULL OR end_date <= $%d)", argIndex))
			args = append(args, *filter.EndDate())
			argIndex++
		}
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}
