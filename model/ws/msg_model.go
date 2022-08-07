/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 0:05
 */

package ws

const (
	messageTypeText = "text"

	OK                 = 200  // Success
	NotLoggedIn        = 1000 // 未登录
	ParameterIllegal   = 1001 // 参数不合法
	UnauthorizedUserId = 1002 // 非法的用户Id
	Unauthorized       = 1003 // 未授权
	ServerError        = 1004 // 系统错误
	NotData            = 1005 // 没有数据
	ModelAddError      = 1006 // 添加错误
	ModelDeleteError   = 1007 // 删除错误
	ModelStoreError    = 1008 // 存储错误
	OperationFailure   = 1009 // 操作失败
	RoutingNotExist    = 1010 // 路由不存在
)

// Message 消息的定义
type Message struct {
	Target string `json:"target"` // 目标
	Type   string `json:"type"`   // 消息类型 text/img/
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
}

type CreatePlayerMessage struct {
	Event    string `json:"event"`
	UUid     string `json:"uuid"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

func NewTestMsg(from string, Msg string) (message *Message) {
	message = &Message{
		Type: messageTypeText,
		From: from,
		Msg:  Msg,
	}
	return
}

func NewCreatePlayerMsg(uuid string, username string, photo string) (message *CreatePlayerMessage) {
	message = &CreatePlayerMessage{
		Event:    "create_player",
		UUid:     uuid,
		Username: username,
		Photo:    photo,
	}
	return
}

func getCreatePlayerMsgData(event, uuId, username, photo, msgId string) string {
	textMsg := NewCreatePlayerMsg(uuId, username, photo)
	head := NewResponseHead(msgId, event, OK, "success", textMsg)
	return head.String()
}

func getTextMsgData(event, uuId, msgId, message string) string {
	textMsg := NewTestMsg(uuId, message)
	head := NewResponseHead(msgId, event, OK, "Ok", textMsg)
	return head.String()
}

// GetTextMsgData 文本消息
func GetTextMsgData(uuId, msgId, message string) string {
	return getTextMsgData("msg", uuId, msgId, message)
}

// GetCreatePlayerMsgData 用户被创建的消息
func GetCreatePlayerMsgData(uuId, username, photo, msgId string) string {
	return getCreatePlayerMsgData("create_player", uuId, username, photo, msgId)
}

// GetTextMsgDataEnter 用户进入消息
func GetTextMsgDataEnter(uuId, msgId, message string) string {

	return getTextMsgData("enter", uuId, msgId, message)
}

// GetTextMsgDataExit 用户退出消息
func GetTextMsgDataExit(uuId, msgId, message string) string {

	return getTextMsgData("exit", uuId, msgId, message)
}

// GetErrorMessage 根据错误码 获取错误信息
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		OK:                 "success",
		NotLoggedIn:        "未登录",
		ParameterIllegal:   "参数不合法",
		UnauthorizedUserId: "非法的用户Id",
		Unauthorized:       "未授权",
		NotData:            "没有数据",
		ServerError:        "系统错误",
		ModelAddError:      "添加错误",
		ModelDeleteError:   "删除错误",
		ModelStoreError:    "存储错误",
		OperationFailure:   "操作失败",
		RoutingNotExist:    "路由不存在",
	}

	if message == "" {
		if value, ok := codeMap[code]; ok {
			// 存在
			codeMessage = value
		} else {
			codeMessage = "未定义错误类型!"
		}
	} else {
		codeMessage = message
	}

	return codeMessage
}
