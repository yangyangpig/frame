package main

import (
	//"sync"
	"fmt"
	"framework/rpcclient/core"
	"time"
	//	"framework/rpcclient/srpc"
	//	"framework/rpcclient"
	//"framework/rpcclient/bgf"
	//"framework/rpcclient/srpc"
	//"putil/byteorder"
	"putil/log"

	"../proto"
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
	svrid := int64(23)
	svrtype := int64(13)
	version := int32(1)
	groupids := []int64{1, 2, 3}

	client.SetSvr(svrtype, svrid, version, groupids)
	client.SetSvrNames("default", "default", []string{})
	client.SetLocalAddress("192.168.201.80", 55567)

	mydisp := new(MyDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp, client.Szkclient)
	client.WriteNormalLog("test", "here is normal log!")
	client.WriteRealLog("test", "here is real log!")
	client.WriteDebugLog("test", "here is debug log!")
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	for i := 0; i < 20; i++ {

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
		request := new(RPCProto.ArithRequest)
		request.A1 = 5
		request.A2 = 45
		req_bytes, err := request.Marshal()
		if err != nil {
			fmt.Println(req_bytes)
		}
		response := client.SendAndRecvRespRpcMsg("arith", "Add", req_bytes, 5000)
		if response.ReturnCode != 0 {
			//rpc返回结果异常
			plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		} else {
			arithResp := new(RPCProto.ArithResponse)
			arithResp.Unmarshal(response.Body)
			plog.Debug("return value  = ", arithResp.A3)
		}

		//测试点对点无返回值
		//		client.SendNoRespRpcMsg("arith", "Add", req_bytes)

		//测试点对组无响应
		//		client.SendGroupNotifyMsg("arith", 1, []byte("1"))
		time.Sleep(5 * time.Second)
	}

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
