package exporter

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
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
		),
	)
	return r
}

// NewJaegerExporter 导出器
func NewJaegerExporter() (trace.SpanExporter, error) {
	return jaeger.New(
		// TODO: 配置文件传进来
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
}

// NewJaegerProvider jaeger-mode提供者
func NewJaegerProvider() *trace.TracerProvider {
	exporter, err := NewJaegerExporter()
	if err != nil {
		log.Fatalln(err)
	}
	res := NewJaegerResource()

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp

}
