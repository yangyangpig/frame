
package main

import (
	//"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"fmt"
)


//待处理，
//网络出现问题 做通知

type MyDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(rpcclient.GetRpcHeadString(req.Rpchead))
	plog.Debug("receive:", string(body))

	responsestr := "hi I'm Server return your data"
	plog.Debug("return", responsestr)
	dis.client.SendPacket(req, []byte(responsestr))
}

func main() {

	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}
	svrid := int32(22)
	svrtype := int32(12)
	//	uid := int32(100)
	client.SetSvr(svrtype, svrid)
	client.SetLocalAddress("192.168.201.80", 55568)
//	gids := make([]int64, 1)
//	gids[0] = 20
//	client.SetGroupIds(gids)
	mydisp := new(MyDispatch)
	mydisp.client = client
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	//response := client.SendAndRecvRespRpcMsg("function1", uid, 7, 18, []byte("helloword"), 10000)

	//plog.Debug(response.String())

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
