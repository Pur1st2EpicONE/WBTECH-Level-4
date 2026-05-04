package v1

import (
	"L4.3/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetEventsDay(c *gin.Context) {
	h.getEvents(c, models.Day)
}

func (h *Handler) GetEventsWeek(c *gin.Context) {
	h.getEvents(c, models.Week)
}

func (h *Handler) GetEventsMonth(c *gin.Context) {
	h.getEvents(c, models.Month)
}

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
