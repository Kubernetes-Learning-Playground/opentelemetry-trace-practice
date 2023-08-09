package main

import (
//"context"
//"time"

// "go.opentelemetry.io/otel"
// "go.opentelemetry.io/otel/metric"
// "go.opentelemetry.io/otel/sdk/export/metric/prometheus"
// "go.opentelemetry.io/otel/sdk/metric/controller/push"
// "go.opentelemetry.io/otel/sdk/metric/processor/basic"
// "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func main() {
	//// 创建一个 MeterProvider
	//meterProvider := otel.GetMeterProvider()
	//
	//// 创建一个 Meter
	//meter := meterProvider.Meter("example-meter")
	//
	//// 创建一个 Counter 指标
	//counter := metric.Must(meter).NewInt64Counter("example-counter")
	//
	//// 创建一个 Prometheus Exporter
	//exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//defer exporter.Stop()
	//
	//// 创建一个 BatchRecorder，用于批量记录指标
	//batcher := basic.New(
	//	simple.NewWithExactDistribution(),
	//	push.WithPeriod(10*time.Second),
	//	push.WithExporter(exporter),
	//)
	//pusher := push.New(batcher, nil)
	//pusher.Start()
	//defer pusher.Stop()
	//
	//// 记录指标
	//ctx := context.Background()
	//counter.Add(ctx, 1)
}
