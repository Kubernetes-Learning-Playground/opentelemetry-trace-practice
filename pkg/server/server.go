package server

import (
	"github.com/gin-gonic/gin"
	"github.com/practice/opentelemetry-practice/pkg/server/handler"
	"github.com/practice/opentelemetry-practice/pkg/server/middleware"
)

func HttpServer() {

	r := gin.New()

	// 使用中间件的方式引入链路追踪
	r.Use(middleware.OpenTelemetryTraceMiddleware())

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

	r.Run(":8080")
}
