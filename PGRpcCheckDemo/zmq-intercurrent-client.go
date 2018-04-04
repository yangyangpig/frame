package main

import (
	_ "PGRpcCheckDemo/app/bigPackData"
	"PGRpcCheckDemo/app/controllers"
	"PGRpcCheckDemo/app/proto"
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"time"
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

var (
	count int = 10000
	chs       = make([]chan *rpcclient.RpcResponse, count)
)

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

	//向zmq推送三种日志
	//	client.WriteNormalLog("test", "here is normal log!")
	//	client.WriteRealLog("test", "here is real log!")
	//	client.WriteDebugLog("test", "here is debug log!")

	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	request := new(RPCProto.BigPackDataRequest)
	request.ContainerReq = createBigData()
	req_bytes, err := request.Marshal()
	plog.Debug("req_bytes length  = ", len(req_bytes))

	if err != nil {
		fmt.Println(req_bytes)
	}
	//		"arith"表示服务名
	//		Add表示调用的方法名
	//		req_bytes表示请求的参数（经过了protobuf的marshal后）
	//		5000表示5000毫秒后无响应就超时！
	for i := 0; i < count; i++ {
		plog.Debug("并发请求", i)
		chs[i] = make(chan *rpcclient.RpcResponse)
		go concurrentRequest(client, i, req_bytes)
		time.Sleep(2 * time.Second)

	}

	for key, _ := range chs {
		plog.Debug("执行变量channel次数为", key)
		res := <-chs[key]
		plog.Debug("执行变量channel变量为", res.Seq)
		if res.ReturnCode != 0 {
			//rpc返回结果异常
			//RPC_RESPONSE_COMPLET         = 0 //完成
			//	RPC_RESPONSE_TIMEOUT         = 1 //超时
			//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
			//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
			//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
			plog.Debug("rpc return code = ", res.ReturnCode, " return err = ", res.Err)
		} else {
			plog.Debug("返回开始输出", key)
			arithResp := new(RPCProto.BigPackDataResponse)
			arithResp.Unmarshal(res.Body)
			plog.Debug("并发请求回复序列  = ", res.Seq)
			plog.Debug("return value  = ", arithResp)
		}
	}

	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		return
	}
}

func createBigData() map[string]string {
	bigPack := new(controllers.BigPack)
	bigPack.Route = 7
	bigPack.Container = make(map[string]string)
	container := bigPack.CreateBigPack()
	return container
}

//并发请求
func concurrentRequest(client *rpcclient.RpcCall, j int, req []byte) {
	plog.Debug("并发请求开始", j)
	response := client.SendAndRecvRespRpcMsg("bigdata.GetData", req, 5000, 0)
	chs[j] <- response

	//向zmq推送三种日志

	client.WriteNormalLog("test", "here is normal log!")
	//	client.WriteRealLog("test", "here is real log!")
	//	client.WriteDebugLog("test", "here is debug log!")
	plog.Debug("并发请求结束", j)

}
