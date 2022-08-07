/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 10:55
 */

package ws

import (
	"github.com/solitudealma/warlock/servers/websocket"
)

// WebsocketInit Websocket 路由
func WebsocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("create_player", websocket.CreatePlayer)
	websocket.Register("move_to", websocket.MoveTo)
	websocket.Register("shoot_fireball", websocket.ShootFireball)
	websocket.Register("attack", websocket.Attack)
	websocket.Register("blink", websocket.Blink)
	websocket.Register("message", websocket.Message)
	websocket.Register("heartbeat", websocket.HeartbeatController)
}
