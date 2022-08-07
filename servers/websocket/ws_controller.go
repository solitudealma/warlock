/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 20:46
 */

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/system"
	"github.com/solitudealma/warlock/model/ws"
	"github.com/solitudealma/warlock/servers/grpcclient"
	"time"
)

var (
	serverIp   string
	serverPort string
)

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	loginInfo := &ws.PlayerLogin{}
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
	err := wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
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
	WsClientManager.Login <- login

	//fmt.Println("用户登录 成功", seq, grpcclient.Addr, loginInfo.UserId)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	heartBeatInfo := &ws.HeartBeat{}
	if err := json.Unmarshal(message, heartBeatInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	//fmt.Println("webSocket_request 心跳接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	return
}

func GroupSendEvent(client *Client, data []byte) {
	if client.RoomName == "" {
		keys := global.WlWSRedis.Keys(context.Background(), "*"+client.UserId+"*").Val()
		if keys != nil {
			client.RoomName = keys[0]
		}
	}
	for _, c := range WsClientManager.Rooms[client.RoomName] {
		c.Send <- data
	}
}

func CreatePlayer(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	playerInfo := &ws.CreatePlayer{}
	if err := json.Unmarshal(message, playerInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("创建玩家接口 解析数据失败", seq, err)

		return
	}

	//fmt.Println("webSocket_request 创建玩家接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("创建玩家接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	var (
		user     = system.SysUser{}
		username = playerInfo.Username
	)

	user, err = userService.GetUserInfo(username)
	player := &ws.Player{Uuid: playerInfo.UUid, Username: playerInfo.Username, Photo: playerInfo.Photo,
		Score: user.Score}
	grpcclient.AddPlayer(player, playerInfo.AppId)

	return
}

func MoveTo(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	moveToInfo := &ws.MoveTo{}
	if err := json.Unmarshal(message, moveToInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家移动接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("moveToInfo: %+v\n", moveToInfo)
	//fmt.Println("webSocket_request 玩家移动接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家移动接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
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

	GroupSendEvent(client, textData)
	return
}

func ShootFireball(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	shootFireBallInfo := &ws.ShootFireball{}
	if err := json.Unmarshal(message, shootFireBallInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家发射火球接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("shootFireBallInfo: %+v\n", shootFireBallInfo)
	//fmt.Println("webSocket_request 玩家移动接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家发射火球接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
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

	GroupSendEvent(client, textData)

	return
}

func Attack(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	attackInfo := &ws.Attack{}
	if err := json.Unmarshal(message, attackInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家被攻击接口 解析数据失败", seq, err)

		return
	}
	//fmt.Println("webSocket_request 玩家移动接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家被攻击接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = ws.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}
	if client.RoomName == "" {
		return
	}

	playersString := global.WlWSRedis.Get(context.Background(), "warlock:room:info:"+client.RoomName).Val()
	players := make([]*ws.Player, 0)
	err = json.Unmarshal([]byte(playersString), &players)
	if err != nil {
		global.WlLog.Errorf("playersString Unmarshal err: %v", err)
	}
	if len(players) == 0 {
		return
	}

	for i := 0; i < len(players); i++ {
		if players[i].Uuid == attackInfo.AttackedUuid {
			players[i].Hp -= 25
		}
	}

	remainCnt := 0
	for i := 0; i < len(players); i++ {
		if players[i].Hp > 0 {
			remainCnt += 1
		}
	}

	if remainCnt > 1 {
		if client.RoomName != "" {
			playersData, _ := json.Marshal(players)
			global.WlWSRedis.Set(context.Background(), "warlock:room:info:"+client.RoomName, string(playersData), 3600*time.Second)
		}
	} else {
		var score int32 = 0
		for i := 0; i < len(players); i++ {
			if players[i].Hp <= 0 {
				score = -5
			} else {
				score = 10
			}
			err := userService.UpdateUserScore(players[i].Username, score)
			if err != nil {
				return 0, "", nil
			}
		}

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

	GroupSendEvent(client, textData)

	return
}

func Blink(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	blinkInfo := &ws.Blink{}
	if err := json.Unmarshal(message, blinkInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家闪现接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("blinkInfo: %+v\n", blinkInfo)
	//fmt.Println("webSocket_request 玩家闪现接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家闪现接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
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

	GroupSendEvent(client, textData)

	return
}

func Message(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = ws.OK
	currentTime := uint64(time.Now().Unix())

	messageInfo := &ws.ChatMessage{}
	if err := json.Unmarshal(message, messageInfo); err != nil {
		code = ws.ParameterIllegal
		fmt.Println("玩家发送聊天消息接口 解析数据失败", seq, err)

		return
	}
	fmt.Printf("messageInfo: %+v\n", messageInfo)
	//fmt.Println("webSocket_request 玩家发送聊天消息接口", grpcclient.AppId, grpcclient.UserId)

	if !client.IsLogin() {
		fmt.Println("玩家发送聊天消息接口 用户未登录", client.AppId, client.UserId, seq)
		code = ws.NotLoggedIn

		return
	}

	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
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
	err = wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
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

	GroupSendEvent(client, textData)

	return
}
