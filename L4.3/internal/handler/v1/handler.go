// Package v1 provides version 1 of the API handlers for the event management system.
//
// It defines a Handler struct that wraps the service layer and logger, and exposes
// HTTP endpoints to create, update, delete, and retrieve user events. Each method
// is annotated for Swagger documentation generation and uses JSON for request
// and response payloads.
package v1

import (
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	"L4.3/internal/service"
	"L4.3/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Handler represents the API handler for version 1 of the event service.
//
// It holds references to the service layer (business logic) and logger.
// All methods on Handler are HTTP endpoints that operate on events.
type Handler struct {
	service service.Service // service handles the business logic for events
	logger  logger.Logger   // logger is used to log request processing and errors
}

// NewHandler creates a new Handler instance with the given service and logger.
//
// service: the business logic layer that the handler will call for event operations.
// logger: structured logger to log request and error information.
func NewHandler(service service.Service, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateEvent handles HTTP POST requests to create a new event.
// It validates the request body and calls the service layer to create the event.
//
// @Summary Create a new event
// @Description Creates an event for a user
// @Tags events
// @Accept json
// @Produce json
// @Param request body CreateRequestV1 true "Event data"
// @Success 200 {object} CreateResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/create_event [post]
func (h *Handler) CreateEvent(c *gin.Context) {

	var request CreateRequestV1

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	eventDate, err := parseDate(request.EventDate)
	if err != nil {
		respondError(c, err)
		return
	}

	event := models.Event{
		Meta: models.Meta{UserID: request.UserID, EventDate: eventDate},
		Data: models.Data{Text: request.Text},
	}

	eventID, err := h.service.CreateEvent(&event)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, CreateResponseV1{EventID: eventID})

}

// UpdateEvent handles HTTP POST requests to update an existing event.
// It validates the request body and updates the event's text and/or date.
//
// @Summary Update an existing event
// @Description Updates an event's text or date
// @Tags events
// @Accept json
// @Produce json
// @Param request body UpdateRequestV1 true "Event update data"
// @Success 200 {object} UpdateResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/update_event [post]
func (h *Handler) UpdateEvent(c *gin.Context) {

	var request UpdateRequestV1

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	var date time.Time
	var err error

	if request.NewDate != "" {
		date, err = parseDate(request.NewDate)
		if err != nil {
			respondError(c, errs.ErrInvalidDateFormat)
			return
		}

	}

	event := models.Event{
		Meta: models.Meta{UserID: request.UserID, EventID: request.EventID, NewDate: date},
		Data: models.Data{Text: request.Text}}

	if err := h.service.UpdateEvent(&event); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, UpdateResponseV1{Updated: true})

}

// DeleteEvent handles HTTP POST requests to delete an existing event.
//
// @Summary Delete an event
// @Description Deletes an event for a user by ID
// @Tags events
// @Accept json
// @Produce json
// @Param request body DeleteRequestV1 true "Event delete data"
// @Success 200 {object} DeleteResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/delete_event [post]
func (h *Handler) DeleteEvent(c *gin.Context) {

	var request DeleteRequestV1

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	meta := models.Meta{UserID: request.UserID, EventID: request.EventID}

	if err := h.service.DeleteEvent(&meta); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, DeleteResponseV1{Deleted: true})

}

// GetEventsDay handles HTTP GET requests to retrieve all events for a specific day.
//
// @Summary Get events for a day
// @Description Returns all events for a given day for a user
// @Tags events
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} ListOfEventsResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/events_for_day [get]
func (h *Handler) GetEventsDay(c *gin.Context) {
	h.getEvents(c, models.Day)
}

// GetEventsWeek handles HTTP GET requests to retrieve all events for a specific week.
//
// @Summary Get events for a week
// @Description Returns all events for a given week for a user
// @Tags events
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} ListOfEventsResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/events_for_week [get]
func (h *Handler) GetEventsWeek(c *gin.Context) {
	h.getEvents(c, models.Week)
}

// GetEventsMonth handles HTTP GET requests to retrieve all events for a specific month.
//
// @Summary Get events for a month
// @Description Returns all events for a given month for a user
// @Tags events
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} ListOfEventsResponseV1
// @Failure 400 {object} ErrorResponse400
// @Failure 500 {object} ErrorResponse500
// @Router /api/v1/events_for_month [get]
func (h *Handler) GetEventsMonth(c *gin.Context) {
	h.getEvents(c, models.Month)
}

// getEvents is a helper method to fetch events based on the given period type (day, week, month).
// It parses query parameters, calls the service layer, and returns the formatted response.
func (h *Handler) getEvents(c *gin.Context, period models.Period) {

	userId, eventDate, err := parseQuery(c.Query("user_id"), c.Query("date"))
	if err != nil {
		respondError(c, err)
		return
	}

	events, err := h.service.GetEvents(&models.Meta{UserID: userId, EventDate: eventDate}, period)
	if err != nil {
		respondError(c, err)
		return
	}

	respEvents := make([]EventDtoV1, len(events))

	for i, e := range events {
		respEvents[i] = EventDtoV1{Text: e.Data.Text, EventDate: e.Meta.EventDate.Format("2006-01-02"), EventID: e.Meta.EventID}
	}

	respondOK(c, ListOfEventsResponseV1{Events: respEvents})

}
