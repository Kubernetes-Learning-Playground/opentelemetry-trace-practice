package exporter

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
)

// NewResource 资源：可观测实体
func NewJaegerResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("myweb"),
			semconv.ServiceVersionKey.String("v1.0.0"),
		),
	)
	return r
}

// NewJaegerExporter 导出器
func NewJaegerExporter(endpoint string) (trace.SpanExporter, error) {
	return jaeger.New(
		// TODO: 配置文件传进来
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)),
	)
}

// NewJaegerProvider jaeger-mode提供者
func NewJaegerProvider(endpoint string) *trace.TracerProvider {
	exporter, err := NewJaegerExporter(endpoint)
	if err != nil {
		log.Fatalln(err)
	}
	res := NewJaegerResource()

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp

}
