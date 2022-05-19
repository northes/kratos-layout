package orm

import (
	"errors"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	zapgorm "moul.io/zapgorm2"

	"github.com/northes/kratos-layout/internal/app/component/orm/mysql"
)

var (
	ErrUnsupportedType = errors.New("unsupported database type")
)

type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	MySQL       Driver = "mysql"
	PostgresSQL Driver = "postgres"
)

type LogLevel string

const (
	Silent LogLevel = "silent"
	Error  LogLevel = "error"
	Warn   LogLevel = "warn"
	Info   LogLevel = "info"
)

func (l LogLevel) Convert() logger.LogLevel {
	switch l {
	case Silent:
		return logger.Silent
	case Error:
		return logger.Error
	case Warn:
		return logger.Warn
	case Info:
		return logger.Info
	default:
		return logger.Info
	}
}

type Config struct {
	Driver          Driver
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Options         []string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
	LogLevel        LogLevel
}

func New(config *Config, kLogger klog.Logger, zLogger *zap.Logger) (db *gorm.DB, cleanup func(), err error) {
	if config == nil {
		return nil, func() {}, nil
	}

	gLogger := zapgorm.New(zLogger.WithOptions(zap.AddCallerSkip(3)))
	gLogger.SetAsDefault()

	db, err = mysql.New(mysql.Config{
		Driver:                    config.Driver.String(),
		Host:                      config.Host,
		Port:                      config.Port,
		Database:                  config.Database,
		Username:                  config.Username,
		Password:                  config.Password,
		Options:                   config.Options,
		MaxIdleConn:               config.MaxIdleConn,
		MaxOpenConn:               config.MaxOpenConn,
		ConnMaxIdleTime:           time.Second * time.Duration(config.ConnMaxIdleTime),
		ConnMaxLifeTime:           time.Second * time.Duration(config.ConnMaxLifeTime),
		Logger:                    gLogger.LogMode(config.LogLevel.Convert()),
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         0,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
	})

	cleanup = func() {
		logHelper := klog.NewHelper(kLogger)
		logHelper.Info("closing the database resources")

		sqlDB, err := db.DB()
		if err != nil {
			logHelper.Error(err)
		}

		if err := sqlDB.Close(); err != nil {
			logHelper.Error(err)
		}
	}

	return db, cleanup, nil
}
