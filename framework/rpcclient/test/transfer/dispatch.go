package main

import (
	"sync"
	"net"
	"putil/log"
	"fmt"
	"time"
	"framework/rpcclient/bgf"
	"framework/rpcclient/net"
	"framework/rpcclient/srpc"

	"putil/byteorder"
	"framework/rpcclient/core"
)
var(
	isexit bool
	clmap map[string]*rpcnet.TcpNet
	clMutex sync.Mutex
)

func init(){
	isexit = false
	clmap = make(map[string]*rpcnet.TcpNet, 10)
}

func addcl(id int , tcpnet *rpcnet.TcpNet){	
	clMutex.Lock()
	defer clMutex.Unlock()
	plog.Debug("addcl:", id)
	clmap[id] = tcpnet	
}

func delcl(id int){
	clMutex.Lock()
	defer clMutex.Unlock()
	plog.Debug("delcl:", id)
	delete(clmap, id)
}


func clientinput(tcpnet *rpcnet.TcpNet){
	protolen := make([]byte, 4)
	data := make([]byte, 1600)
	for !isexit {
		if tcpnet.ReadAtLeast(protolen, 4); err != nil {
			plog.Fatal("io exception  ", err)
			tcpnet.Close()
			return
		} 
		
		var datalen int32
		datalen, err := byteOrderOp.Util_ntoh_int32(protolen)
		if err != nil {
			plog.Fatal("decode datapackage len err", err)
		}
		plog.Debug("recv datalen", datalen)

		if err := rpcclient.net.ReadAtLeast(data, int(datalen-4)); err != nil {
			plog.Fatal("io excepton", err)
			rpcclient.close()
			return
		}

		bgfmsg := new(RPCProto.BGFMsg)
		decodedata := data[:datalen-4]
		err = bgfmsg.Unmarshal(decodedata)
		if err != nil {
			plog.Fatal("occure decode err!")
			continue
		}
		
		processPacket(bgfmsg)
	}
}


func processPacket(bgfmsg *RPCProto.BGFMsg) {
	plog.Debug("processPacket", bgfmsg.BGFHead.Cmd)
	switch bgfmsg.BGFHead.Cmd {
	case srpc.DISPATCH_REG_SER:
		fallthrough
	case srpc.DISPATCH_REG_RSP:
		fallthrough
	case srpc.DISPATCH_SENDINFO_SER:
		processRegisterResp(bgfmsg)
	case srpc.RESPONSE_DISPATCH_KEEP_ALIVE_SER:		
	case srpc.RPC_REQUEST_PACKAGE:
		processRpcRequest(bgfmsg)
	case srpc.RPC_RESPONSE_PACKET:
		processRpcResponse(bgfmsg)
	}
	//  return 1
}


func processRpcRequest(bgfmsg *RPCProto.BGFMsg){
	
}

func processRpcResponse(bgfmsg *RPCProto.BGFMsg){
	
}

func main(){


	
	server, err := net.Listen("tcp", ":9081")
	if err != nil{
		plog.Debug("net.Listen exception", err)
	}

	go func(){
		for !isexit {
			con, err := server.Accept()
			if err != nil {
				plog.Debug("server.Accept() ", err)
			}
			tcpnet := new(rpcnet.TcpNet)
			tcpcon ,is := con.(*net.TCPConn)
			if !is{
				plog.Debug("con is note a*net.TCPConn type")
				continue
			}
			tcpnet.CofigConn(tcpcon)
			go clientinput(tcpnet)


		}
	}()

	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		isexit = true
		time.Sleep(time.Second)
		return
	}
}
