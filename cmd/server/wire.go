//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"

	"github.com/northes/kratos-layout/internal/app"
	"github.com/northes/kratos-layout/internal/app/command"
	appConfig "github.com/northes/kratos-layout/internal/app/config"
)

func initApp(*rotatelogs.RotateLogs, log.Logger, *zap.Logger, *appConfig.Config) (*app.App, func(), error) {
	panic(wire.Build(
		app.ProviderSet,
		app.New,
	))
}

func initCommand(logs *rotatelogs.RotateLogs, logger2 log.Logger, logger3 *zap.Logger, config2 *appConfig.Config) (*command.Command, func(), error) {
	panic(wire.Build(
		command.ProviderSet,
		command.New,
	))
}
