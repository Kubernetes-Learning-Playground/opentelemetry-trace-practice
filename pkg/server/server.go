package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/common"
	"github.com/practice/opentelemetry-practice/pkg/server/handler"
	"github.com/practice/opentelemetry-practice/pkg/server/middleware"
)

func HttpServer(c *common.ServerConfig) {

	if !c.Debug{
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 使用中间件的方式引入链路追踪
	r.Use(middleware.OpenTelemetryTraceMiddleware(c.JaegerEndpoint), middleware.MetricsCollector.Metrics())

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "测试用")
	})

	// GET  /users/1101  --- 聚合API
	r.GET("/users/:id", handler.UserInfoAndScore)

	// 子API
	r.GET("/users/score/:id", handler.UserScore)

	// 子API
	r.GET("/users/info/:id", handler.UserInfo)

	r.GET("/orders", handler.Order)

	r.GET("/metrics", handler.PrometheusHandler())

	// 自定义的业务接口：模拟用户的访问量
	r.GET("/users/visit", handler.UserVisit)

	err := r.Run(fmt.Sprintf(":%v", c.Port))
	fmt.Println(err)
}
