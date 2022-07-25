/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 11:37
 */

package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/warlock-backend/api/v1"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("settings")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		baseRouter.GET("register", baseApi.Register) // 管理员注册账号
		baseRouter.GET("login", baseApi.Login)
	}
	return baseRouter
}
