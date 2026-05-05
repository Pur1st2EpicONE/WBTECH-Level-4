// Package handler provides HTTP layer initialization, including routing,
// middleware, and static content handling.
//
// It wires service layer handlers into Gin, applies cross-cutting concerns
// (logging, recovery), and exposes both API and UI endpoints.
package handler

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"L4.3/internal/errs"
	v1 "L4.3/internal/handler/v1"
	"L4.3/internal/service"
	"L4.3/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

// NewHandler constructs and configures the HTTP handler.
func NewHandler(service service.Service, logger logger.Logger) http.Handler {

	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(middleware(logger))
	handler.Static("/static", "./web/static")

	apiV1 := handler.Group("/api/v1")
	handlerV1 := v1.NewHandler(service, logger)

	apiV1.POST("/create_event", handlerV1.CreateEvent)
	apiV1.POST("/update_event", handlerV1.UpdateEvent)
	apiV1.POST("/delete_event", handlerV1.DeleteEvent)

	apiV1.GET("/events_for_day", handlerV1.GetEventsDay)
	apiV1.GET("/events_for_week", handlerV1.GetEventsWeek)
	apiV1.GET("/events_for_month", handlerV1.GetEventsMonth)

	handler.GET("/", renderPage(template.Must(template.ParseFiles("web/templates/index.html"))))

	return handler

}

// middleware is a request logging middleware.
//
// It enriches each request with a unique request_id and logs structured data
// after the request is processed.

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

// renderPage returns a handler that renders an HTML template.
//
// It sets the appropriate content type and executes the provided template.
// In case of rendering error, it responds with a generic internal error.
//
// Intended for serving the root UI page.
func renderPage(tmpl *template.Template) gin.HandlerFunc {
	return func(c *ginext.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			c.String(http.StatusInternalServerError, errs.ErrInternal.Error())
		}
	}
}
