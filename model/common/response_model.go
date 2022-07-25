/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 0:08
 */

package common

import "encoding/json"

/************************  响应数据  **************************/

type Head struct {
	Seq      string    `json:"seq"`      // 消息的Id
	Event    string    `json:"event"`    // 消息的event
	Response *Response `json:"response"` // 消息体
}

type Response struct {
	Code    int         `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

// PushMsg push 数据结构体
type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// NewResponseHead 设置返回消息
func NewResponseHead(seq string, event string, code int, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Event: event, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)

	return
}

func NewResponse(code int, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}
