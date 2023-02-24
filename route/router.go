package route

import (
	"DistributedMemory/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// gin framework,包括Logger，Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "./static")

	// 不需要经过验证就能访问的接口
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.GET("/user/signin", handler.SignInHandler)
	router.POST("/user/signin", handler.DoSinginHandler)

	// 加入中间件，用于校验token的拦截器
	router.Use(handler.HTTPInterceptor)

	// Use之后的所有handler都会经过拦截器进行token校验
}
