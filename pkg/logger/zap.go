package logger

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/Roh-Bot/blog-api/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const RequestIDKey = "request_id"

type AsyncZapLogger struct {
	logger     *zap.Logger
	queue      chan map[string]any
	quit       chan struct{}
	dropped    uint64
	batchSize  int
	flushDelay time.Duration
}

// NewAsyncZapLogger creates a non-blocking, async Zap logger.
func ZapNew(cfg config.Logger, cores ...zapcore.Core) (*AsyncZapLogger, error) {
	var tee zapcore.Core

	if cfg.EnableStdout {
		stdoutCore := NewZapCore(StdoutSink(), cfg.Level)
		tee = zapcore.NewTee(append([]zapcore.Core{stdoutCore}, cores...)...)
	} else {
		tee = zapcore.NewTee(cores...)
	}

	z := zap.New(
		tee,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	l := &AsyncZapLogger{
		logger:     z,
		queue:      make(chan map[string]any, cfg.BufferSize),
		quit:       make(chan struct{}),
		batchSize:  cfg.BatchSize,  // configurable for high throughput
		flushDelay: cfg.FlushDelay, // flush interval for batching
	}

	go l.worker()
	return l, nil
}

// NewZapCore helper for standardized encoding and level parsing
func NewZapCore(syncer zapcore.WriteSyncer, level string) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig()),
		zapcore.Lock(syncer),
		parseLevel(level),
	)
}

func encoderConfig() zapcore.EncoderConfig {
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

func parseLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// --- Worker and async internals ---

func (l *AsyncZapLogger) worker() {
	defer func() {
		if r := recover(); r != nil {
			l.logger.Error("logger worker panic", zap.Any("recover", r))
		}
	}()

	ticker := time.NewTicker(l.flushDelay)
	defer ticker.Stop()

	var batch []map[string]any

	for {
		select {
		case entry := <-l.queue:
			batch = append(batch, entry)
			if len(batch) >= l.batchSize {
				l.writeBatch(batch)
				batch = nil
			}

			// Optional: log queue saturation
			if len(l.queue) > cap(l.queue)*9/10 {
				l.logger.Warn("log queue nearing capacity",
					zap.Int("current", len(l.queue)),
					zap.Int("capacity", cap(l.queue)))
			}

		case <-ticker.C:
			if len(batch) > 0 {
				l.writeBatch(batch)
				batch = nil
			}

		case <-l.quit:
			// drain remaining logs
			for {
				select {
				case entry := <-l.queue:
					batch = append(batch, entry)
				default:
					if len(batch) > 0 {
						l.writeBatch(batch)
					}
					return
				}
			}
		}
	}
}

func (l *AsyncZapLogger) writeBatch(batch []map[string]any) {
	for _, entry := range batch {
		l.writeEntry(entry)
	}
}

func (l *AsyncZapLogger) writeEntry(entry map[string]any) {
	msg, _ := entry["message"].(string)
	if msg == "" {
		msg = "(empty message)"
	}

	level, _ := entry["level"].(string)
	if level == "" {
		level = "info"
	}

	var fields []zap.Field
	for k, v := range entry {
		if k == "message" || k == "level" {
			continue
		}
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case "debug":
		l.logger.Debug(msg, fields...)
	case "info":
		l.logger.Info(msg, fields...)
	case "warn":
		l.logger.Warn(msg, fields...)
	case "error":
		l.logger.Error(msg, fields...)
	default:
		l.logger.Info(msg, fields...)
	}
}

func (l *AsyncZapLogger) enqueue(entry map[string]any) {
	select {
	case l.queue <- entry:
	default:
		atomic.AddUint64(&l.dropped, 1)
	}
}

// --- Public Logging API ---

func (l *AsyncZapLogger) Info(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, "info", msg, fields)
}

func (l *AsyncZapLogger) Error(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, "error", msg, fields)
}

func (l *AsyncZapLogger) Warn(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, "warn", msg, fields)
}

func (l *AsyncZapLogger) Debug(ctx context.Context, msg string, fields map[string]any) {
	l.log(ctx, "debug", msg, fields)
}

func (l *AsyncZapLogger) log(ctx context.Context, level, msg string, fields map[string]any) {
	if fields == nil {
		fields = map[string]any{}
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields[RequestIDKey] = reqID
	}
	fields["level"] = level
	fields["message"] = msg
	fields["timestamp"] = time.Now().Format(time.RFC3339)
	l.enqueue(fields)
}

// Flush gracefully stops the logger and flushes the remaining entries.
func (l *AsyncZapLogger) Flush() error {
	close(l.quit)
	return l.logger.Sync()
}

// DroppedCount returns the number of dropped log entries.
func (l *AsyncZapLogger) DroppedCount() uint64 {
	return atomic.LoadUint64(&l.dropped)
}

// With allows you to add default fields to the logger.
func (l *AsyncZapLogger) With(fields map[string]any) *AsyncZapLogger {
	newLogger := *l
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	newLogger.logger = l.logger.With(zapFields...)
	return &newLogger
}
