package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/server/dal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"net/url"
)

func UserInfoAndScore(c *gin.Context) {
	id := c.Param("id")

	// 模拟请求其他接口
	score, _ := requestForMap(c.Request.Context(), "/users/score/"+id)
	info, _ := requestForMap(c.Request.Context(), "/users/info/"+id)
	c.JSON(200, gin.H{"info": info, "score": score})
}

func UserScore(c *gin.Context) {
	fmt.Println(c.Request.Header)
	c.JSON(200, gin.H{"userid": c.Param("id"), "socre": 100})
}

func UserInfo(c *gin.Context) {

	id := c.Param("id")

	c.JSON(200, gin.H{"userid": c.Param("id"), "name": "user-" + id})
}

func Order(c *gin.Context) {

	if c.Query("error") != "" {
		// 传入error，jaeger中会显示日志
		span := trace.SpanFromContext(c.Request.Context())
		span.RecordError(fmt.Errorf("订单错误信息"))
		c.String(400, "订单错误")
		return
	}

	// 子方法，用来获取子业务信息
	dal.GetOrderExtraInfo(c.Request.Context())
	dal.UpdateOrderState(c.Request.Context())

	c.String(200, "订单列表")
}

const LocalHost = "http://localhost:8080"

// requestForMap 模拟请求其他接口时，记录下header trace
func requestForMap(ctx context.Context, reqUrl string) (gin.H, error) {

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

	// trace 记录在header中，实现不同请求的链路调用
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// http请求
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
