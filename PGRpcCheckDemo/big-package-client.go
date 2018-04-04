package main

import (
	_ "PGRpcCheckDemo/app/bigPackData"
	"PGRpcCheckDemo/app/controllers"
	"PGRpcCheckDemo/app/proto"
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
)

type MyBigPackData struct {
	client *rpcclient.RpcCall
}

func (dis *MyBigPackData) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	//plog.Debug(req)
	//plog.Debug(body)

	responsestr := "hi I'm Server return your data"
	dis.client.SendPacket(req, []byte(responsestr))
}

func main() {
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}

	err = client.RpcInit("./conf/clientapp.conf")
	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}

	mydisp := new(MyBigPackData)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)

	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	request := new(RPCProto.BigPackDataRequest)
	request.ContainerReq = createBigData()
	req_bytes, err := request.Marshal()

	if err != nil {
		fmt.Println(req_bytes)
	}
	//		"arith"表示服务名
	//		Add表示调用的方法名
	//		req_bytes表示请求的参数（经过了protobuf的marshal后）
	//		5000表示5000毫秒后无响应就超时！
	for i := 0; i < 1; i++ {
		response := client.SendAndRecvRespRpcMsg("bigdata.GetData", req_bytes, 5000, 0)
		if response.ReturnCode != 0 {
			//rpc返回结果异常
			//RPC_RESPONSE_COMPLET         = 0 //完成
			//	RPC_RESPONSE_TIMEOUT         = 1 //超时
			//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
			//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
			//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
			plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		} else {
			arithResp := new(RPCProto.BigPackDataResponse)
			arithResp.Unmarshal(response.Body)
			plog.Debug("return value  = ", arithResp)
		}
	}
}

func createBigData() map[string]string {
	bigPack := new(controllers.BigPack)
	bigPack.Route = 500
	bigPack.Container = make(map[string]string)
	container := bigPack.CreateBigPack()
	return container
}
