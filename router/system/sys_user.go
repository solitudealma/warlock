/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 13:22
 */

package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/solitudealma/warlock/api/v1"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("settings")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.GET("getinfo", baseApi.GetUserInfo) // 获取自身信息
	}
}
