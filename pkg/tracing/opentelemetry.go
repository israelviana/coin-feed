package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func InitTracer() *sdktrace.TracerProvider {
	// Recurso SEM schema -> não conflita com o Default()
	baseRes := resource.NewSchemaless(
		attribute.String("service.name", "coin-feed"),
		attribute.String("service.version", "1.0.0"),
		attribute.String("environment", os.Getenv("APP_ENV")),
	)

	res, err := resource.Merge(resource.Default(), baseRes)
	if err != nil {
		panic("failed to create OTel resource: " + err.Error())
	}

	exporterEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if exporterEndpoint == "" {
		tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res))
		otel.SetTracerProvider(tp)
		Tracer = otel.Tracer("coin-feed")
		return tp
	}

	var exporter sdktrace.SpanExporter
	switch os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") {
	case "http/protobuf":
		exp, err := otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpoint(exporterEndpoint),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			panic("failed to create OTLP HTTP exporter: " + err.Error())
		}
		exporter = exp
	default:
		// **Use WithInsecure()** (não passe TLS com credentials "insecure")
		exp, err := otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpoint(exporterEndpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			panic("failed to create OTLP gRPC exporter: " + err.Error())
		}
		exporter = exp
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.5))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	Tracer = otel.Tracer("coin-feed")
	return tp
}
