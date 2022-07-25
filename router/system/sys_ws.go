/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 15:38
 */

package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/warlock-backend/api/v1"
)

type WSRouter struct{}

func (ws *WSRouter) InitWSRouter(Router *gin.RouterGroup) {
	wsRouter := Router.Group("wss")
	wsApi := v1.ApiGroupApp.SystemApiGroup.WSApi
	{
		wsRouter.GET("multiplayer", wsApi.MultiPlayer) // 多人联机对战
	}
}
