package handler

import (
	"github.com/northes/kratos-layout/internal/app/command/handler/greet"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(greet.NewHandler)
