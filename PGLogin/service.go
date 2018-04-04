package main

import (
	"Proto/login"
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
	"strconv"
	"strings"
	"PGLogin/app/controllers"
	"github.com/astaxie/beego/config"
	"PGLogin/app/service"
	"Proto/sms"
)

var (
	quitFlag chan bool          //退出标记
	//loginObj = new(controllers.LoginController) //Server实例，为处理接收到的请求做准备
	loginv2Obj = new(controllers.LoginController) //Server实例，为处理接收到的请求做准备
	cmsObj = new(controllers.SmsController) //Server实例，为处理接收到的请求做准备
	svrtype  int64
	svrid    int64
	version  int32
	groupids []int64
	svrname  string
	RegFuncs []string //所有注册的方法= {"Add", "Multiply"}
	netIp    string
	netPort  int
)

//包含了rpc实例的结构体
type MyServiceDispatch struct {
	client *rpcclient.RpcCall
}

//接收rpc请求并计算返回
func (dis *MyServiceDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
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
	case "Login":
		//接收参数并处理
		arithReq := new(pgLogin.LoginRequest)
		arithReq.Unmarshal(body)            //解参数的pb包
		arithResp := loginv2Obj.Login(arithReq)
		plog.Debug("===mydispatch", methodName)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	case "GetCaptcha":
		//接收参数并处理
		arithReq := new(pgSms.GetCaptchaRequest)
		arithReq.Unmarshal(body)            //解参数的pb包
		arithResp := cmsObj.GetCaptcha(arithReq)
		plog.Debug("===mydispatch", methodName)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	default:
		//rt := []byte{}
	}

	dis.client.SendPacket(req, rt) //处理完毕之后返回数据！
}

func main() {
	service.Init()
	runtime.GOMAXPROCS(runtime.NumCPU())
	quitFlag = make(chan bool)
	client := service.Client
	iniconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		plog.Fatal("new config adapter err: ", err)
		return
	}

	//设置本rpc的服务类型、服务id、版本号
	svrtype, _ = iniconf.Int64("app.svrtype")
	svrid, _ = iniconf.Int64("app.svrid")
	tversion, _ := iniconf.Int("app.version")
	svrname = iniconf.String("app.svrname")
	RegFuncs = iniconf.Strings("app.funcNames")
	netIp = iniconf.String("app.netIp")
	netPort, _ := iniconf.Int("app.netPort")

	if svrtype == 0 || svrid == 0 || tversion == 0 || svrname == "" || len(RegFuncs) == 0 || netIp == "" || netPort == 0 {
		plog.Fatal("some config is illegal,plz check ! svrtype =", svrtype, " |svrid =", svrid, " |tversion = ", tversion, " |svrname =", svrname, " |funcs len = ", len(RegFuncs), "|netIp =", netIp, "|netPort = ", netPort)
		return
	}

	//配置文件内容二次处理
	version = int32(tversion)
	tgroupids := iniconf.Strings("app.groupids") //需要二次装换一下
	groupids = make([]int64, len(tgroupids))
	for k, v := range tgroupids {
		tv, _ := strconv.Atoi(v)
		groupids[k] = int64(tv)
	}

	plog.Debug("Register info is:   svrtype =", svrtype, " |svrid =", svrid, " |version = ", version, " |svrname =", svrname, " |funcs = ", RegFuncs, "|netIp =", netIp, "|netPort = ", netPort)

	client.SetSvr(svrtype, svrid, version, groupids)
	client.SetSvrNames(svrname, svrname, RegFuncs) //设置本rpc服务的svrName、svcName和服务开放的方法名！（此处svrName和svcName设置为相同）
	mydisp := new(MyServiceDispatch)
	mydisp.client = client
	//设置本rpc服务的代理人（Net层）的ip和端口，并启动服务！
	err = client.LaunchRpcClient(netIp, netPort, mydisp, client.Szkclient)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}

	//控制退出
	_ = <-quitFlag //收消息退出！
}
