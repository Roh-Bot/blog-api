package logger

import "context"

type Logger interface {
	Infoln(args ...any)
	Infof(template string, args ...any)
	Errorln(args ...any)
	Errorf(template string, args ...any)
	InfolnWithRequestId(ctx context.Context, args ...any)
	InfofWithRequestId(ctx context.Context, template string, args ...any)
	ErrorlnWithRequestId(ctx context.Context, args ...any)
	ErrorfWithRequestId(ctx context.Context, template string, args ...any)
}
