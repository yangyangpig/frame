package main

import (
	"dfqp/proto/login"
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
	"strings"
	"dfqp/pg-login/controller"
	"dfqp/pg-login/service"
	"dfqp/lang"
)

var (
	quitFlag chan bool          //退出标记
	loginObj = new(controller.LoginController)
	BindObj = new(controller.BindController)
	passwordObj = new(controller.PasswordController)
	logObj  = new(controller.LogInfoController)
	RegFuncs []string //所有注册的方法= {"Add", "Multiply"}
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
	rt := []byte{} //最终返回的字节切片
	//TODO检查方法名称是否存在
	switch methodName {
	case "login":
		//接收参数并处理
		req := new(pgLogin.LoginRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := loginObj.Login(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	case "resetPwd":
		//接收参数并处理
		req := new(pgLogin.ResetPwdRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := passwordObj.Reset(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	case "guestBindPhone":
		req := new(pgLogin.GuestBindPhoneRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := BindObj.GuestBindPhone(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	case "guestBindWechat":
		req := new(pgLogin.GuestBindWechatRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := BindObj.GuestBindWechat(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	case "wechatBindPhone":
		req := new(pgLogin.WechatBindPhoneRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := BindObj.WechatBindPhone(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	case "getLoginInfo":
		req := new(pgLogin.LogInfoRequest)
		req.Unmarshal(body)            //解参数的pb包
		resp := logObj.GetLogInfo(req)
		resp.Msg = lang.Msg(int(resp.Status), lang.ZH)
		plog.Debug("===mydispatch", resp)
		rt, _ = resp.Marshal() //返回数据的pb格式化
	default:

	}

	dis.client.SendPacket(req, rt) //处理完毕之后返回数据！
}

func main() {
	service.Init()
	runtime.GOMAXPROCS(runtime.NumCPU())
	quitFlag = make(chan bool)
	//实例化rpc
	client := service.Client

	err := client.RpcInit("")
	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}

	mydisp := new(MyServiceDispatch)
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