package main

import (
	"PGDemo/app/arithproto"
	"PGDemo/arith" //业务代码
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
	//"strconv"
	"strings"
	//"github.com/astaxie/beego/config"
)

var (
	quitFlag chan bool          //退出标记
	arithobj = new(arith.Arith) //Server实例，为处理接收到的请求做准备
)

//包含了rpc实例的结构体，实现了rpccall里面的接口方法RpcRequest
type MyDispatch struct {
	client *rpcclient.RpcCall
}

//接收rpc请求并计算返回
func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	serviceAndMethodName := string(req.Rpchead.MethodName) //方法名称(eg：User.getUserInfo)
	var methodName string
	//转成纯方法名称
	pos := strings.LastIndex(serviceAndMethodName, ".")
	if pos > 0 && pos < len(serviceAndMethodName) {
		methodName = serviceAndMethodName[pos+1:] //避免越界
	}
	//	if len(methodName) == 0 {
	//		//参数异常的情况下
	//	}

	rt := []byte{} //最终返回的字节切片
	//TODO检查方法名称是否存在
	switch methodName {
	case "Add":
		//接收参数并处理
		arithReq := new(arithproto.ArithRequest)
		arithReq.Unmarshal(body)            //解参数的pb包
		arithResp := arithobj.Add(arithReq) //计算
		plog.Debug("return", arithResp.A3)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	case "Multiply":
		//接收参数并处理
		arithReq := new(arithproto.ArithRequest)
		arithReq.Unmarshal(body)                      //解参数的pb包
		arithResp := arithobj.Multiply(arithReq)      //计算
		plog.Debug("Multiply returns:", arithResp.A3) //日志
		rt, _ = arithResp.Marshal()                   //返回数据的pb格式化
	default:
		//rt := []byte{}
	}
	//panic("here comes a panic!")
	dis.client.SendPacket(req, rt) //处理完毕之后返回数据！

}

func main() {
	plog.Info("开始:")
	plog.Debug("结束:")
	plog.Warn("warn等级")
	plog.Fatal("fatal等级")
	plog.Core("core等级")

	runtime.GOMAXPROCS(runtime.NumCPU())
	defer plog.CatchPanic()
	quitFlag = make(chan bool)
	//实例化rpc
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		return
	}

	err = client.RpcInit("")

	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}

	mydisp := new(MyDispatch)
	mydisp.client = client
	//设置本rpc服务的代理人（Net层）的ip和端口，并启动服务！
	err = client.LaunchRpcClient(mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	//panic("the main panic")

	//控制退出
	_ = <-quitFlag //收消息退出！
}
