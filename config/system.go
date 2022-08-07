/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:42
 */

package config

type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`                                  // 环境值
	HttpPort      int    `mapstructure:"httpPort" json:"httpPort" yaml:"httpPort"`                   // 后端端口值
	WebsocketPort int    `mapstructure:"websocketPort" json:"websocketPort" yaml:"websocketPort"`    // websocket端口值
	RpcPort       int    `mapstructure:"rpcPort" json:"rpcPort" yaml:"rpcPort"`                      // rpc服务端口值
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	LimitCountIP  int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP   int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
}
