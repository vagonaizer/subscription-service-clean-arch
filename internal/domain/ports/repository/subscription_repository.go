package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error)
	GetAll(ctx context.Context, filter *models.SubscriptionFilter, limit, offset int) ([]*models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCostForPeriod(ctx context.Context, filter *models.SubscriptionFilter, period *models.DatePeriod) (int, error)
	Count(ctx context.Context, filter *models.SubscriptionFilter) (int, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
