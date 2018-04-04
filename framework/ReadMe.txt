RPC接入说明文档：
参考：PGDemo/server.go  和 PGDemo/client.go 文件

一、原理说明：
1、业务代码 + framework 共同组成rpc服务；
2、在main中实例化rpc的实例：
	关键代码：
	import("framework/rpcclient/core")
	client, err := rpcclient.NewRpcCall()

3、业务代码需要向framework中注册其自身各种属性信息，包括：
	//设置svrtype、svrid、version、组id
	client.SetSvr(svrtype, svrid, version, []int64{1})
	svrtype：服务类型			服务类型的id，全公司统一分配！
	svrid：服务id				同一svrtype下具体的服务实例的id
	version：服务版本			服务升级使用，通过该值来选择最新版本的服务！
	[]int64{1}：服务的组id		服务组id，用来协助x对组的消息发送！
	
	//设置服务名称，服务的方法名slice
	client.SetSvrNames(svrName, svcName, RegFuncs)
	svrName：服务名称					服务名与svrtype一一对应，全公司统一分配！eg："arith"
	svcName：服务名称					svcName与svrName暂为一致（由c++同学定义）eg:"arith"
	RegFuncs：rpc服务开放的方法名		服务开放的方法名eg：[]string{"Add", "Multiply"} 表示arith中放开了Add和Multiply方法供rpc调用！
	
	//设置本地的ip和端口（可不设置！）
	client.SetLocalAddress("127.0.0.1",5648)
	
4、rpc向Net层注册
	关键代码：
	err = client.LaunchRpcClient(netip, netport, mydisp, client.Szkclient)
	netip ：			表示Net服务的ip（rpc要连的net节点，一般会是本机）eg:"192.168.202.25"
	netport:			表示Net服务监听的端口 eg:7000
	mydisp:				main模块的核心结构体，其中主要的成员是rpccall、主要方法RpcRequest负责处理rpc请求。
	client.Szkclient：	zk客户端，直接传该值即可！
	
	
作为客户端发起rpc调用方式示例：
	request := new(RPCProto.ArithRequest)
	request.A1 = 5
	request.A2 = 45
	req_bytes, err := request.Marshal()
	if err != nil {
		fmt.Println(req_bytes)
	}
	//		"arith"表示服务名
		//		Add表示调用的方法名
		//		req_bytes表示请求的参数（经过了protobuf的marshal后）
		//		5000表示5000毫秒后无响应就超时！
	response := client.SendAndRecvRespRpcMsg("arith", "Add", req_bytes, 5000)
	if response.ReturnCode != 0 {
		//rpc返回结果异常
		//RPC_RESPONSE_COMPLET         = 0 //完成
		//	RPC_RESPONSE_TIMEOUT         = 1 //超时
		//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
		//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {
		arithResp := new(RPCProto.ArithResponse)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return value  = ", arithResp.A3)
	}
	
作为服务端提供rpc服务示例：
var arithobj = new(arith.Arith) //类似serveice
func (dis *MyDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {

	//业务接收数据并开始处理：
	methodName := string(req.Rpchead.MethodName) //方法名称
	rt := []byte{}                               //最终返回的字节切片
	//TODO检查方法名称是否存在
	switch methodName {
	case "Add":
		//接收参数并处理
		arithReq := new(RPCProto.ArithRequest)
		arithReq.Unmarshal(body)            //解参数的pb包
		arithResp := arithobj.Add(arithReq) //计算
		plog.Debug("return", arithResp.A3)
		rt, _ = arithResp.Marshal() //返回数据的pb格式化
	case "Multiply":
		//接收参数并处理
		arithReq := new(RPCProto.ArithRequest)
		arithReq.Unmarshal(body)                      //解参数的pb包
		arithResp := arithobj.Multiply(arithReq)      //计算
		plog.Debug("Multiply returns:", arithResp.A3) //日志
		rt, _ = arithResp.Marshal()                   //返回数据的pb格式化
	default:
		//rt := []byte{}
	}

	dis.client.SendPacket(req, rt)
}
	
	
	
	