package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/server/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// 模拟一个业务函数
func GetOrderExtraInfo(parentCtx context.Context) {

	_, span := middleware.GinTp.Tracer(middleware.TracerName).Start(parentCtx, "order-extrainfo")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{
		Key: "备注", Value: attribute.StringValue("获取订单额外信息"),
	})

	time.Sleep(time.Second * 2)
}

// 模拟 更新订单状态
func UpdateOrderState(parentCtx context.Context) {
	_, span := middleware.GinTp.Tracer(middleware.TracerName).Start(parentCtx, "order-update-status")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{
		Key: "备注", Value: attribute.StringValue("更新订单状态"),
	})
	time.Sleep(time.Second * 1) //假设这个是业务函数
}

// ----------用户相关的演示
const LocalHost = "http://localhost:8080"

func RequestForMap(ctx context.Context, reqUrl string) (gin.H, error) {

	ret := gin.H{}
	u, err := url.Parse(reqUrl)
	if err != nil {
		return ret, err
	}
	if u.Host == "" {
		reqUrl = LocalHost + u.Path
	}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return ret, nil
	}

	otel.GetTextMapPropagator().Inject(ctx,
		propagation.HeaderCarrier(req.Header))

	// go自带的 http client 请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ret, err
	}
	defer rsp.Body.Close()
	b, _ := io.ReadAll(rsp.Body)

	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.Println(err)
		return ret, err
	}
	return ret, nil
}
