package main

import (
	"PGReport/app/controllers"
	//	"PGNotices/app/proto"
	//"PGReport/app/service"
	//"fmt"
	//"framework/rpcclient/core"
	// "framework/rpcclient/tcleaest/dest/arith"
	//	"putil/log"
	//	"github.com/astaxie/beego"
	"framework/rpcclient/core"
	"putil/log"
	"PGReport/app/proto"
	"fmt"
	//"PGNotices/app/service"
	"strings"
	"PGLibrary"
)

type MyDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(rpcclient.GetRpcHeadString(req.Rpchead))
	plog.Debug("receive:", string(body))

	serviceAndMethodName := string(req.Rpchead.MethodName) //方法名称(eg：User.getUserInfo)
	var methodName string
	//转成纯方法名称
	pos := strings.LastIndex(serviceAndMethodName, ".")
	if pos > 0 && pos < len(serviceAndMethodName) {
		methodName = serviceAndMethodName[pos+1:] //避免越界
	}
	//业务接收数据并开始处理：
	//methodName := string(req.Rpchead.MethodName) //方法名称
	rt := []byte{}                               //最终返回的字节切片
	//TODO检查方法名称是否存在
	plog.Debug("method:")

	// 公共检测特殊包
	var endRt, endReturn []byte
	//comReq := new(Report.CommonRequest)
	//comReq.Unmarshal(body)
	//switch comReq.XSpecialFlag_ {
	//case 1: // unzip
	//	endBody = PGLibrary.DoUnGzip(comReq.XSpecialString_)
	//case 2:
	//	endBody = PGLibrary.Base64Decode(comReq.XSpecialString_)
	//case 3:
	//	tmpBody := PGLibrary.Base64Decode(comReq.XSpecialString_)
	//	endBody = PGLibrary.DoUnGzip(string(tmpBody))
	//default:
	//}
	//if endBody == nil {
	//	endBody = body
	//}
	//plog.Debug("decode:=========")
	//plog.Debug(endBody)
	switch methodName {
	case "Devices":
		//接收参数并处理
		plog.Debug("method:Devices")
		req := new(Report.DevicesRequest)
		req.Unmarshal(body)                             //解参数的pb包
		reportObj := new(controllers.ReportController) //计算
		resp := reportObj.Devices(req)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	default:

	}
	if  rt != nil { //len(rt) > 7*1024
		comReq := new(Report.CommonRequest)
		endRt = PGLibrary.DoGzip(rt)
		//endRt = PGLibrary.Base64Encode(rt)

		comReq.XSpecialString_ = string(endRt)
		comReq.XSpecialFlag_ = 1
		endReturn, _ = comReq.Marshal()
	} else {
		endReturn = rt
	}

	plog.Debug("byte len:=========")
	plog.Debug(len(endRt))
	plog.Debug(endRt)
	plog.Debug("pb len:=========")
	plog.Debug(len(endReturn))
	dis.client.SendPacket(req, endReturn)
}

var RegFuncs = []string{"Devices"}

func main() {
	//service.Init()
	//web http测试
	//beego.Router("/", &controllers.NoticesController{})
	//beego.Run()

	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}
	svrtype := int64(1001)
	svrid := int64(10001)
	client.SetSvr(svrtype, svrid, 1, nil)
	client.SetSvrNames("Report", "devices", RegFuncs)
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
