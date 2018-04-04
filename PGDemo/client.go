package main

import (
	//"sync"
	"PGDemo/app/arithproto"
	"fmt"
	"framework/rpcclient/core"
	"runtime"
	"time"
	//	"framework/rpcclient/srpc"
	//	"framework/rpcclient"
	//"framework/rpcclient/bgf"
	//"framework/rpcclient/srpc"
	//"putil/byteorder"
	"putil/log"
	//"framework/rpcclient/proto"
)

var (
	quitFlag chan bool //退出标记
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
	plog.Info("开始")
	plog.Debug("开始")
	plog.Warn("开始")
	plog.Fatal("开始")
	plog.Core("开始")
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer plog.CatchPanic()
	quitFlag = make(chan bool)
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

	mydisp := new(MyDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)
	plog.Debug("开始")
	//client.WriteNormalLog("try_log", "here is normal log!")
	//client.WriteRealLog("try_log", "here is real log!")
	//client.WriteDebugLog("try_log", "here is debug log!")
	plog.Debug("结束")
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	time.Sleep(2 * time.Second)
	//for i := 0; i < 1000; i++ {

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

	/***************************点对点有返回--并发包个数*****************************/
	var concurrent_times = 1 //并发个数 1000
	var reqs [][]byte = make([][]byte, 0, concurrent_times)
	var result []int32 = make([]int32, 0, concurrent_times)
	//超载包后面加一个正常包
	request := new(arithproto.ArithRequest)
	request.A1 = 1
	for i := 0; i < concurrent_times; i++ {
		request.A2 = uint32(i)
		req_bytes, err := request.Marshal() //把本进程的数据结构序列化
		if err != nil {
			plog.Debug(err)
		}
		plog.Debug(req_bytes)
		reqs = append(reqs, req_bytes)
		plog.Debug(reqs)
	}

	for i := 0; i < 200000; i++ {
		//		"arith"表示服务名
		//		Add表示调用的方法名
		//		req_bytes表示请求的参数（经过了protobuf的marshal后）
		//		5000表示5000毫秒后无响应就超时！
		for j := 0; j < concurrent_times; j++ {
			go func(j int) {
				plog.Debug("the len is:", len(reqs), "the j is:", j)
				response := client.SendAndRecvRespRpcMsg("arith.Add", reqs[j], 5000, 0)
				if response.ReturnCode != 0 {
					//rpc返回结果异常
					//RPC_RESPONSE_COMPLET         = 0 //完成
					//	RPC_RESPONSE_TIMEOUT         = 1 //超时
					//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
					//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
					//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
					plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
				} else {
					arithResp := new(arithproto.ArithResponse)
					arithResp.Unmarshal(response.Body)
					plog.Debug("return value  = ", arithResp.A3)
					result = append(result, arithResp.A3)
				}
			}(j)
		}

		time.Sleep(1000 * time.Millisecond)
		//time.Sleep(8 * time.Second)
		plog.Debug("the concurrent test result is:", result)
		plog.Debug("the concurrent test result count:", len(result))
		//time.Sleep(2 * time.Second)

	}

	/***************************测试点对点无返回值*****************************/
	//	request := new(arithproto.ArithRequest)
	//	request.A1 = 1
	//	request.A2 = 5
	//	req_bytes, _ := request.Marshal()
	//	res := client.SendNoRespRpcMsg("arith.Add", req_bytes, 5000, 0)
	//	plog.Debug(res)

	/***************************测试点对组无返回值*****************************/
	//测试点对组无响应
	//		client.SendGroupNotifyMsg("arith", 1, []byte("1"))
	//time.Sleep(5 * time.Second)
	//}

	//response = client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//plog.Debug("return code = ", response.ReturnCode)
	//response = client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//plog.Debug("return code = ", response.ReturnCode)
	//response = client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//plog.Debug("return code = ", response.ReturnCode)

	//response := client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//if response.ReturnCode != rpcclient.RPC_RESPONSE_COMPLET {
	//	plog.Debug("return ok")
	//}
	//plog.Debug("return:", response.String())

	//===============================================

	//var sw sync.WaitGroup
	//sw.Add(4)
	//for i := 1; i < 5; i++ {
	//	go func(i int) {
	//		plog.Debug("i = ", i)
	//		response := client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//		plog.Debug("i=", i, "return code = ", response.ReturnCode)
	//		//sw.Done()
	//	}(i)
	//}
	//sw.Wait()

	/*response = client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	if response.ReturnCode != rpcclient.RPC_RESPONSE_COMPLET {
		plog.Debug("return ok")
	}
	plog.Debug("return:", response.String())

	response = client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	if response.ReturnCode != rpcclient.RPC_RESPONSE_COMPLET {
		plog.Debug("return ok")
	}
	plog.Debug("return:", response.String())*/

	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		return
	}

	/*
		s1 := new(srpc.CRpcHead)
		s1.FlowId = 1
		s1.Sequence = 1
		s1.MethodName = "function1"

		s2 := []byte("helloworld")
		p0, err := srpc.SrpcPackPkg(s1, s2)
		s3, s4, err := srpc.SrpcUnpackPkgHeadandBody(p0)
		fmt.Println(s3, string(s4), err)
	*/
}
