package config

import (
	"github.com/northes/kratos-layout/internal/app/component/discovery"
	"github.com/northes/kratos-layout/internal/app/component/orm"
	"github.com/northes/kratos-layout/internal/app/component/redis"
)

var SupportedEnvs = []string{Local.String(), Test.String(), Prod.String()}

type Env string

func (e Env) String() string {
	return string(e)
}

const (
	Local Env = "local"
	Test  Env = "test"
	Prod  Env = "prod"
)

type Config struct {
	App       *App              `json:"app"`
	HTTP      *HTTP             `json:"http"`
	GRPC      *GRPC             `json:"grpc"`
	DB        *orm.Config       `json:"db"`
	Redis     *redis.Config     `json:"redis"`
	Discovery *discovery.Config `json:"discovery"`
	Services  *Services         `json:"services"`
	JWT       *JWT              `json:"jwt"`
}

type App struct {
	Name    string `json:"name"`
	Env     Env    `json:"env"`
	Timeout int64  `json:"timeout"`
}

type HTTP struct {
	Network      string `json:"network"`
	Addr         string `json:"addr"`
	Timeout      int64  `json:"timeout"`
	ExternalAddr string `json:"external_addr"`
}

type GRPC struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Timeout int64  `json:"timeout"`
}

type Services struct {
	Self string `json:"self"`
}

type JWT struct {
	Key string `json:"key"`
}
