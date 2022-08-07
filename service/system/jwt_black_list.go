/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:37
 */

package system

import (
	"context"
	"fmt"
	"github.com/solitudealma/warlock/config"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/system"
	"time"
)

type JwtService struct{}

//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	querySql := "INSERT INTO jwt_blacklists(`created_at`, `updated_at`, `jwt`) VALUES(?, ?, ?)"
	//currentTime := time.Now().Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	_, err = config.MasterDB.Exec(querySql, currentTime, currentTime, jwtList.Jwt)
	if err != nil {
		fmt.Println(err)
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.WlJWTRedis.Get(context.Background(), userName).Result()
	return redisJWT, err
}

//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	global.WlLog.Info("SetRedisJWT")
	timer := time.Duration(global.WlConfig.JWT.ExpiresTime) * time.Second
	err = global.WlJWTRedis.Set(context.Background(), userName, jwt, timer).Err()
	global.WlLog.Errorf("SetRedisJWT err: %v", err)
	return err
}

func LoadAll() {
	var data []string
	querySql := "SELECT `jwt` FROM jwt_blacklists"
	err := config.MasterDB.Select(&data, querySql)
	if err != nil {
		global.WlLog.Errorf("加载数据库jwt黑名单失败%+v\n", err)
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
