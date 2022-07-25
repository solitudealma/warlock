package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/warlock-backend/core"
	"github.com/warlock-backend/global"
	"github.com/warlock-backend/initialize"
	"github.com/warlock-backend/router/ws"
	"github.com/warlock-backend/server/task"
	websocket "github.com/warlock-backend/server/ws"
	"github.com/warlock-backend/service/system"
	"time"
)

func main() {
	global.WlVp = core.Viper()
	global.WlLog = core.Logrus() // 初始化Logrus日志库
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(global.WlConfig.JWT.ExpiresTime)),
	)
	gin.SetMode(gin.ReleaseMode) //开发环境
	if global.WlConfig.System.UseMultipoint || global.WlConfig.System.UseRedis {
		// 初始化redis服务
		initialize.JWTRedis()
	}

	initialize.WSRedis()
	ws.WebsocketInit()
	// 定时任务
	task.Init()
	go websocket.StartWebSocket()
	// 从db加载jwt数据
	system.LoadAll()

	router := initialize.Routers()
	port := fmt.Sprintf(":%d", global.WlConfig.System.Addr)
	err := router.Run(port)
	if err != nil {
		return
	}
}
