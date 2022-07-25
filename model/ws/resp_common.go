/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/16 0:14
 */

package ws

type JsonResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(code int, message string, data interface{}) JsonResult {

	message = GetErrorMessage(code, message)
	jsonMap := grantMap(code, message, data)

	return jsonMap
}

// 按照接口格式生成原数据数组
func grantMap(code int, message string, data interface{}) JsonResult {

	jsonMap := JsonResult{
		Code: code,
		Msg:  message,
		Data: data,
	}
	return jsonMap
}
