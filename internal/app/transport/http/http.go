package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/northes/kratos-layout/internal/app/config"
	"github.com/northes/kratos-layout/internal/app/transport/http/handler"
	"github.com/northes/kratos-layout/internal/app/transport/http/router"
)

var ProviderSet = wire.NewSet(
	handler.ProviderSet,
	router.ProviderSet,
	NewServer,
)

func NewServer(
	logger log.Logger,
	httpConf *config.HTTP,
	router *gin.Engine,
) *khttp.Server {
	if router == nil {
		return nil
	}
	var opts = []khttp.ServerOption{
		khttp.Logger(logger),
	}

	if httpConf.Network != "" {
		opts = append(opts, khttp.Network(httpConf.Network))
	}

	if httpConf.Addr != "" {
		opts = append(opts, khttp.Address(httpConf.Addr))
	}

	if httpConf.Timeout != 0 {
		opts = append(opts, khttp.Timeout(time.Duration(httpConf.Timeout)*time.Second))
	}

	srv := khttp.NewServer(opts...)
	srv.HandlePrefix("/", router)

	return srv
}
