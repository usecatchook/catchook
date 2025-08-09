package tracer

import (
	"context"
	"strings"
	"time"

	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var tracerProvider *sdktrace.TracerProvider
var mainTracer trace.Tracer
var isEnabled bool

func Initialize(cfg config.TracerConfig, appLogger logger.Logger) error {
	if !cfg.Enabled {
		isEnabled = false
		appLogger.Info(context.Background(), "tracing disabled, using no-op tracer")
		mainTracer = noop.NewTracerProvider().Tracer("no-op")
		return nil
	}

	isEnabled = true
	httpClientOpts := []otlptracehttp.Option{}
	if cfg.Endpoint != "" {
		if hasScheme(cfg.Endpoint) {
			httpClientOpts = append(httpClientOpts, otlptracehttp.WithEndpointURL(cfg.Endpoint))
		} else {
			httpClientOpts = append(httpClientOpts, otlptracehttp.WithEndpoint(cfg.Endpoint))
		}
	}
	if cfg.Insecure {
		httpClientOpts = append(httpClientOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(httpClientOpts...),
	)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
			),
		),
	)

	otel.SetTracerProvider(tp)
	tracerProvider = tp
	mainTracer = tp.Tracer(cfg.ServiceName)

	appLogger.Info(context.Background(), "tracing initialized",
		logger.String("service", cfg.ServiceName),
		logger.String("endpoint", cfg.Endpoint),
	)
	return nil
}

// StartSpan starts a new span with the given context
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// If tracing is not enabled, return a noop span from context
	if !isEnabled {
		return ctx, trace.SpanFromContext(ctx)
	}

	// Add request_id attribute if present in context
	if requestID, ok := ctx.Value(logger.RequestIDKey).(string); ok {
		opts = append(opts, trace.WithAttributes(
			attribute.String("request_id", requestID),
		))
	}
	return mainTracer.Start(ctx, name, opts...)
}

func Close(ctx context.Context) error {
	if !isEnabled || tracerProvider == nil {
		return nil
	}
	return tracerProvider.Shutdown(ctx)
}

func hasScheme(endpoint string) bool {
	return strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")
}

func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()

	if err := fn(ctx); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
