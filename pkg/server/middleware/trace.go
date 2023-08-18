package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/opentelemetry/exporter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	TracerName = "gin"
)

var TraceProvider *trace.TracerProvider

// OpenTelemetryTraceMiddleware 中间件
func OpenTelemetryTraceMiddleware(endpoint string) gin.HandlerFunc {
	// 直接使用otel collector sdk
	//TraceProvider = exporter.NewJaegerProvider(endpoint, exporter.ServiceHttp)
	TraceProvider = exporter.NewOTLProvider(exporter.ServiceHttp)
	tracer := TraceProvider.Tracer(TracerName)
	return func(c *gin.Context) {

		// 路由完整路径
		spanName := c.FullPath()
		// 如果是notFound路由
		if spanName == "" {
			spanName = "notFoundRoute-" + c.Request.Method
		}
		ctx := c.Request.Context()
		// 需要把 Propagator 表头加入到 context中
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)) //++
		ctx, span := tracer.Start(ctx, spanName)
		//ctx, span := TraceProvider.Tracer(TracerName).Start(c, spanName)
		defer span.End()
		c.Request = c.Request.WithContext(ctx) // 设置spanContext
		c.Next()
		// handler执行完成，才能拿到status
		status := c.Writer.Status() //ex: 200,400,503,404

		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)

		code, msg := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetStatus(code, msg)

	}
}
