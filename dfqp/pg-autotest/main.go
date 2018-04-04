package main

import (
	"framework/rpcclient/core"
	"putil/log"
	"dfqp/pg-autotest/service"
	"runtime"
)

var (
	quitFlag chan bool	//退出标记
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
	err := client.RpcInit("./conf/app.conf")
	if err != nil {
		plog.Debug("RPC 初始化异常")
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
	//控制退出
	_ = <-quitFlag
}

