package v1

import (
	"L4.3/internal/errs"
	"L4.3/internal/models"
	"github.com/gin-gonic/gin"
)

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
