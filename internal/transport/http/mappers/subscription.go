package mappers

import (
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/models"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/dto/response"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/utils"
)

func SubscriptionToResponse(subscription *models.Subscription) response.SubscriptionResponse {
	resp := response.SubscriptionResponse{
		ID:          subscription.ID().String(),
		ServiceName: subscription.ServiceName(),
		Price:       subscription.Price(),
		UserID:      subscription.UserID().String(),
		StartDate:   utils.FormatMonthYear(subscription.StartDate()),
		CreatedAt:   subscription.CreatedAt(),
		UpdatedAt:   subscription.UpdatedAt(),
	}

	if subscription.EndDate() != nil {
		endDate := utils.FormatMonthYear(*subscription.EndDate())
		resp.EndDate = &endDate
	}

	return resp
}

func SubscriptionsToListResponse(subscriptions []*models.Subscription, pagination response.PaginationResponse) response.SubscriptionsListResponse {
	data := make([]response.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		data[i] = SubscriptionToResponse(subscription)
	}

	return response.SubscriptionsListResponse{
		Data:       data,
		Pagination: pagination,
	}
}

func CostSummaryToResponse(summary *models.CostSummary) response.CostSummaryResponse {
	period := summary.Period()
	return response.CostSummaryResponse{
		TotalCost: summary.TotalCost(),
		Period: response.PeriodResponse{
			StartDate: utils.FormatMonthYear(period.From()),
			EndDate:   utils.FormatMonthYear(period.To()),
		},
		Currency: "RUB",
	}
}

func SubscriptionFilterFromRequest(userID *string, serviceName *string, startDate *string, endDate *string) (*models.SubscriptionFilter, error) {
	filter := models.NewSubscriptionFilter()

	if userID != nil && *userID != "" {
		parsedUserID, err := utils.ValidateUUID(*userID, "user_id")
		if err != nil {
			return nil, err
		}
		filter.SetUserID(&parsedUserID)
	}

	if serviceName != nil && *serviceName != "" {
		normalized := utils.NormalizeString(*serviceName)
		filter.SetServiceName(&normalized)
	}

	if startDate != nil && *startDate != "" {
		start, err := utils.ParseMonthYear(*startDate)
		if err != nil {
			return nil, err
		}
		start = utils.StartOfMonth(start)
		filter.SetStartDate(&start)
	}

	if endDate != nil && *endDate != "" {
		end, err := utils.ParseMonthYear(*endDate)
		if err != nil {
			return nil, err
		}
		end = utils.EndOfMonth(end)
		filter.SetEndDate(&end)
	}

	return filter, nil
}
