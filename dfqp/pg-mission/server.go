package main

import (
	"dfqp/pg-mission/controller"
	//"PGMission/app/service"
	"dfqp/proto/mission"
	"framework/rpcclient/core"
	//"github.com/astaxie/beego/config"
	"putil/log"
	"runtime"
	//"strconv"
	"strings"
)

var (
	quitFlag   chan bool                            //退出标记
	missionObj = new(controllers.MissionController) //Server实例，为处理接收到的请求做准备
	svrtype    int64
	svrid      int64
	version    int32
	groupids   []int64
	svrname    string
	RegFuncs   []string //所有注册的方法= {"Add", "Multiply"}
	netIp      string
	netPort    int
)



//包含了rpc实例的结构体
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

	case "MissionList":
		//接收参数并处理
		arithReq := pgMission.MissionListRequest{}
		arithReq.Unmarshal(body) //解参数的pb包
		arithResp := missionObj.MissionList(arithReq)
		plog.Debug("===mydispatch", methodName)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	default:

		plog.Debug("===mydispatch", methodName)
		//rt := []byte{}
	}

	dis.client.SendPacket(req, rt) //处理完毕之后返回数据！
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	quitFlag = make(chan bool)
	//实例化rpc
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		return
	}

	err = client.RpcInit()
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

	//控制退出
	_ = <-quitFlag //收消息退出！
}
