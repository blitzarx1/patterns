package tracing

import (
	pkgContext "context"
	"time"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/boson-research/patterns/internal/context"
)

// InitTracerProvider creates and globally registers new tracer provider. To use it, call otel.GetTracerProvider().
func InitTracerProvider(
	ctx context.Context,
	serviceName string,
	serviceVersion string,
	agentEndpoint string,
	batchTimeout time.Duration,
) (func(pkgContext.Context) error, error) {
	exporter, err := newExporter(ctx, agentEndpoint)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes([]attribute.KeyValue{
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		}...),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(batchTimeout)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(jaeger.Jaeger{})
	return tp.Shutdown, nil
}

func newExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient([]otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(endpoint),
		}...),
	)

	return exporter, err
}
