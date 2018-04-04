package main

import (
	"net/http"
	"putil/log"
	"dfqp/pg-http/service"
	"runtime"
	"framework/rpcclient/core"
	"dfqp/pg-http/controller"
)

type MyClientDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyClientDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(req)
	plog.Debug(body)

	responsestr := "hi I'm Server return your data"
	dis.client.SendPacket(req, []byte(responsestr))
}

func main() {
	service.Init()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//实例化rpc
	client := service.Client
	err := client.RpcInit("./conf/clientapp.conf")
	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}
	//服务基础配置
	mydisp := new(MyClientDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}

	//http服务
	//路由
	upload := new(controller.UploadApi)
	http.HandleFunc("/uploadAvatar", upload.UpAvatar)

	pay := new(controller.PayHttp)
	http.HandleFunc("/sendMoney", pay.SendMoney)

	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		plog.Fatal("start service fail:", err)
		return
	}
}
