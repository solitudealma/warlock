/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 15:40
 */

package system

import (
	"github.com/gin-gonic/gin"
	websocket "github.com/warlock-backend/server/ws"
)

type WSApi struct{}

func (ws *WSApi) MultiPlayer(ctx *gin.Context) {
	websocket.UpGraderWs(ctx)
}
