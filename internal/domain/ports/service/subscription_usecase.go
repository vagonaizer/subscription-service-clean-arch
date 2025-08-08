package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/models"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, serviceName string, price int, userID uuid.UUID, startDate string, endDate *string) (*models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetSubscriptionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Subscription, error)
	GetAllSubscriptions(ctx context.Context, filter *models.SubscriptionFilter, limit, offset int) ([]*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, serviceName *string, price *int, startDate *string, endDate *string) (*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate string) (*models.CostSummary, error)
	GetSubscriptionStats(ctx context.Context, userID *uuid.UUID) (int, error)
}
