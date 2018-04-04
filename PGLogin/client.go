package main

import (
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"Proto/login"
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
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}

	//desttype := int64(12)
	//destid := int64(22)
	//服务基础配置
	svrid := int64(1538)
	svrtype := int64(1981)
	version := int32(1)
	groupids := []int64{1, 2, 3}

	client.SetSvr(svrtype, svrid, version, groupids)
	client.SetSvrNames("default", "default", []string{})
	//client.SetLocalAddress("192.168.201.80", 55567)

	mydisp := new(MyClientDispatch)
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
	//		"arith"表示服务名
	//		Add表示调用的方法名
	//		req_bytes表示请求的参数（经过了protobuf的marshal后）
	//		5000表示5000毫秒后无响应就超时！
	for i:=0; i<1; i++ {
		go func() {
			//超载包后面加一个正常包
			request := new(pgLogin.LoginRequest)
			request.Guid = "0017A793-D6BC-4843-A3B8-9D3FC9190BB5"
			request.LoginType = 1
			request.AppId = 100100
			req_bytes, err := request.Marshal()
			if err != nil {
				fmt.Println(req_bytes)
			}

			response := client.SendAndRecvRespRpcMsg("PGLogin.Login", req_bytes, 5000, 0)
			if response.ReturnCode != 0 {
				//rpc返回结果异常
				//RPC_RESPONSE_COMPLET         = 0 //完成
				//	RPC_RESPONSE_TIMEOUT         = 1 //超时
				//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
				//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
				//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
				plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
			} else {
				arithResp := new(pgLogin.LoginResponse)
				arithResp.Unmarshal(response.Body)
				plog.Debug("return value  = ", arithResp)
			}
		}()
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
