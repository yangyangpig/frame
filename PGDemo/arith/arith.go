package arith

import (
	"PGDemo/app/arithproto"
	"framework/rpcclient/szmq" //zmq日志操作
)

type Arith int

//var (
//	Logger *szmq.SzmqPushClient //通过zmq写日志！
//)

//方法需要手动注册

//加法服务
func (arith *Arith) Add(rq *arithproto.ArithRequest) *arithproto.ArithResponse {
	rp := new(arithproto.ArithResponse)
	rp.A3 = int32(rq.A1 + rq.A2)

	go szmq.Logger.WriteNormalLog("test", "hello,the log comes from arith!")

	return rp
}

//乘法服务
func (arith *Arith) Multiply(rq *arithproto.ArithRequest) *arithproto.ArithResponse {
	rp := new(arithproto.ArithResponse)
	rp.A3 = int32(rq.A1 * rq.A2)
	return rp
}
