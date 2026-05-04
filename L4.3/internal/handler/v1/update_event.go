package v1

import (
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	"github.com/gin-gonic/gin"
)

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
