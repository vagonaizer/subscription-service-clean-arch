package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vagonaizer/effective-mobile/subscription-service/pkg/logger"
)

type Router struct {
	engine *gin.Engine
	logger *logger.Logger
}

type RouterConfig struct {
	Debug  bool
	Logger *logger.Logger
}

func New(config RouterConfig) *Router {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	return &Router{
		engine: engine,
		logger: config.Logger,
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) SetupMiddleware(middlewares ...gin.HandlerFunc) {
	r.engine.Use(middlewares...)
}

func (r *Router) RegisterHealthRoutes() {
	health := r.engine.Group("/health")
	{
		health.GET("/", r.handleHealthCheck)
		health.GET("/ready", r.handleReadiness)
		health.GET("/live", r.handleLiveness)
	}
}

func (r *Router) RegisterAPIRoutes(handlers ...RouteHandler) {
	api := r.engine.Group("/api")
	v1 := api.Group("/v1")

	for _, handler := range handlers {
		handler.RegisterRoutes(v1)
	}
}

func (r *Router) RegisterSwaggerRoutes() {
	r.logger.Info("registering swagger routes")

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	r.engine.Static("/api-docs", "./api/swagger")

	r.logger.Info("swagger documentation available at /swagger/index.html")
}

func (r *Router) handleHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"service":   "subscription-service",
		"timestamp": gin.H{},
	})
}

func (r *Router) handleReadiness(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ready",
	})
}

func (r *Router) handleLiveness(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

type RouteHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}
