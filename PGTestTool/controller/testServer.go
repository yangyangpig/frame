package controller

import (
	"dfqp/lang"
	"dfqp/pg-http/service"
	"encoding/base64"
	"encoding/json"
	_ "fmt"
	_ "net"
	"net/http"
	"putil/log"
)

type TestServerController struct{}

const (
	Success = iota
	Faild
)

func (this *TestServerController) RequestData(writer http.ResponseWriter, request *http.Request) {
	//获取请求数据
	request.ParseForm()
	servername := request.PostFormValue("servername")
	funcname := request.PostFormValue("funcname")

	reqBytesTmp := request.PostFormValue("params")
	reqBytes, _ := base64.StdEncoding.DecodeString(reqBytesTmp)

	// 初始化
	code := Success

	requestFunc := string(servername) + "." + string(funcname)
	response := service.Client.SendAndRecvRespRpcMsg(requestFunc, reqBytes, 5000, 0)
	if response.ReturnCode != 0 {
		code = Faild
	}
	//response.Body是[]byte类型
	res := base64.StdEncoding.EncodeToString(response.Body)
	this.output(writer, code, res, "")
	return
}

func (this *TestServerController) ResponseData(writer http.ResponseWriter, request *http.Request) {

}

func (this *TestServerController) output(writer http.ResponseWriter, code int, data interface{}, msg string) {
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
