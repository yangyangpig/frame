package main

import (
	//"sync"
//	"fmt"
	"framework/rpcclient/core"
	//	"framework/rpcclient/srpc"
	//	"framework/rpcclient"
	"putil/log"
	"framework/rpcclient/bgf"
	"putil/byteorder"
	"framework/rpcclient/srpc"
	
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
/*
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}

	svrid := int32(23)
	svrtype := int32(13)
	
	desttype := int32(12)
	destid := int32(22)

	client.SetSvr(svrtype, svrid)
	client.SetLocalAddress("172.20.153.46", 55567)
*/

	var data77 []byte =                []byte{0x0,0x0,0x0,0x4d,0xa,0x15,0x8,0x1,0x10,0x89,0x2,0x18,0x0,0x20,0x0,0x28,0x17,0x30,0xd,0x38,0x16,0x40,0xc,0x58,0x0,0x60,0x7b,0x12,0x30,0x28,0x0,0x0,0x0,0x1d,0x0,0x0,0x0,0x9,0x50,0x7b,0xa0,0x1,0x0,0xa8,0x1,0x0,0xf0,0x1,0x7b,0xc0,0x2,0x0,0xc8,0x2,0x0,0x9a,0x3,0x9,0x66,0x75,0x6e,0x63,0x74,0x69,0x6f,0x6e,0x31,0x68,0x65,0x6c,0x6c,0x6f,0x77,0x6f,0x72,0x64,0x29}
	//[]byte{0x0,0x0,0x0,0x4d,0xa,0x15,0x8,0x1,0x10,0x89,0x2,0x18,0x0,0x20,0x0,0x28,0x17,0x30,0xd,0x38,0x16,0x40,0xc,0x58,0x0,0x60,0x3,0x12,0x30,0x28,0x0,0x0,0x0,0x1d,0x0,0x0,0x0,0x9,0x50,0x3,0xa0,0x1,0x0,0xa8,0x1,0x0,0xf0,0x1,0x3,0xc0,0x2,0x0,0xc8,0x2,0x0,0x9a,0x3,0x9,0x66,0x75,0x6e,0x63,0x74,0x69,0x6f,0x6e,0x31,0x68,0x65,0x6c,0x6c,0x6f,0x77,0x6f,0x72,0x64,0x29}
	var data80 []byte = []byte{0x0,0x0,0x0,0x50,0xa,0x16,0x8,0x1,0x10,0x89,0x2,0x18,0x0,0x20,0x0,0x28,0x17,0x30,0xd,0x38,0x16,0x40,0xc,0x58,0x0,0x60,0x98,0x2,0x12,0x32,0x28,0x0,0x0,0x0,0x1f,0x0,0x0,0x0,0x9,0x50,0x98,0x2,0xa0,0x1,0x0,0xa8,0x1,0x0,0xf0,0x1,0x98,0x2,0xc0,0x2,0x0,0xc8,0x2,0x0,0x9a,0x3,0x9,0x66,0x75,0x6e,0x63,0x74,0x69,0x6f,0x6e,0x31,0x68,0x65,0x6c,0x6c,0x6f,0x77,0x6f,0x72,0x64,0x29}





	datalen1, err := byteOrderOp.Util_ntoh_int32(data77)
	if err != nil {
		plog.Fatal("decode datapackage len err", err)
	}
	plog.Debug("data77 datalen:=", datalen1)
	var bgfmsg RPCProto.BGFMsg
	decodedata1 := data77[4:datalen1]
	err = bgfmsg.Unmarshal(decodedata1)

	rpchead1, body1, err1 := srpc.SrpcUnpackPkgHeadandBody(bgfmsg.Body)
	if err1 != nil {
		plog.Debug("SrpcUnpackPkgHeadandBody")
	}

	plog.Debug(rpcclient.GetRpcHeadString(rpchead1))
	plog.Debug(string(body1))


	datalen, err := byteOrderOp.Util_ntoh_int32(data80)
	if err != nil {
		plog.Fatal("decode datapackage len err", err)
	}
	plog.Debug("data80 datalen:=", datalen)

	decodedata := data80[4:datalen]
	err = bgfmsg.Unmarshal(decodedata)
	

	rpchead, body, err1 := srpc.SrpcUnpackPkgHeadandBody(bgfmsg.Body)
	if err1 != nil {
		plog.Debug("SrpcUnpackPkgHeadandBody")
	}

	plog.Debug(rpcclient.GetRpcHeadString(rpchead))
	plog.Debug(string(body))



	for i := 0 ; i < 10 ; i++ {
		var head RPCProto.BGFHead
		head.Uid = 1
		head.Cmd = srpc.RPC_REQUEST_PACKAGE
		head.DSerId = 12
		head.DSerType = 22
		head.SSerId = 23
		head.SSerType = 13
		head.TransTypes = srpc.TRANS_P2P
		head.MsgType = 0
		head.MtId = 0

		//rpchead := new(srpc.CRpcHead)
		var rpchead srpc.CRpcHead
		rpchead.MethodName = "function1"
		rpchead.Sequence = uint64(i)
		rpchead.FlowId = uint64(i)

		head.MtId = int32(i)



		rpcpack, err := srpc.SrpcPackPkg(&rpchead, []byte("helloword"))
		if err != nil {
			plog.Fatal(err)
			return
		}


		bfgmsg := new(RPCProto.BGFMsg)
		bfgmsg.BGFHead = head
		bfgmsg.Body = rpcpack

		lenbuf := make([]byte, 4, 1600)
		//var lenbuf [4] byte

		datalen := bfgmsg.Size()

		byteOrderOp.Util_hton_int32(int32(datalen+4), lenbuf)
		lenbuf1 := lenbuf[4 : datalen+4]
		marshalint, err := bfgmsg.MarshalTo(lenbuf1)
		if err != nil {
			plog.Fatal("Marshall err :", err)
		}
		lenbuf = lenbuf[:datalen+4]
		plog.Debug("i =", i, "send msg size :", len(lenbuf), "data:", "marshalint = ", marshalint)


	}
















	//mydisp := new(MyDispatch)
	//mydisp.client = client
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	//if err != nil {
	//	plog.Fatal("lauch failed", err)
	//	return
	//}
	//plog.Debug("LaunchRpcClient succeed!!!")
	//response := client.SendAndRecvRespRpcMsg("function1", desttype, destid, []byte("helloword"), 5000)
	//plog.Debug("return code = ", response.ReturnCode)
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
	
//	var controlstr string
//	fmt.Scanln(&controlstr)

//	if controlstr == "q" || controlstr == "Q" {
//		return
//	}

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
