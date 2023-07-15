package dal

import (
	"context"
	"github.com/practice/opentelemetry-practice/pkg/server/middleware"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

// GetOrderExtraInfo 模拟对db操作
func GetOrderExtraInfo(parentCtx context.Context) {

	_, span := middleware.TraceProvider.Tracer(middleware.TracerName).Start(parentCtx, "order-extrainfo")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{
		Key: "备注", Value: attribute.StringValue("获取订单额外信息"),
	})

	time.Sleep(time.Second * 2)
}

// UpdateOrderState 模拟对db操作
func UpdateOrderState(parentCtx context.Context) {
	_, span := middleware.TraceProvider.Tracer(middleware.TracerName).Start(parentCtx, "order-update-status")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{
		Key: "备注", Value: attribute.StringValue("更新订单状态"),
	})
	time.Sleep(time.Second * 1) //假设这个是业务函数
}
