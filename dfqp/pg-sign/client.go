package main

import (
	"dfqp/proto/sign"
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
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
	// 实例化rpc
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

	mydisp := new(MyClientDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)
	client.WriteNormalLog("test", "here is normal log!")
	client.WriteRealLog("test", "here is real log!")
	client.WriteDebugLog("test", "here is debug log!")
	if err != nil {
		plog.Fatal("lauch failed", err)
	}
	plog.Debug("LaunchRpcClient succeed!!!")

	// test signin接口
	request := new(pgSign.SigninRequest)
	request.Mid = 23232323
	request.Day = 1
	req_bytes, err := request.Marshal()
	plog.Debug("req_bytes", req_bytes)
	if err != nil {
		fmt.Println(req_bytes)
	}

	response := client.SendAndRecvRespRpcMsg("pgSign.signin", req_bytes, 5000, 0)
	plog.Debug("client", client)
	plog.Debug("response", response)
	plog.Debug("req_bytes", req_bytes)
	if response.ReturnCode != 0 {
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {
		arithResp := new(pgSign.SigninResponse)
		plog.Debug("response.Body", response.Body)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return value  = ", arithResp)
	}

	<-quitFlag
}
