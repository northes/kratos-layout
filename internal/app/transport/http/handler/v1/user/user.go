package user

import (
	"github.com/gin-gonic/gin"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/northes/kratos-layout/internal/app/service/user"
)

type HandlerInterface interface {
	Create(c *gin.Context)
}

var _ HandlerInterface = (*Handler)(nil)

type Handler struct {
	logger  *log.Helper
	service user.IService
}

func NewHandler(logger log.Logger, service user.IService) *Handler {
	return &Handler{
		logger:  log.NewHelper(logger),
		service: service,
	}
}
