package logs

import (
	"context"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func Tracer(name string) func() {
	exporter, err := stdouttrace.New()
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()), // Always sample to collect all traces
	)
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer(name)

	ctx, span := tracer.Start(context.Background(), "main-span")
	return func() {
		defer span.End()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}

	}
}
