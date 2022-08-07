package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/solitudealma/warlock/config"
	"github.com/solitudealma/warlock/core"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/initialize"
	"github.com/solitudealma/warlock/match_system/grpcserver"
	"github.com/solitudealma/warlock/router/ws"
	"github.com/solitudealma/warlock/servers/task"
	"github.com/solitudealma/warlock/servers/websocket"
	"github.com/solitudealma/warlock/service/system"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"time"
)

func main() {
	global.WlVp = core.Viper()
	global.WlLog = core.Logrus() // 初始化Logrus日志库
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(global.WlConfig.JWT.ExpiresTime)),
	)

	querySql := `
		CREATE TABLE
		IF
			NOT EXISTS jwt_blacklists (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				created_at datetime NULL DEFAULT NULL,
				updated_at datetime NULL DEFAULT NULL,
				deleted_at datetime NULL DEFAULT NULL,
				jwt text NULL
			);
		CREATE TABLE
		IF
			NOT EXISTS sys_users (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				created_at datetime NULL DEFAULT NULL,
				updated_at datetime NULL DEFAULT NULL,
				deleted_at datetime NULL DEFAULT NULL,
				uuid VARCHAR (191) NULL DEFAULT NULL,
				username VARCHAR (191) NULL DEFAULT NULL,
				password VARCHAR (191) NULL DEFAULT NULL, 
				avatar VARCHAR (191) DEFAULT 'https://cdn.acwing.com/media/user/profile/photo/71847_lg_f844104f10.jpg',
				score INTEGER NOT NULL DEFAULT 1500
			)`
	_, err := config.MasterDB.Exec(querySql)
	if err != nil {
		global.WlLog.Errorf("create table err: %v", err)
		return
	}
	gin.SetMode(gin.DebugMode) //开发环境
	if global.WlConfig.System.UseMultipoint || global.WlConfig.System.UseRedis {
		// 初始化redis服务
		initialize.JWTRedis()
	}

	initialize.WSRedis()
	ws.WebsocketInit()
	// 定时任务
	task.Init()
	go websocket.StartWebSocket()
	go grpcserver.Init()
	// 从db加载jwt数据
	system.LoadAll()

	router := initialize.Routers()
	port := fmt.Sprintf(":%d", global.WlConfig.System.HttpPort)
	err = router.Run(port)
	if err != nil {
		return
	}
}
