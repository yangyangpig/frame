package main

import (
	//"sync"
	"fmt"
	"framework/rpcclient/core"
	//	"framework/rpcclient/srpc"
	//	"framework/rpcclient"
	//"framework/rpcclient/bgf"
	//"framework/rpcclient/srpc"
	//"putil/byteorder"
	"putil/log"

	"PGNotices/app/proto"
)

type MyDispatch1 struct {
	client *rpcclient.RpcCall
}

func (dis *MyDispatch1) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
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

	svrid := int64(24)
	svrtype := int64(14)

	client.SetSvr(svrtype, svrid, 1, nil)
	//client.SetLocalAddress("192.168.56.102", 55570)

	mydisp := new(MyDispatch1)
	mydisp.client = client
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp, client.Szkclient)
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	request := new(Notices.GetListRequest)
	request.App = 101000
	req_bytes, _ := request.Marshal()
	response := client.SendAndRecvRespRpcMsg("noticesserver.GetList", req_bytes, 5000)
	plog.Debug("return code = ", response.ReturnCode)
	arithResp := new(Notices.GetListResponse)
	arithResp.Unmarshal(response.Body)
	plog.Debug("return value  = ", arithResp)
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
