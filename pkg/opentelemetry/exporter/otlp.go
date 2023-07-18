package exporter

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
)

// NewResource 资源：可观测实体
func NewOTLPResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("myotlpweb"),
		),
	)
	return r
}

// NewOTLPExporter 导出器
func NewOTLPExporter() (trace.SpanExporter, error) {
	// 跳过证书
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

// NewOTLProvider otlp-mode提供者
func NewOTLProvider() *trace.TracerProvider {
	exporter, err := NewOTLPExporter()
	if err != nil {
		log.Fatalln(err)
	}
	res := NewOTLPResource()

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp

}
