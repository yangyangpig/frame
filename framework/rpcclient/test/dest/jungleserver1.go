package main

import (
	//"fmt"
	_ "fmt"
	"framework/rpcclient/core"
	"os"
	"putil/log"
	"strconv"
	"time"

	"framework/rpcclient/test/dest/arith" //业务代码

	"framework/rpcclient/proto"
)

//待处理，
//网络出现问题 做通知
type MyDispatch struct {
	client *rpcclient.RpcCall
}

//所有注册的方法，先简单的做
var RegFuncs = []string{"Add", "Multiply"}
var arithobj = new(arith.Arith) //类似serveice

func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(rpcclient.GetRpcHeadString(req.Rpchead))
	plog.Debug("receive:", string(body))

	//业务接收数据并开始处理：
	methodName := string(req.Rpchead.MethodName) //方法名称
	rt := []byte{}                               //最终返回的字节切片
	//TODO检查方法名称是否存在
	switch methodName {
	case "Add":
		//接收参数并处理
		arithReq := new(RPCProto.ArithRequest)
		arithReq.Unmarshal(body)            //解参数的pb包
		arithResp := arithobj.Add(arithReq) //计算
		plog.Debug("return", arithResp.A3)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	case "Multiply":
		//接收参数并处理
		arithReq := new(RPCProto.ArithRequest)
		arithReq.Unmarshal(body)                      //解参数的pb包
		arithResp := arithobj.Multiply(arithReq)      //计算
		plog.Debug("Multiply returns:", arithResp.A3) //日志
		rt, _ = arithResp.Marshal()                   //返回数据的pb格式化
	default:
		//rt := []byte{}
	}

	dis.client.SendPacket(req, rt)
}

func main() {

	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}
	//默认值
	svrtype := int64(12)
	svrid := int64(27)
	version := int32(1)

	//cli模式下，获取args来设置svrid和svrtype以及version
	if len(os.Args) >= 4 {
		tsvrtype, _ := strconv.Atoi(os.Args[1])
		tsvrid, _ := strconv.Atoi(os.Args[2])
		tversion, _ := strconv.Atoi(os.Args[3])

		svrtype = int64(tsvrtype)
		svrid = int64(tsvrid)
		version = int32(tversion)
	}

	client.SetSvr(svrtype, svrid, version, []int64{1})
	client.SetSvrNames("arith", "arith", RegFuncs)
	client.SetLocalAddress("192.168.201.80", 55568)
	//client.SetLocalAddress("192.168.201.80", 55569)

	//gids := []int64{1, 2, 3}
	//	gids := []int64{1, 2, 5}
	//gids[0] = 20
	//client.SetGroupIds(gids)
	mydisp := new(MyDispatch)
	mydisp.client = client

	//	fmt.Println(client)
	//	return

	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	err = client.LaunchRpcClient("192.168.202.25", 7000, mydisp, client.Szkclient)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	arith.Logger = client.Szmqpushclient //给业务的日志实例赋值，rpc和业务共用一个！
	plog.Debug("LaunchRpcClient succeed!!!")
	//response := client.SendAndRecvRespRpcMsg("user.GetMoney", uid, 7, 18, []byte("helloword"), 10000)

	//plog.Debug(response.String())

	//	var controlstr string
	//	fmt.Scanln(&controlstr)
	//	if controlstr == "q" || controlstr == "Q" {
	//		return
	//	}
	time.Sleep(10000 * time.Second)

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
