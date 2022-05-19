package service

import (
	"github.com/google/wire"
	"github.com/northes/kratos-layout/internal/app/service/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(user.IService), new(*user.Service)), user.NewService),
)
