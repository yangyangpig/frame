package arith

import (
	"framework/rpcclient/proto"
	"framework/rpcclient/szmq" //zmq日志操作
)

type Arith int

var (
	Logger *szmq.SzmqPushClient //通过zmq写日志！
)

//方法需要手动注册

//加法服务
func (arith *Arith) Add(rq *RPCProto.ArithRequest) *RPCProto.ArithResponse {
	rp := new(RPCProto.ArithResponse)
	rp.A3 = int32(rq.A1 + rq.A2)
	Logger.WriteNormalLog("test", "hello,the log comes from business package!")
	return rp
}

//乘法服务
func (arith *Arith) Multiply(rq *RPCProto.ArithRequest) *RPCProto.ArithResponse {
	rp := new(RPCProto.ArithResponse)
	rp.A3 = int32(rq.A1 * rq.A2)
	return rp
}
