package response

import "time"

type SubscriptionResponse struct {
	ID          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" example:"07-2025"`
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-15T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-01-15T10:30:00Z"`
}

type SubscriptionsListResponse struct {
	Data       []SubscriptionResponse `json:"data"`
	Pagination PaginationResponse     `json:"pagination"`
}

type CostSummaryResponse struct {
	TotalCost int            `json:"total_cost" example:"2400"`
	Period    PeriodResponse `json:"period"`
	Currency  string         `json:"currency" example:"RUB"`
}

type PeriodResponse struct {
	StartDate string `json:"start_date" example:"01-2025"`
	EndDate   string `json:"end_date" example:"06-2025"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

type StatsResponse struct {
	TotalSubscriptions int `json:"total_subscriptions"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
