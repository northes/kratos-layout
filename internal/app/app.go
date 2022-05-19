package app

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/northes/kratos-layout/internal/app/component"
	"github.com/northes/kratos-layout/internal/app/config"
	"github.com/northes/kratos-layout/internal/app/model"
	"github.com/northes/kratos-layout/internal/app/repo"
	"github.com/northes/kratos-layout/internal/app/service"
	"github.com/northes/kratos-layout/internal/app/transport"
)

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	component.ProviderSet,
	transport.ProviderSet,
	repo.ProviderSet,
	service.ProviderSet,
)

type App struct {
	logger    *log.Helper
	db        *gorm.DB
	transport *transport.Transport
}

func New(
	logger log.Logger,
	db *gorm.DB,
	transport *transport.Transport,
) *App {
	return &App{
		logger:    log.NewHelper(logger),
		db:        db,
		transport: transport,
	}
}

// Start 启动应用
func (a *App) Start(cancel context.CancelFunc) (err error) {
	// 数据迁移
	if a.db != nil {
		if err = model.Migrate(a.db); err != nil {
			return
		}

		a.logger.Info("database migration completed")
	}

	go func() {
		if err = a.transport.Start(); err != nil {
			a.logger.Error(err)
			cancel()
			return
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) (err error) {
	if err = a.transport.Stop(); err != nil {
		return
	}
	return nil
}
