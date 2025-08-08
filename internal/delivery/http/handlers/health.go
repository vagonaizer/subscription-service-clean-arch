package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/dto/response"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type HealthHandler struct {
	logger      *logger.Logger
	healthCheck func(ctx context.Context) error
}

func NewHealthHandler(logger *logger.Logger, healthCheck func(ctx context.Context) error) *HealthHandler {
	return &HealthHandler{
		logger:      logger.Named("health-handler"),
		healthCheck: healthCheck,
	}
}

func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("/", h.Health)
		health.GET("/ready", h.Ready)
		health.GET("/live", h.Live)
	}
}

// Health godoc
// @Summary Health check
// @Description Get overall health status of the service and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} response.HealthResponse
// @Failure 503 {object} response.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)
	overallStatus := "healthy"

	if h.healthCheck != nil {
		if err := h.healthCheck(ctx); err != nil {
			h.logger.Error("health check failed", zap.Error(err))
			services["database"] = "unhealthy"
			overallStatus = "unhealthy"

			healthResp := response.HealthResponse{
				Status:    overallStatus,
				Timestamp: time.Now(),
				Services:  services,
			}

			c.JSON(http.StatusServiceUnavailable, healthResp)
			return
		}
		services["database"] = "healthy"
	}

	healthResp := response.HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
	}

	c.JSON(http.StatusOK, healthResp)
}

// Ready godoc
// @Summary Readiness check
// @Description Check if service is ready to accept traffic
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if h.healthCheck != nil {
		if err := h.healthCheck(ctx); err != nil {
			h.logger.Warn("readiness check failed", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "service dependencies unavailable",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Live godoc
// @Summary Liveness check
// @Description Check if service is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
