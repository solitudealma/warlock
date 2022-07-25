/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 20:46
 */

package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/warlock-backend/global"
	"github.com/warlock-backend/model/ws"
	"github.com/warlock-backend/service"
	"time"
)

func getRoomNameFromRooms(roomName string) bool {
	var flag = false
	for k := range clientManager.Rooms {
		if k == roomName {
			flag = true
		}
	}
	global.WlLog.Info(fmt.Sprintf("%v", flag))
	return flag
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	loginInfo := &ws.Login{}
	if err := json.Unmarshal(message, loginInfo); err != nil {
		code = ws.ParameterIllegal
		global.WlLog.Errorf("用户登录 解析数据失败 请求id: "+seq+"err: %v", err)

		return
	}

	//fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", loginInfo.ServiceToken)
	//fmt.Printf("%+v %d\n", loginInfo, len(loginInfo.UserId))
	if loginInfo.UserId == "" || len(loginInfo.UserId) >= 40 {
		code = ws.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, loginInfo.UserId)

		return
	}

	if !InAppIds(loginInfo.AppId) {
		code = ws.Unauthorized
		fmt.Println("用户登录 不支持的平台", seq, loginInfo.AppId)

		return
	}

	client.Login(loginInfo.AppId, loginInfo.UserId, loginInfo.Username, loginInfo.Photo, currentTime)

	// 存储数据
	userOnline := ws.UserLogin(serverIp, serverPort, loginInfo.AppId, loginInfo.UserId, loginInfo.Username,
		loginInfo.Photo, client.Addr, currentTime)
	err := service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)

		return
	}

	// 用户登录
	login := &Login{
		AppId:  loginInfo.AppId,
		UserId: loginInfo.UserId,
		Client: client,
	}
	clientManager.Login <- login

	//fmt.Println("用户登录 成功", seq, client.Addr, loginInfo.UserId)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	heartBeatInfo := &ws.HeartBeat{}
	if err := json.Unmarshal(message, heartBeatInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	//fmt.Println("webSocket_request 心跳接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("心跳接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)

			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	return
}

func CreatePlayer(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	playerInfo := &ws.CreatePlayer{}
	if err := json.Unmarshal(message, playerInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("创建玩家接口 解析数据失败", seq, err)

		return
	}

	//fmt.Println("webSocket_request 创建玩家接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("创建玩家接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("创建玩家接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("创建玩家接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	var (
		roomName = ""
	)
	for i := 1; i < 100000000; i++ {
		name := fmt.Sprintf("room-%d", i)
		fmt.Printf("%d", len(clientManager.Rooms[name]))
		if !getRoomNameFromRooms(name) || len(clientManager.Rooms[name]) < RoomCapacity {
			roomName = name
			break
		}
	}

	if roomName == "" {
		return
	}

	for _, value := range clientManager.Rooms[roomName] {
		client.Send <- []byte(ws.GetCreatePlayerMsgData(value.UserId, value.Username, value.Photo, seq))
	}

	client.RoomName = roomName
	clientManager.AddToRooms(roomName, client)

	player := &ws.CreatePlayer{UUid: playerInfo.UUid, Username: playerInfo.Username, Photo: playerInfo.Photo, AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "create_player",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    player,
		},
	})

	fmt.Printf("rooms length: %d\n", len(clientManager.Rooms[client.RoomName]))
	fmt.Printf("房间 %s 有：", client.RoomName)
	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}
	return
}

func MoveTo(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	moveToInfo := &ws.MoveTo{}
	if err := json.Unmarshal(message, moveToInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家移动接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("moveToInfo: %+v\n", moveToInfo)
	//fmt.Println("webSocket_request 玩家移动接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家移动接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("玩家移动接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("玩家移动接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	moveTo := &ws.MoveTo{UUid: moveToInfo.UUid, Tx: moveToInfo.Tx, Ty: moveToInfo.Ty, AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "move_to",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    moveTo,
		},
	})

	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}

	return
}

func ShootFireball(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	shootFireBallInfo := &ws.ShootFireball{}
	if err := json.Unmarshal(message, shootFireBallInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家发射火球接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("shootFireBallInfo: %+v\n", shootFireBallInfo)
	//fmt.Println("webSocket_request 玩家移动接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家发射火球接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("玩家发射火球接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("玩家发射火球接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	shootFireBall := &ws.ShootFireball{UUid: shootFireBallInfo.UUid, Tx: shootFireBallInfo.Tx, Ty: shootFireBallInfo.Ty,
		BallUuid: shootFireBallInfo.BallUuid, AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "shoot_fireball",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    shootFireBall,
		},
	})

	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}

	return
}

func Attack(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	attackInfo := &ws.Attack{}
	if err := json.Unmarshal(message, attackInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家被攻击接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("attackInfo: %+v\n", attackInfo)
	//fmt.Println("webSocket_request 玩家移动接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家被攻击接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("玩家被攻击接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("玩家被攻击接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	attack := &ws.Attack{UUid: attackInfo.UUid, X: attackInfo.X, Y: attackInfo.Y, AttackedUuid: attackInfo.AttackedUuid,
		BallUuid: attackInfo.BallUuid, Angle: attackInfo.Angle, Damage: attackInfo.Damage, AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "attack",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    attack,
		},
	})

	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}

	return
}

func Blink(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	blinkInfo := &ws.Blink{}
	if err := json.Unmarshal(message, blinkInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家闪现接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("blinkInfo: %+v\n", blinkInfo)
	//fmt.Println("webSocket_request 玩家闪现接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家闪现接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("玩家闪现接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("玩家闪现接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	blink := &ws.Blink{UUid: blinkInfo.UUid, Tx: blinkInfo.Tx, Ty: blinkInfo.Ty, AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "blink",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    blink,
		},
	})

	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}

	return
}

func Message(client *Client, seq string, message []byte) (code int, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	messageInfo := &ws.ChatMessage{}
	if err := json.Unmarshal(message, messageInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家发送聊天消息接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("messageInfo: %+v\n", messageInfo)
	//fmt.Println("webSocket_request 玩家发送聊天消息接口", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家发送聊天消息接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := service.ServiceGroupApp.SystemServiceGroup.WSService.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = ws.NotLoggedIn
			fmt.Println("玩家发送聊天消息接口 用户未登录", seq, client.AppId, client.UserId)

			return
		} else {
			code = ws.ServerError
			fmt.Println("玩家发送聊天消息接口 GetUserOnlineInfo", seq, client.AppId, client.UserId, err)
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = service.ServiceGroupApp.SystemServiceGroup.WSService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	chatMessage := &ws.ChatMessage{UUid: messageInfo.UUid, Text: messageInfo.Text, Username: messageInfo.Username,
		AppId: 101}

	textData, _ := json.Marshal(gin.H{
		"seq":   seq,
		"type":  "group_send_event",
		"event": "message",
		"response": gin.H{
			"code":    200,
			"codeMsg": "success",
			"data":    chatMessage,
		},
	})

	for _, c := range clientManager.Rooms[client.RoomName] {
		fmt.Printf("%+v\n", c)
		c.Send <- textData
	}

	return
}
