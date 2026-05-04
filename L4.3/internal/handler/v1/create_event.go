package v1

import (
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	"github.com/gin-gonic/gin"
)

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
		Data: models.Data{Text: request.Text, Reminder: time.Duration(request.Reminder) * time.Minute}}

	eventID, err := h.service.CreateEvent(&event)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, CreateResponseV1{EventID: eventID})

}
