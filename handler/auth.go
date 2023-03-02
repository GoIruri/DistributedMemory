package handler

import (
	"DistributedMemory/common"
	"DistributedMemory/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HTTPInterceptor http请求拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(username) < 3 || !IsTokenValid(username, token) {
			// token校验失败则跳转到直接返回失败提示
			c.Abort()
			resp := util.NewRespMsg(int(common.StatusParamInvalid),
				"token 无效",
				nil)
			c.JSON(http.StatusForbidden, resp)
			return
		}

		c.Next()
	}
}
