package handler

import (
	"github.com/google/wire"
	userpb "github.com/northes/kratos-layout/internal/app/transport/grpc/api/v1/user"
	"github.com/northes/kratos-layout/internal/app/transport/grpc/handler/v1/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(userpb.UserServer), new(*user.Handler)), user.NewHandler),
)
