package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/transport/http/dto/response"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/apperror"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		stack := string(debug.Stack())

		log.Error("panic recovered",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Any("error", recovered),
			zap.String("stack", stack))

		errorResp := response.NewErrorResponse(
			apperror.CodeInternalError,
			"Internal server error occurred",
			map[string]string{
				"panic": fmt.Sprintf("%v", recovered),
			},
			requestID,
		)

		c.Header("Content-Type", "application/json")
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResp)
	})
}

func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		err := c.Errors.Last().Err

		if appErr, ok := apperror.IsAppError(err); ok {
			log.Warn("application error occurred",
				zap.String("request_id", requestID),
				zap.String("error_code", appErr.Code()),
				zap.String("error_message", appErr.Message()),
				zap.Error(appErr.Cause()))

			errorResp := response.NewErrorResponse(
				appErr.Code(),
				appErr.Message(),
				appErr.Details(),
				requestID,
			)

			c.Header("Content-Type", "application/json")
			c.AbortWithStatusJSON(appErr.HTTPStatus(), errorResp)
			return
		}

		log.Error("unexpected error occurred",
			zap.String("request_id", requestID),
			zap.Error(err))

		errorResp := response.NewErrorResponse(
			apperror.CodeInternalError,
			"An unexpected error occurred",
			nil,
			requestID,
		)

		c.Header("Content-Type", "application/json")
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResp)
	}
}
