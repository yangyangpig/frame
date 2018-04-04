package service

import "framework/rpcclient/core"

type pgGoodsDispatch struct {
	client *rpcclient.RpcCall
}

// RpcRequest 接收rpc请求并计算返回
func (dis *pgGoodsDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {



}