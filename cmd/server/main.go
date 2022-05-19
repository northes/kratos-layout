package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/pflag"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"

	"github.com/northes/kratos-layout/internal/app/command"
	appconfig "github.com/northes/kratos-layout/internal/app/config"
	"github.com/northes/kratos-layout/pkg/log"
	"github.com/northes/kratos-layout/pkg/path"
)

var (
	appName     = "kratos-layout"
	hostname, _ = os.Hostname()
)

var (
	rootPath      = path.RootPath()
	logPath       string
	logLevel      string
	logFormat     string
	logCallerSkip int
	configPath    string
)

func init() {
	pflag.StringVarP(&logPath, "log.path", "", "logs/%Y%m%d.log", "日志输出路径")
	pflag.StringVarP(&logLevel, "log.level", "", "info", "日志等级（debug、info、warn、error、panic、fatal）")
	pflag.StringVarP(&logFormat, "log.format", "", "json", "日志输出格式（text、json）")
	pflag.IntVarP(&logCallerSkip, "log.caller-skip", "", 4, "日志 caller 跳过层数")
	pflag.StringVarP(&configPath, "config", "f", filepath.Join(rootPath, "etc/config.yaml"), "配置文件路径")

	cobra.OnInitialize(setup)
}

var (
	loggerWriter *rotatelogs.RotateLogs
	logger       klog.Logger
	hLogger      *klog.Helper
	zLogger      *zap.Logger
	config       kconfig.Config
	configModel  = new(appconfig.Config)
)

func main() {
	defer cleanup()

	cmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			hLogger.Info("start app ...")

			// 监听退出信号
			signalCtx, signalStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer signalStop()

			appServ, appCleanup, err := initApp(loggerWriter, logger, zLogger, configModel)
			if err != nil {
				panic(err)
			}
			defer appCleanup()

			// 调用 app 启动钩子
			if err := appServ.Start(signalStop); err != nil {
				panic(err)
			}

			// 等待退出信号
			<-signalCtx.Done()
			signalStop()

			hLogger.Info("the app is shutting down ...")

			// 延时关闭以处理定时任务
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configModel.App.Timeout)*time.Second)
			defer cancel()

			// 关闭应用
			if err := appServ.Stop(ctx); err != nil {
				panic(err)
			}
		},
	}

	command.Setup(cmd, func() (*command.Command, func(), error) {
		return initCommand(loggerWriter, logger, zLogger, configModel)
	})

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

}

func setup() {
	var err error

	if logPath != "" {
		if !filepath.IsAbs(logPath) {
			logPath = filepath.Join(rootPath, logPath)
		}

		loggerWriter, err = rotatelogs.New(
			logPath,
			rotatelogs.WithClock(rotatelogs.Local),
		)
		if err != nil {
			panic(err)
		}
	}

	var writer io.Writer
	if loggerWriter == nil {
		writer = os.Stdout
	} else {
		writer = io.MultiWriter(os.Stdout, loggerWriter)
	}
	zLogger = log.New(
		log.WithLevel(log.Level(log.Level(logLevel))),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(writer),
		log.WithCallerSkip(logCallerSkip),
	)
	logger = klog.With(
		kzap.NewLogger(zLogger),
		"service.id", hostname,
		"service.name", appName,
	)
	hLogger = klog.NewHelper(logger)

	hLogger.Info("initializing resource ...")
	hLogger.Infof("the log output directory: %s", filepath.Dir(logPath))

	if configPath == "" {
		panic("config path is missing")
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(rootPath, configPath)
	}

	hLogger.Infof("load config from: %s", configPath)

	configResources := []kconfig.Source{file.NewSource(configPath)}
	config = kconfig.New(
		kconfig.WithSource(configResources...),
		kconfig.WithLogger(logger),
	)

	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := config.Scan(configModel); err != nil {
		panic(err)
	}

	if err := appconfig.Loaded(logger, config, configModel); err != nil {
		panic(err)
	}

	if !funk.ContainsString(appconfig.SupportedEnvs, configModel.App.Env.String()) {
		panic("unsupported env value: " + configModel.App.Env)
	}

	hLogger.Infof("current env: %s", configModel.App.Env)
}

func cleanup() {
	if hLogger != nil {
		hLogger.Info("resource cleaning ...")
	}

	if config != nil {
		if err := config.Close(); err != nil {
			hLogger.Error(err.Error())
		}
	}

	if loggerWriter != nil {
		if err := loggerWriter.Close(); err != nil {
			panic(err)
		}
	}

	if zLogger != nil {
		if err := zLogger.Sync(); err != nil {
			hLogger.Error(err.Error())
		}
	}
}
