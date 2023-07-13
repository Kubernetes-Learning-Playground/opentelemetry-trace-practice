package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/server/handler"
	"github.com/practice/opentelemetry-practice/pkg/server/middleware"
	"go.opentelemetry.io/otel/trace"
)

func HttpServer() {

	r := gin.New()

	// 使用中间件的方式引入链路追踪
	r.Use(middleware.OpenTelemetryTraceMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "测试用")
	})

	// GET  /users/1101  --- 聚合API
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		score, _ := handler.RequestForMap(c.Request.Context(), "/users/score/"+id)
		info, _ := handler.RequestForMap(c.Request.Context(), "/users/info/"+id)
		c.JSON(200, gin.H{"info": info, "score": score})
	})

	// 子API
	r.GET("/users/score/:id", func(c *gin.Context) {
		fmt.Println(c.Request.Header)
		c.JSON(200, gin.H{"userid": c.Param("id"), "socre": 100})
	})

	// 子API
	r.GET("/users/info/:id", func(c *gin.Context) {

		id := c.Param("id")

		c.JSON(200, gin.H{"userid": c.Param("id"), "name": "user-" + id})
	})

	r.GET("/orders", func(c *gin.Context) {

		if c.Query("error") != "" {
			span := trace.SpanFromContext(c.Request.Context())
			span.RecordError(fmt.Errorf("订单错误信息"))
			c.String(400, "订单错误")
			return
		}

		handler.GetOrderExtraInfo(c.Request.Context()) // 好比是子方法，用来获取子业务信息
		handler.UpdateOrderState(c.Request.Context())

		c.String(200, "订单列表")
	})

	r.Run(":8080")
}
