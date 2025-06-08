package logger

import (
	"context"
	"github.com/Roh-Bot/blog-api/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	config config.Logger
	zap    *zap.SugaredLogger
}

func ZapNew(config config.Logger, cores ...zapcore.Core) (*ZapLogger, error) {
	// Combine the cores if provided
	var tee zapcore.Core
	if len(cores) > 0 {
		if config.IsDevelopment {
			stdoutLogger := ConsoleSink()
			stdoutCore := NewZapCore(zapcore.Lock(stdoutLogger), config.Level)
			tee = zapcore.NewTee(append([]zapcore.Core{stdoutCore}, cores...)...)
		} else {
			tee = zapcore.NewTee(cores...)
		}
	}

	z := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).
		Sugar()
	zap.Fields()
	return &ZapLogger{zap: z}, nil
}

func LogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dPanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func CustomEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "Msg",
		LevelKey:       "Level",
		TimeKey:        "Timestamp",
		NameKey:        "logger",
		CallerKey:      "Caller",
		StacktraceKey:  "stacktrace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func NewZapCore(syncer zapcore.WriteSyncer, level string) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(CustomEncoderConfig()),
		syncer,
		LogLevel(level))
}

func NewZapCoreLevelEnabler(syncer zapcore.WriteSyncer, level string) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(CustomEncoderConfig()),
		syncer,
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == LogLevel(level)
		}))
}

func NewZapCoreLevelsEnabler(syncer zapcore.WriteSyncer, level string) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(CustomEncoderConfig()),
		syncer,
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl <= LogLevel(level)
		}))
}

func (z *ZapLogger) Sync() error {
	return z.zap.Sync()
}

func (z *ZapLogger) Infoln(args ...any) {
	z.zap.Infoln(args)
}

func (z *ZapLogger) Errorln(args ...any) {
	z.zap.Infoln(args)
}

func (z *ZapLogger) Infof(template string, args ...any) {
	if len(args) == 0 || args == nil {
		z.zap.Infof(template)
		return
	}
	z.zap.Infof(template, args)
}

func (z *ZapLogger) Errorf(template string, args ...any) {
	if len(args) == 0 || args == nil {
		z.zap.Errorf(template)
		return
	}
	z.zap.Errorf(template, args)
}

func (z *ZapLogger) InfolnWithRequestId(ctx context.Context, args ...any) {
	logger := z.zap.With()
	if requestId, ok := ctx.Value("request_id").(string); ok {
		logger = z.zap.With(zap.String("Request", requestId))
	}
	logger.Infoln(args)
}

func (z *ZapLogger) ErrorlnWithRequestId(ctx context.Context, args ...any) {
	logger := z.zap.With()
	if requestId, ok := ctx.Value("request_id").(string); ok {
		logger = z.zap.With(zap.String("Request", requestId))
	}
	logger.Errorln(args)
}

func (z *ZapLogger) InfofWithRequestId(ctx context.Context, template string, args ...any) {
	logger := z.zap.With()
	if requestId, ok := ctx.Value("request_id").(string); ok {
		logger = z.zap.With(zap.String("Request", requestId))
	}
	if len(args) == 0 || args == nil {
		logger.Infof(template)
		return
	}
	logger.Infof(template, args)
}

func (z *ZapLogger) ErrorfWithRequestId(ctx context.Context, template string, args ...any) {
	logger := z.zap.With()
	if requestId, ok := ctx.Value("request_id").(string); ok {
		logger = z.zap.With(zap.String("Request", requestId))
	}
	if len(args) == 0 || args == nil {
		logger.Errorf(template)
		return
	}
	logger.Errorf(template, args)
}
