/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:52
 */

package initialize

import (
	"context"

	"github.com/warlock-backend/global"

	"github.com/go-redis/redis/v8"
)

func JWTRedis() {
	redisCfg := global.WlConfig.JWTRedis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.WlLog.Errorf("redis connect ping failed, err:", err)
	} else {
		global.WlLog.Infof("redis connect ping response: %s", pong)
		global.WlJWTRedis = client
	}
}

func WSRedis() {
	redisCfg := global.WlConfig.WSRedis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.WlLog.Errorf("redis connect ping failed, err: %v", err)
	} else {
		global.WlLog.Infof("redis connect ping response: %s", pong)
		global.WlWSRedis = client
	}
}
