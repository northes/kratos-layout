package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(
		new(*Config),
		"App",
		"HTTP",
		"GRPC",
		"JWT",
		"DB",
		"Redis",
		"Discovery",
	),
)

var watchKeys = []string{
	"services.self",
}

func Loaded(hLogger log.Logger, cfg config.Config, conf *Config) error {
	if err := Watch(hLogger, cfg, conf); err != nil {
		return err
	}
	return nil
}

func Watch(hLogger log.Logger, cfg config.Config, conf *Config) error {
	var logger = log.NewHelper(hLogger)

	for _, key := range watchKeys {
		logger.Infof("the config is being watching,key: %s", key)

		if err := cfg.Watch(key, func(s string, value config.Value) {
			logger.Infof("config has changed, key: %s", s)

			if err := cfg.Scan(conf); err != nil {
				logger.Errorf("scan config to model failed,err: %v", err)
			}
		}); err != nil {
			return err
		}
	}

	return nil
}
