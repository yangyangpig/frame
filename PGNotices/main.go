package main

import (
	"PGNotices/app/controllers"
//	"PGNotices/app/proto"
	"PGNotices/app/service"
	//"fmt"
	//"framework/rpcclient/core"
	// "framework/rpcclient/tcleaest/dest/arith"
//	"putil/log"
//	"github.com/astaxie/beego"
	"framework/rpcclient/core"
	"putil/log"
	"PGNotices/app/proto"
	"fmt"
)

type MyDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(rpcclient.GetRpcHeadString(req.Rpchead))
	plog.Debug("receive:", string(body))

	//业务接收数据并开始处理：
	methodName := string(req.Rpchead.MethodName) //方法名称
	rt := []byte{}                               //最终返回的字节切片
	//TODO检查方法名称是否存在
	plog.Debug("method:")
	switch methodName {
	case "GetList":
		//接收参数并处理
		plog.Debug("method:getlist")
		req := new(Notices.GetListRequest)
		req.Unmarshal(body)                             //解参数的pb包
		noticeObj := new(controllers.NoticesController) //计算
		resp := noticeObj.GetAll(req)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	default:

	}

	dis.client.SendPacket(req, rt)
}

var RegFuncs = []string{"GetList"}

func main() {
	service.Init()
	//web http测试
	//beego.Router("/", &controllers.NoticesController{})
	//beego.Run()

	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}
	svrid := int64(25)
	svrtype := int64(15)
	client.SetSvr(svrtype, svrid, 1, nil)
	client.SetSvrNames("noticesserver", "noticesserver", RegFuncs)
	//client.SetLocalAddress("192.168.56.102", 55569)

	//service.Logger = client.Szmqpushclient

	mydisp := new(MyDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp, client.Szkclient)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		return
	}
}
