package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/domain/ports/service"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/dto/request"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/dto/response"
	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/mappers"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/utils"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *logger.Logger
}

func NewSubscriptionHandler(service service.SubscriptionService, logger *logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger.Named("subscription-handler"),
	}
}

func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.CreateSubscription)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
		subscriptions.GET("/", h.GetSubscriptions)
	}

	users := router.Group("/users")
	{
		users.GET("/:user_id/subscriptions", h.GetUserSubscriptions)
		users.GET("/:user_id/subscriptions/stats", h.GetUserStats)
	}

	costs := router.Group("/costs")
	{
		costs.GET("/calculate", h.CalculateTotalCost)
	}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body request.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} response.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ValidationErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req request.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.Error(apperror.InvalidInput("request_body", err.Error()))
		return
	}

	userID, err := req.GetUserID()
	if err != nil {
		c.Error(apperror.InvalidUserID(req.UserID))
		return
	}

	subscription, err := h.service.CreateSubscription(
		c.Request.Context(),
		req.ServiceName,
		req.Price,
		userID,
		req.StartDate,
		utils.StringPtr(req.EndDate),
	)
	if err != nil {
		c.Error(err)
		return
	}

	resp := mappers.SubscriptionToResponse(subscription)
	h.logger.Info("subscription created successfully",
		zap.String("subscription_id", resp.ID),
		zap.String("service_name", resp.ServiceName))

	c.JSON(http.StatusCreated, resp)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get a single subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID" format(uuid)
// @Success 200 {object} response.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	req := request.GetSubscriptionRequest{
		ID: c.Param("id"),
	}

	id, err := req.GetID()
	if err != nil {
		c.Error(apperror.InvalidInput("id", err.Error()))
		return
	}

	subscription, err := h.service.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	resp := mappers.SubscriptionToResponse(subscription)
	c.JSON(http.StatusOK, resp)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" format(uuid)
// @Param subscription body request.UpdateSubscriptionRequest true "Updated subscription data"
// @Success 200 {object} response.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ValidationErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := utils.ValidateUUID(id, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req request.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.Error(apperror.InvalidInput("request_body", err.Error()))
		return
	}

	subscription, err := h.service.UpdateSubscription(
		c.Request.Context(),
		parsedID,
		req.ServiceName,
		req.Price,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.Error(err)
		return
	}

	resp := mappers.SubscriptionToResponse(subscription)
	h.logger.Info("subscription updated successfully",
		zap.String("subscription_id", resp.ID))

	c.JSON(http.StatusOK, resp)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID" format(uuid)
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	req := request.DeleteSubscriptionRequest{
		ID: c.Param("id"),
	}

	id, err := req.GetID()
	if err != nil {
		c.Error(apperror.InvalidInput("id", err.Error()))
		return
	}

	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	h.logger.Info("subscription deleted successfully",
		zap.String("subscription_id", id.String()))

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "Subscription deleted successfully",
	})
}

// GetSubscriptions godoc
// @Summary List subscriptions
// @Description Get list of subscriptions with optional filtering
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter" format(uuid)
// @Param service_name query string false "Service name filter"
// @Param start_date query string false "Start date filter (MM-YYYY format)"
// @Param end_date query string false "End date filter (MM-YYYY format)"
// @Param limit query int false "Limit number of results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} response.SubscriptionsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	req := h.parseGetSubscriptionsRequest(c)

	filter, err := mappers.SubscriptionFilterFromRequest(
		req.UserID,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.Error(err)
		return
	}

	subscriptions, err := h.service.GetAllSubscriptions(
		c.Request.Context(),
		filter,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		c.Error(err)
		return
	}

	pagination := response.NewPaginationResponse(req.Limit, req.Offset, nil)
	resp := mappers.SubscriptionsToListResponse(subscriptions, pagination)

	h.logger.Debug("subscriptions retrieved",
		zap.Int("count", len(subscriptions)),
		zap.Int("limit", req.Limit),
		zap.Int("offset", req.Offset))

	c.JSON(http.StatusOK, resp)
}

