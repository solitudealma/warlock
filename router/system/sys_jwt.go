/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:44
 */

package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/solitudealma/warlock/api/v1"
)

type JwtRouter struct{}

func (s *JwtRouter) InitJwtRouter(Router *gin.RouterGroup) {
	jwtRouter := Router.Group("settings")
	jwtApi := v1.ApiGroupApp.SystemApiGroup.JwtApi
	{
		jwtRouter.GET("logout", jwtApi.JsonInBlacklist) // jwt加入黑名单
	}
}
