/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 15:42
 */

package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/ws"
)

const (
	userOnlinePrefix    = "warlock:user:online:" // 用户在线状态
	roomPrefix          = "warlock:room:info:"   // 房间信息
	userOnlineCacheTime = 10 * 60
	roomLifeTime        = 30 * 60
)

type WSService struct{}

/*********************  查询用户是否在线  ************************/
func getUserOnlineKey(userKey string) (key string) {
	return fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
}

/*********************  查询用户是否在线  ************************/
func getRoomInfoKey(roomKey string) (key string) {
	return fmt.Sprintf("%s%s", roomPrefix, roomKey)
}

func (wsService *WSService) GetRoomInfo(roomKey string) ([]*ws.CreatePlayer, error) {
	key := getRoomInfoKey(roomKey)
	data, _ := global.WlWSRedis.Get(context.Background(), key).Bytes()
	roomInfo := make([]*ws.CreatePlayer, 0)
	err := json.Unmarshal(data, &roomInfo)
	for _, player := range roomInfo {
		fmt.Printf("getPlayerInfo: %v\n", player)
	}
	return roomInfo, err
}

func (wsService *WSService) SetRoomInfo(roomKey string, roomInfo []*ws.CreatePlayer) (err error) {
	key := getRoomInfoKey(roomKey)

	valueByte, _ := json.Marshal(roomInfo)

	_, err = global.WlWSRedis.Do(context.Background(), "setEx", key, roomLifeTime, string(valueByte)).Result()
	return err
}

func (wsService *WSService) GetUserOnlineInfo(userKey string) (userOnline *ws.UserOnline, err error) {

	key := getUserOnlineKey(userKey)

	data, err := global.WlWSRedis.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", userKey, err)

			return
		}

		fmt.Println("GetUserOnlineInfo", userKey, err)

		return
	}

	userOnline = &ws.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		fmt.Println("获取用户在线数据 json Unmarshal", userKey, err)

		return
	}

	fmt.Println("获取用户在线数据", userKey, "time", userOnline.LoginTime, userOnline.HeartbeatTime, "AccIp", userOnline.AccIp, userOnline.IsLogoff)

	return
}

// SetUserOnlineInfo 设置用户在线数据
func (wsService *WSService) SetUserOnlineInfo(userKey string, userOnline *ws.UserOnline) (err error) {

	key := getUserOnlineKey(userKey)

	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		fmt.Println("设置用户在线数据 json Marshal", key, err)

		return
	}

	_, err = global.WlWSRedis.Do(context.Background(), "setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置用户在线数据 ", key, err)

		return
	}

	return
}
