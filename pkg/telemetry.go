package pkg

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"os"
	"time"
)

func InitMetrics(ctx context.Context) (*metric.MeterProvider, error) {
	exp, err := otelMetricsExporter(ctx)
	if err != nil {
		return nil, err
	}

	res, err := otelResource(ctx)
	if err != nil {
		return nil, err
	}

	meterProv := otelMetricsStart(exp, res)

	return meterProv, nil
}

func otelMetricsExporter(ctx context.Context) (metric.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("127.0.0.1:4317"),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
	)
}

func otelResource(ctx context.Context) (*resource.Resource, error) {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	return resource.New(ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("serviceName"),
			semconv.ServiceInstanceIDKey.String(host),
			attribute.String("env", "dev"),
		),
	)
}

func otelMetricsStart(exporter metric.Exporter, res *resource.Resource) *metric.MeterProvider {
	meterProv := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	global.SetMeterProvider(meterProv)

	return meterProv
}
