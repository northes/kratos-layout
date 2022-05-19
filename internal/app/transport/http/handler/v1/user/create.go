package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/northes/kratos-layout/internal/app/service/user"
)

func (h *Handler) Create(c *gin.Context) {
	req := new(user.CreateRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error(err)
		c.String(http.StatusOK, err.Error())
		return
	}

	_, err := h.service.Create(c.Request.Context(), *req)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
	return
}
