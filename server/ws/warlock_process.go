/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 20:46
 */

package ws

import (
	"encoding/json"
	"fmt"
	"github.com/warlock-backend/global"
	"github.com/warlock-backend/model/common"
	"github.com/warlock-backend/model/ws"
	"sync"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code int, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {

	//global.WL_Log.Info("处理数据:"+client.Addr, zap.String("message", string(message)))

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()

	requestInfo := &ws.Request{}

	err := json.Unmarshal(message, requestInfo)
	if err != nil {
		global.WlLog.Errorf("处理数据 json Unmarshal, err: %v", err)
		client.SendMsg([]byte("发送数据不合法"))

		return
	}

	requestData, err := json.Marshal(requestInfo.Data)
	if err != nil {
		global.WlLog.Errorf("处理数据 json Marshal, err: %v", err)
		client.SendMsg([]byte("处理数据失败"))

		return
	}

	seq := requestInfo.Seq
	event := requestInfo.Event

	var (
		code int
		msg  string
		data interface{}
	)

	// request
	//global.WL_Log.Info("warlock_request", zap.String("request", event+" "+client.Addr))

	// 采用 map 注册的方式 获取对应controller函数 进行响应
	if value, ok := getHandlers(event); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = ws.RoutingNotExist
		global.WlLog.Info("处理数据 路由不存在, " + client.Addr + " " + "event: " + event)
	}

	msg = ws.GetErrorMessage(code, msg)

	responseHead := common.NewResponseHead(seq, event, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)

		return
	}

	if event == "login" {
		client.Send <- headByte
	}

	//fmt.Println("warlock_response send", client.Addr, client.AppId, client.UserId, "event", event, "code", code)

	return
}
