/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:19
 */

package config

type Server struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	JWTRedis JWTRedis `mapstructure:"jwt_redis" json:"jwt_redis" yaml:"jwt_redis"`
	WSRedis  WSRedis  `mapstructure:"ws_redis" json:"ws_redis" yaml:"ws_redis"`
	System   System   `mapstructure:"system" json:"system" yaml:"system"`

	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}
