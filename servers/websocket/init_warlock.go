/**
* @Author: SolitudeAlma
* @Date: 2022 2022/7/15 20:46
*/

package websocket

import (
    "fmt"
    "github.com/gorilla/websocket"
    "github.com/solitudealma/warlock/global"
    "net/http"
    "time"
)

var (
    WsClientManager = NewClientManager() // 管理者
    appIds          = []uint32{101, 102} // 全部的平台
)

func GetAppIds() []uint32 {

    return appIds
}

func InAppIds(appId uint32) (inAppId bool) {

    for _, value := range appIds {
        if value == appId {
            inAppId = true

            return
        }
    }

    return
}

// StartWebSocket 启动程序
func StartWebSocket() {
    webSocketPort := fmt.Sprintf(":%d", global.WlConfig.System.WebsocketPort)

    http.HandleFunc("/wss/multiplayer", UpGraderWs)
    // 添加处理程序
    go WsClientManager.start()

    err := http.ListenAndServe(webSocketPort, nil)
    if err != nil {
        return
    }
}

func UpGraderWs(w http.ResponseWriter, req *http.Request) {
    // 升级协议
    conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w,
    req, nil)
    if err != nil {
        http.NotFound(w, req)
        return
    }

    fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

    currentTime := uint64(time.Now().Unix())
    client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

    go client.Read()
    go client.Write()

    // 用户连接事件
    WsClientManager.Register <- client
}
