package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/practice/opentelemetry-practice/pkg/opentelemetry/exporter"
	"go.opentelemetry.io/otel"
)

// file-mode模式 链路追踪

const TraceName = "mytrace"

func main() {
	f, err := os.OpenFile("file-mode_trace.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tp := exporter.NewFileProvider(f)
	ctx, span := otel.Tracer(TraceName).Start(context.Background(), "main")

	// 模拟执行业务逻辑
	time.Sleep(time.Second * 3)

	span.End()
	tp.ForceFlush(ctx)
}