// GetUserSubscriptions godoc
// @Summary Get user subscriptions
// @Description Get all subscriptions for a specific user
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param limit query int false "Limit number of results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} response.SubscriptionsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{user_id}/subscriptions [get]
func (h *SubscriptionHandler) GetUserSubscriptions(c *gin.Context) {
	req := request.GetUserSubscriptionsRequest{
		UserID: c.Param("user_id"),
		Limit:  h.parseIntQuery(c, "limit", 20),
		Offset: h.parseIntQuery(c, "offset", 0),
	}

	userID, err := req.GetUserID()
	if err != nil {
		c.Error(apperror.InvalidUserID(req.UserID))
		return
	}

	subscriptions, err := h.service.GetSubscriptionsByUser(
		c.Request.Context(),
		userID,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		c.Error(err)
		return
	}

	pagination := response.NewPaginationResponse(req.Limit, req.Offset, nil)
	resp := mappers.SubscriptionsToListResponse(subscriptions, pagination)

	h.logger.Debug("user subscriptions retrieved",
		zap.String("user_id", userID.String()),
		zap.Int("count", len(subscriptions)))

	c.JSON(http.StatusOK, resp)
}

// GetUserStats godoc
// @Summary Get user subscription statistics
// @Description Get total number of subscriptions for a user
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Success 200 {object} response.StatsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{user_id}/subscriptions/stats [get]
func (h *SubscriptionHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("user_id")
	parsedUserID, err := utils.ValidateUUID(userID, "user_id")
	if err != nil {
		c.Error(err)
		return
	}

	count, err := h.service.GetSubscriptionStats(c.Request.Context(), &parsedUserID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := response.StatsResponse{
		TotalSubscriptions: count,
	}

	c.JSON(http.StatusOK, resp)
}

// CalculateTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculate total cost of subscriptions for a given period with optional filtering
// @Tags costs
// @Produce json
// @Param user_id query string false "User ID filter" format(uuid)
// @Param service_name query string false "Service name filter"
// @Param start_date query string true "Start date (MM-YYYY format)"
// @Param end_date query string true "End date (MM-YYYY format)"
// @Success 200 {object} response.CostSummaryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ValidationErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /costs/calculate [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	req := h.parseCalculateCostRequest(c)

	var userID *uuid.UUID
	if req.UserID != nil && *req.UserID != "" {
		parsedUserID, err := utils.ValidateUUID(*req.UserID, "user_id")
		if err != nil {
			c.Error(err)
			return
		}
		userID = &parsedUserID
	}

	summary, err := h.service.CalculateTotalCost(
		c.Request.Context(),
		userID,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.Error(err)
		return
	}

	resp := mappers.CostSummaryToResponse(summary)

	h.logger.Info("cost calculated successfully",
		zap.Int("total_cost", resp.TotalCost),
		zap.String("period", req.StartDate+" to "+req.EndDate))

	c.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) parseGetSubscriptionsRequest(c *gin.Context) request.GetSubscriptionsRequest {
	return request.GetSubscriptionsRequest{
		UserID:      h.parseStringQuery(c, "user_id"),
		ServiceName: h.parseStringQuery(c, "service_name"),
		StartDate:   h.parseStringQuery(c, "start_date"),
		EndDate:     h.parseStringQuery(c, "end_date"),
		Limit:       h.parseIntQuery(c, "limit", 20),
		Offset:      h.parseIntQuery(c, "offset", 0),
	}
}

func (h *SubscriptionHandler) parseCalculateCostRequest(c *gin.Context) request.CalculateCostRequest {
	return request.CalculateCostRequest{
		UserID:      h.parseStringQuery(c, "user_id"),
		ServiceName: h.parseStringQuery(c, "service_name"),
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
	}
}

func (h *SubscriptionHandler) parseStringQuery(c *gin.Context, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}

func (h *SubscriptionHandler) parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
