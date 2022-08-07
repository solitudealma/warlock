/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/7 18:04
 */

package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/solitudealma/warlock/config"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
)

var (
	WlJWTRedis            *redis.Client
	WlWSRedis             *redis.Client
	WlConfig              config.Server
	WlVp                  *viper.Viper
	WlLog                 *logrus.Logger
	GvaConcurrencyControl = &singleflight.Group{}

	BlackCache local_cache.Cache
)
