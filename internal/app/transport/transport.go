package transport

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"github.com/northes/kratos-layout/internal/app/component/discovery"
	"github.com/northes/kratos-layout/internal/app/config"
	gtr "github.com/northes/kratos-layout/internal/app/transport/grpc"
	htr "github.com/northes/kratos-layout/internal/app/transport/http"
)

var hostname, _ = os.Hostname()

var ProviderSet = wire.NewSet(
	gtr.ProviderSet,
	htr.ProviderSet,
	New,
)

type Transport struct {
	logger *log.Helper
	server *kratos.App
}

func New(
	logger log.Logger,
	appConf *config.App,
	hs *http.Server,
	gs *grpc.Server,
	discovery discovery.Discovery,
) *Transport {
	var servers []transport.Server
	if hs != nil {
		servers = append(servers, hs)
	}
	if gs != nil {
		servers = append(servers, gs)
	}

	options := []kratos.Option{
		kratos.ID(hostname),
		kratos.Name(appConf.Name),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(servers...),
	}

	if discovery != nil {
		options = append(options, kratos.Registrar(discovery))
	}

	server := kratos.New(options...)

	return &Transport{
		logger: log.NewHelper(logger),
		server: server,
	}
}

func (t *Transport) Start() error {
	t.logger.Infof("transport server starting ...")

	if err := t.server.Run(); err != nil {
		return err
	}
	return nil
}

func (t *Transport) Stop() error {
	if err := t.server.Stop(); err != nil {
		return err
	}

	t.logger.Info("transport server stopping ...")
	return nil
}
