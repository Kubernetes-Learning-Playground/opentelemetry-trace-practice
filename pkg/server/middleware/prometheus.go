package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector prometheus收集器
var MetricsCollector *PrometheusCollector

func init() {
	MetricsCollector = NewPrometheusCollector()
}

// PrometheusCollector prometheus收集器
type PrometheusCollector struct {
	VisitCounterVec   *prometheus.CounterVec
	OrderCounterVec   *prometheus.CounterVec
	RequestCounterVec *prometheus.CounterVec
	GaugeVec          *prometheus.GaugeVec
}

// NewPrometheusCollector prometheus collector
func NewPrometheusCollector() *PrometheusCollector {
	return &PrometheusCollector{
		VisitCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "visit_counter",
			Help: "The total number of processed events",
		}, []string{"userid"}), // 这里必须写全部的
		RequestCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "request_counter",
			Help: "The total number of processed events",
		}, []string{"method", "path", "statuscode"}),
		OrderCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "order_counter",
			Help: "The total number of processed events",
		}, []string{"ordername"}),
		GaugeVec: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "opentelemetry_prometheus_gauge",
			Help: "The total number of processed events",
		}, []string{}),
	}
}

// Metrics 使用中间件记录请求成功的纪录
func (pm *PrometheusCollector) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {

		pm.RequestCounterVec.With(
			prometheus.Labels{
				"method":     c.Request.Method,
				"path":       c.Request.URL.String(),
				"statuscode": "200",
			},
		).Inc()

		c.Next()
	}
}
