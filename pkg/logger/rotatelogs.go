package logger

import (
	"github.com/Roh-Bot/blog-api/internal/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
	"time"
)

func RotateLogsWriteSyncer(config config.Config) (zapcore.WriteSyncer, error) {
	logf, err := rotatelogs.New(
		config.Logger.FilePath,
		rotatelogs.WithLinkName(config.Logger.FilePath),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(logf), err
}
