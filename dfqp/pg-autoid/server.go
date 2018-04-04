package main

import (
	"dfqp/pg-autoid/app" //业务代码
	"dfqp/proto/autoidpro"
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
	//"strconv"
	"dfqp/pg-autoid/service"
	//"sort"
	//"fmt"
	//"putil/prof"
	"strings"
	//"sync"

	//"github.com/astaxie/beego/config"
)

var (
	quitFlag chan bool            //退出标记
	arithobj = new(app.AutoidApi) //Server实例，为处理接收到的请求做准备
	svrtype  int64
	svrid    int64
	version  int32
	groupids []int64
	svrname  string
	RegFuncs []string //所有注册的方法= {"Add", "Multiply"}
	netIp    string
	netPort  int

	//计数器
	//ct int32

	//mux sync.Mutex
)

//包含了rpc实例的结构体
type MyDispatch struct {
	client *rpcclient.RpcCall
}

//接收rpc请求并计算返回
func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	//	mux.Lock()
	//	ct += 1
	//	plog.Debug("I get a request:", ct)
	//	mux.Unlock()
	//	if ct%1000 == 0 {
	//		plog.Fatal("now ct value is: ", ct)
	//	}
	//	if ct%5000 == 0 {
	//		prof.StopProfile()
	//	}
	//var concurrent_times = 1000 //并发个数 1000
	//var result []int64 = make([]int64, 0, concurrent_times)

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
	case "GetId":
		//接收参数并处理
		arithReq := new(autoidpro.AutoidRequest)
		arithReq.Unmarshal(body)              //解参数的pb包
		arithResp := arithobj.GetId(arithReq) //计算 TODO^^^^^^^^^^
		//		rt = []byte(fmt.Sprintf("%d", arithResp.Bid))
		//		plog.Debug("return result is: ", arithResp.Bid)
		//result = append(result, arithResp.Bid)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	default:
		//rt := []byte{}
	}
	//panic("here comes a panic!")

	//	if len(result) == 1000 {
	//		sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	//	}
	//plog.Fatal("server.go result length is:", len(result))
	dis.client.SendPacket(req, rt) //处理完毕之后返回数据！

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	service.Init()
	defer plog.CatchPanic()
	//prof.StartProfile("", "")
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
	//prof.StopProfile()
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
