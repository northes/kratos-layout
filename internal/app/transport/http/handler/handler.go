package handler

import (
	"github.com/google/wire"
	"github.com/northes/kratos-layout/internal/app/transport/http/handler/v1/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(user.HandlerInterface), new(*user.Handler)), user.NewHandler),
)
