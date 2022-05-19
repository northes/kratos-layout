package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/northes/kratos-layout/internal/app/transport/http/handler/v1/user"
)

var ProviderSet = wire.NewSet(
	New,
)

func New(userHandler user.HandlerInterface) *gin.Engine {
	r := gin.New()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
		return
	})
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/user", userHandler.Create)
		}
	}
	return r
}
