package main

import (
	"framework/rpcclient/core"
	"putil/log"
	"fmt"
	"PGMatch/app/proto"
	"PGMatch/app/controller"
	"PGMatch/app/service"
	"github.com/astaxie/beego/logs"
)

// 分发服务
type MyDispatch struct {
	client *rpcclient.RpcCall
}

// 注册的方法
var RegFuncs = []string{
	"Lists", // 获取比赛列表
}

// dispatch 分发调用命令
func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(rpcclient.GetRpcHeadString(req.Rpchead))
	plog.Debug("receive:", string(body))

	// 接受业务数据并处理
	// 方法名
	methodName := string(req.Rpchead.MethodName)
	// 最终返回的字符切片
	rt := []byte{}
	// TODO 检查方法名是否存在
	switch methodName {
	case "Lists": // 获取比赛列表信息
		// 接受参数并处理
		matchReq := new(proto.ListsRequest)
		// 解析请求参数 body
		matchReq.Unmarshal(body)
		// 创建比赛结构体对象
		var matchObj = new(controller.MatchController)
		// 调用业务代码函数
		matchRes := matchObj.Lists(matchReq)
		// 返回数据格式化为 pb 格式
		responsePb, err := matchRes.Marshal()
		if err != nil {
			plog.Fatal("Match.%v reponse protobuf parsing failed!\n err:", err)
		}
		rt = responsePb
	default:
		plog.Warn("Match.%v is not exist!", methodName)
		rt = []byte{}
	}
	// 发送结果数据给客户端
	dis.client.SendPacket(req, rt)
}

func main()  {
	// 初始化服务
	service.Init()
	// 初始化 RpcClient
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
	}
	// 设置服务 id 和 type （暂时没有具体定义）
	serverId := int64(27)
	serverType := int64(12)

	// 加入 1 组(agent 分组, 便于对不同的分组进行广播或者通知)
	client.SetSvr(serverType, serverId, 1, nil)
	client.SetSvrNames("Match", "Match", RegFuncs)

	myDisp := new(MyDispatch)
	myDisp.client = client
	// 设置代理（agent） IP 和 Host
	err = client.LaunchRpcClient("192.168.202.25", 7000, myDisp, client.Szkclient)
	if err != nil {
		plog.Fatal("Launch RpcClient failed", err)
		return
	}
	plog.Debug("Launch RpcClient succeed!!!")

	// 保持运行
	var controlStr string
	fmt.Scanln(&controlStr)

	if controlStr == "exit" {
		return
	}
}
