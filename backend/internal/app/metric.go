package app

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
)

func NewMetricProvider(ctx context.Context) error {
	exporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithEndpoint("localhost:4317"),
		otlpmetrichttp.WithTimeout(5*time.Second),
		otlpmetrichttp.WithInsecure(),
	)

	_, _ = exporter, err

	return nil
}
