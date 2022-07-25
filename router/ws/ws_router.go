/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 10:55
 */

package ws

import "github.com/warlock-backend/server/ws"

// WebsocketInit Websocket 路由
func WebsocketInit() {
	ws.Register("login", ws.LoginController)
	ws.Register("create_player", ws.CreatePlayer)
	ws.Register("move_to", ws.MoveTo)
	ws.Register("shoot_fireball", ws.ShootFireball)
	ws.Register("attack", ws.Attack)
	ws.Register("blink", ws.Blink)
	ws.Register("message", ws.Message)
	ws.Register("heartbeat", ws.HeartbeatController)
}
