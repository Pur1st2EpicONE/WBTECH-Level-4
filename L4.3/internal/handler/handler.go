// Package handler provides HTTP handler initialization for the application,
// including API routing, middleware, and Swagger documentation endpoint.
package handler

import (
	"fmt"
	"net/http"
	"time"

	v1 "L4.3/internal/handler/v1"
	"L4.3/internal/service"
	"L4.3/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHandler creates and configures the HTTP handler for the application.
//
// It sets up the Gin engine, registers middleware, API v1 routes, and the
// Swagger documentation endpoint.
//
// Parameters:
// - service: the service layer instance that provides business logic
// - logger: logger instance to log requests and errors
//
// Returns:
// - http.Handler instance ready to be served by a HTTP server
func NewHandler(service service.Service, logger logger.Logger) http.Handler {

	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(middleware(logger))

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service, logger)

	apiV1.POST("/create_event", handlerV1.CreateEvent)
	apiV1.POST("/update_event", handlerV1.UpdateEvent)
	apiV1.POST("/delete_event", handlerV1.DeleteEvent)

	apiV1.GET("/events_for_day", handlerV1.GetEventsDay)
	apiV1.GET("/events_for_week", handlerV1.GetEventsWeek)
	apiV1.GET("/events_for_month", handlerV1.GetEventsMonth)

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return handler

}

// middleware creates a Gin middleware that logs incoming HTTP requests and their outcomes.
//
// It generates a request ID, measures request latency, and logs request details
// including method, path, query string, client IP, HTTP status, user agent, and Gin errors.
//
// Logging behavior based on HTTP status:
// - 500: LogError
// - 400, 503: LogWarn
// - others: LogInfo
//
// Parameters:
// - logger: logger instance to log request details
//
// Returns:
// - gin.HandlerFunc that can be used as middleware
func middleware(logger logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []any{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"status", status,
			"query", query,
			"proto", c.Request.Proto,
			"user_agent", c.Request.UserAgent(),
			"gin_errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"layer", "handler",
		}

		msg := fmt.Sprintf("handler — received %s request to %s", c.Request.Method, path)

		switch status {
		case 500:
			logger.LogError(msg, nil, fields...)
		case 400, 503:
			logger.LogWarn(msg, fields...)
		default:
			logger.LogInfo(msg, fields...)
		}

	}

}
