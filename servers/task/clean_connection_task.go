/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 23:33
 */

package task

import (
	"fmt"
	"github.com/solitudealma/warlock/servers/websocket"
	"runtime/debug"
	"time"
)

func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "")

}

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	//fmt.Println("定时任务，清理超时连接", param)

	websocket.ClearTimeoutConnections()

	return
}
