package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/theotruvelot/catchook/internal/config"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
}

type zapLogger struct {
	*zap.Logger
}

func New(cfg config.LoggerConfig) (Logger, error) {
	config := zap.NewProductionConfig()

	if cfg.Development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Format == "console" {
		config.Encoding = "console"
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &zapLogger{Logger: logger}, nil
}

func (l *zapLogger) getRequestID(ctx context.Context) string {
	if ctx == nil {
		return "no_context"
	}

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return "no_request_id"
}

func (l *zapLogger) log(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	allFields := make([]zap.Field, 0, len(fields)+1)
	allFields = append(allFields, zap.String("request_id", l.getRequestID(ctx)))
	allFields = append(allFields, fields...)

	switch level {
	case zapcore.DebugLevel:
		l.Logger.Debug(msg, allFields...)
	case zapcore.InfoLevel:
		l.Logger.Info(msg, allFields...)
	case zapcore.WarnLevel:
		l.Logger.Warn(msg, allFields...)
	case zapcore.ErrorLevel:
		l.Logger.Error(msg, allFields...)
	case zapcore.FatalLevel:
		l.Logger.Fatal(msg, allFields...)
	}
}

func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.DebugLevel, msg, fields...)
}

func (l *zapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.InfoLevel, msg, fields...)
}

func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.WarnLevel, msg, fields...)
}

func (l *zapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.ErrorLevel, msg, fields...)
}

func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.FatalLevel, msg, fields...)
}

// Fields helpers
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, val interface{}) zap.Field {
	if d, ok := val.(int64); ok {
		return zap.Int64(key, d)
	}
	return zap.Any(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
