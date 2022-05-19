package discovery

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/northes/kratos-layout/internal/app/component/discovery/etcd"
	"go.uber.org/zap"
)

type Discovery interface {
	registry.Registrar
	registry.Discovery
}

type Config struct {
	Etcd *etcd.Config
}

func New(config *Config, logger *zap.Logger) (Discovery, error) {
	if config == nil {
		return nil, nil
	}
	if config.Etcd != nil {
		return etcd.New(config.Etcd, logger)
	}
	return nil, nil
}
