package controller

import (
	"net/http"
	"dfqp/lang"
	"encoding/json"
)

type base struct {}
//生成响应结果
//code api接口状态码
//err 错误信息
//data 处理结果
func (bs *base) output(writer http.ResponseWriter, code int, data interface{}, msg string) {
	resp := make(map[string]interface{})
	resp["code"] = code
	if msg == "" {
		//取语言包
		msg = lang.Msg(code, lang.ZH)
	}
	resp["msg"] = msg
	resp["data"] = data
	jsonStr, _ := json.Marshal(resp)
	writer.Header().Set("Content-Type", "application/json;charset=utf-8")
	writer.Write(jsonStr)
	return
}
