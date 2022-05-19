package repo

import (
	"github.com/google/wire"
	"github.com/northes/kratos-layout/internal/app/repo/user"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(wire.Bind(new(user.IRepo), new(*user.Repository)), user.NewRepository),
)
