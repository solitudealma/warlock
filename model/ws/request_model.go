/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 0:09
 */

package ws

/************************  请求数据  **************************/

// Request 通用请求数据格式
type Request struct {
	Seq   string      `json:"seq"`            // 消息的唯一Id
	Event string      `json:"event"`          // 请求命令字
	Data  interface{} `json:"data,omitempty"` // 数据 json
}

// PlayerLogin 登录请求数据
type PlayerLogin struct {
	ServiceToken string `json:"serviceToken"` // 验证用户是否登录
	AppId        uint32 `json:"appId"`
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	Photo        string `json:"photo"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}

type Player struct {
	Uuid        string `json:"uuid"`
	Username    string `json:"username"`
	Photo       string `json:"photo"`
	Score       uint32 `json:"score,omitempty"`
	Hp          uint32 `json:"hp"`
	WaitingTime uint32 `json:"waiting_time,omitempty"`
}

type CreatePlayer struct {
	AppId    uint32 `json:"appId"`
	UUid     string `json:"uuid"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
	Hp       uint32 `json:"hp"`
}

type MoveTo struct {
	AppId uint32  `json:"appId"`
	UUid  string  `json:"uuid"`
	Tx    float64 `json:"tx"`
	Ty    float64 `json:"ty"`
}

type ShootFireball struct {
	AppId    uint32  `json:"appId"`
	UUid     string  `json:"uuid"`
	Tx       float64 `json:"tx"`
	Ty       float64 `json:"ty"`
	BallUuid string  `json:"ball_uuid"`
}

type Attack struct {
	AppId        uint32  `json:"appId"`
	UUid         string  `json:"uuid"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	AttackedUuid string  `json:"attacked_uuid"`
	Angle        float64 `json:"angle"`
	Damage       float64 `json:"damage"`
	BallUuid     string  `json:"ball_uuid"`
}

type Blink struct {
	AppId uint32  `json:"appId"`
	UUid  string  `json:"uuid"`
	Tx    float64 `json:"tx"`
	Ty    float64 `json:"ty"`
}

type ChatMessage struct {
	AppId    uint32 `json:"appId"`
	UUid     string `json:"uuid"`
	Username string `json:"username"`
	Text     string `json:"text"`
}
