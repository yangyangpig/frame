package main

import (
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"time"

	"PGDemo/proto/yihua"
)

type MyDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(req)
	plog.Debug(body)

	responsestr := "hi I'm Server return your data"
	dis.client.SendPacket(req, []byte(responsestr))
}

func main() {
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}

	//desttype := int64(12)
	//destid := int64(22)
	//服务基础配置
	svrid := int64(1538)
	svrtype := int64(13)
	version := int32(1)
	groupids := []int64{1, 2, 3}

	client.SetSvr(svrtype, svrid, version, groupids)
	client.SetSvrNames("default", "default", []string{})
	//client.SetLocalAddress("192.168.201.80", 55567)

	mydisp := new(MyDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp, client.Szkclient)

	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	for i := 0; i < 3; i++ {

		//超载的包
		//		req_bytes := make([]byte, 12300)
		//		//测试点对点又返回值
		//		response := client.SendAndRecvRespRpcMsg("arith", "Add", req_bytes, 5000)
		//		if response.ReturnCode != 0 {
		//			//rpc返回结果异常
		//			plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		//		} else {
		//			arithResp := new(RPCProto.ArithResponse)
		//			arithResp.Unmarshal(response.Body)
		//			plog.Debug("return value  = ", arithResp.A3)
		//		}

		//超载包后面加一个正常包
		request := new(data.GetMoneyRequest)
		request.Uid = 9001386
		req_bytes, err := request.Marshal()
		if err != nil {
			fmt.Println(req_bytes)
		}
		//		"arith"表示服务名
		//		Add表示调用的方法名
		//		req_bytes表示请求的参数（经过了protobuf的marshal后）
		//		5000表示5000毫秒后无响应就超时！
		response := client.SendAndRecvRespRpcMsg("data.DataService.GetMoney", req_bytes, 5000)
		if response.ReturnCode != 0 {
			//rpc返回结果异常
			//RPC_RESPONSE_COMPLET         = 0 //完成
			//	RPC_RESPONSE_TIMEOUT         = 1 //超时
			//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
			//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
			//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
			plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		} else {
			arithResp := new(data.GetMoneyResponse)
			arithResp.Unmarshal(response.Body)
			plog.Debug("return value  = ", arithResp.Money)
		}

		time.Sleep(5 * time.Second)
	}

	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		return
	}
}
