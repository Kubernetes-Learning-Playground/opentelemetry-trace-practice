package exporter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
)

// NewResource 资源：可观测实体
func NewOTLPResource(serviceName string) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
			semconv.ServiceVersionKey.String("v1.20.0"),
		),
	)
	if err != nil {
		fmt.Println("err: ", err)
	}
	return r
}

// NewOTLPExporter 导出器
func NewOTLPExporter() (trace.SpanExporter, error) {
	// 跳过证书，使用http部署
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

// NewOTLProvider otlp-mode提供者
// 使用OTLProvider，可以直接由otel-collector对接所有的蒐集器
func NewOTLProvider(serviceName string) *trace.TracerProvider {
	exporter, err := NewOTLPExporter()
	if err != nil {
		log.Fatalln(err)
	}
	res := NewOTLPResource(serviceName)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}
