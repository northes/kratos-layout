package component

import (
	"github.com/google/wire"

	"github.com/northes/kratos-layout/internal/app/component/discovery"
	"github.com/northes/kratos-layout/internal/app/component/orm"
	"github.com/northes/kratos-layout/internal/app/component/redis"
)

var ProviderSet = wire.NewSet(
	wire.NewSet(discovery.New),
	wire.NewSet(orm.New),
	wire.NewSet(redis.New),
)
