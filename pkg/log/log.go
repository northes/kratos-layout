package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var (
	DefaultLevel      = Info
	DefaultFormat     = Json
	DefaultWriter     = os.Stdout
	DefaultCallerSkip = 0
)

// Logger 日志
type Logger struct {
	Level      Level
	Format     Format
	Writer     io.Writer
	CallerSkip int
}

func New(options ...Option) *zap.Logger {
	logger := &Logger{
		Level:      DefaultLevel,
		Format:     DefaultFormat,
		Writer:     DefaultWriter,
		CallerSkip: DefaultCallerSkip,
	}

	for _, option := range options {
		option(logger)
	}

	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		writerSyncer  zapcore.WriteSyncer
	)

	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	switch logger.Format {
	case Text:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case Json:
		fallthrough // 继续执行
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	writerSyncer = zapcore.AddSync(logger.Writer)

	core := zapcore.NewCore(
		encoder,
		writerSyncer,
		logger.Level.Convert(),
	)

	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(logger.CallerSkip),
	)

	return zapLogger
}

type Option func(logger *Logger)

// WithLevel 设置日志等级
func WithLevel(level Level) Option {
	return func(logger *Logger) {
		logger.Level = level
	}
}

// WithFormat 设置日志输出的格式
func WithFormat(format Format) Option {
	return func(logger *Logger) {
		logger.Format = format
	}
}

// WithWriter 设置日志输出的 Writer
func WithWriter(writer io.Writer) Option {
	return func(logger *Logger) {
		logger.Writer = writer
	}
}

// WithCallerSkip 设置日志跳过的层级
func WithCallerSkip(skip int) Option {
	return func(logger *Logger) {
		logger.CallerSkip = skip
	}
}

// Format 日志格式
type Format string

const (
	Text Format = "text"
	Json Format = "json"
)

// Level 日志等级
type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Panic Level = "panic"
	Fatal Level = "fatal"
)

func (l Level) Convert() zapcore.Level {
	switch l {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	case Panic:
		return zapcore.PanicLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
