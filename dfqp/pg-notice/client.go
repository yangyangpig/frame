package main

import (
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
	"fmt"
	"dfqp/proto/notice"
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	quitFlag := make(chan bool)
	//实例化rpc
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		return
	}

	err = client.RpcInit("./conf/clientapp.conf")
	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}
	//服务基础配置
	mydisp := new(MyClientDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)
	client.WriteNormalLog("test", "here is normal log!")
	client.WriteRealLog("test", "here is real log!")
	client.WriteDebugLog("test", "here is debug log!")
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	//超载包后面加一个正常包
	request := new(pgNotice.GetListRequest)
	request.Mid = 2000025
	req_bytes, err := request.Marshal()
	if err != nil {
		fmt.Println(req_bytes)
	}
	//		"arith"表示服务名
	//		Add表示调用的方法名
	//		req_bytes表示请求的参数（经过了protobuf的marshal后）
	//		5000表示5000毫秒后无响应就超时！
	response := client.SendAndRecvRespRpcMsg("pgNotice.NoticeList", req_bytes, 5000, 0)
	plog.Debug("rpc response = ",response)
	if response.ReturnCode != 0 {
		//rpc返回结果异常
		//RPC_RESPONSE_COMPLET         = 0 //完成
		//	RPC_RESPONSE_TIMEOUT         = 1 //超时
		//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
		//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {
		arithResp := new(pgNotice.GetListResponse)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return value  = ", arithResp)
	}

	<-quitFlag
}