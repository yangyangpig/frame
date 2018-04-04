package rpcclient

import (
	"framework/rpcclient/srpc"
	"framework/rpcclient/szk"
	"framework/rpcclient/szmq"

	"putil/log"
	"strconv"
	//"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
)

//statefulMap 存放有状态服务，string为服务名，int暂未使用，未来可以作为取模的值
var (
	statefulMap map[string]int = map[string]int{"userserver.UserServerService": 1, "prop.PropService": 1, "coincache.CoinCacheService": 1}
	svrtype     int64
	svrid       int64
	version     int32
	groupids    []int64
	svrname     string
	RegFuncs    []string //所有注册的方法= {"Add", "Multiply"}
	netIp       string
	netPort     int
)

/**
rpc 请求返回
*/
type RpcRecvReq struct {
	Rpchead   *srpc.CRpcHead //rpchead
	Ssertype  int64          //源servertype ,不容许修改
	Sserid    int64          //源serverid， 不容许修改
	smtid     int32          //源mtid(兼容C++现有业务的mtid) 不容许修改
	Inputtime int64          //请求进入的时间戳（单位为纳秒）
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func (rrq *RpcRecvReq) String() string {
	return fmt.Sprintf("[%s,sstype=%d, ssvrid=%d, smtid=%d]", GetRpcHeadString(rrq.Rpchead), rrq.Ssertype, rrq.Sserid, rrq.smtid)

}

type Dispatch interface {
	RpcRequest(req *RpcRecvReq, body []byte)
}

//RpcResponse的返回码（ReturnCode）
const (
	RPC_RESPONSE_NULL            = -1
	RPC_RESPONSE_COMPLET         = 0 //完成
	RPC_RESPONSE_TIMEOUT         = 1 //超时
	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
	RPC_RESPONSE_NETERR          = 3 //发生网络错误
	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例

)

/*
用户请求rpc时的返回结构
返回后首先得判断ReturnCode
*/
type RpcResponse struct {
	Seq        uint64 //请求时的序号
	ReturnCode int    //返回码
	Err        error  //如果发生错误，错误详情
	//请求时的rpchead
	Head *srpc.CRpcHead
	Body []byte //返回数据
}

type RpcCall struct {
	rpcclient      *RpcClient
	Szkclient      *szk.SzkClient       //zookeeper客户端
	Szmqpushclient *szmq.SzmqPushClient //zmq客户端
	netIp          string
	netPort        int
}

/*
*创建rpc调用接口
 */
func NewRpcCall() (rpccall *RpcCall, err error) {
	rpccall = new(RpcCall)
	rpccall.rpcclient, err = newRpcClient()
	if err != nil {
		plog.Fatal("NewClient err", err)
	}
	//zookeeper客户端
	rpccall.Szkclient, err = szk.NewSzkClient()
	if err != nil {
		plog.Fatal("new szkclient err", err)
	}
	//zeromq push客户端
	rpccall.Szmqpushclient, err = szmq.NewSzmqPushClient()
	if err != nil {
		plog.Fatal("new szmqpushclient err", err)
	}
	return
}

/**
* 初始化rpc服务
 */
func (rpccall *RpcCall) RpcInit(confPath string) error {
	//基础配置文件
	if confPath == "" {
		confPath = "./conf/app.conf"
	}
	iniconf, err := config.NewConfig("ini", confPath)
	if err != nil {
		plog.Fatal("new config adapter err: ", err)
		return err
	}

	//设置本rpc的服务类型、服务id、版本号
	svrtype, _ = iniconf.Int64("app.svrtype")
	svrid, _ = iniconf.Int64("app.svrid")
	tversion, _ := iniconf.Int("app.version")
	svrname = iniconf.String("app.svrname")
	RegFuncs = iniconf.Strings("app.funcNames")
	netIp = iniconf.String("app.netIp")
	netPort, _ := iniconf.Int("app.netPort")

	rpccall.netIp = netIp
	rpccall.netPort = netPort

	if svrtype == 0 || svrid == 0 || tversion == 0 || svrname == "" || len(RegFuncs) == 0 || netIp == "" || netPort == 0 {
		plog.Fatal("some config is illegal,plz check ! svrtype =", svrtype, " |svrid =", svrid, " |tversion = ", tversion, " |svrname =", svrname, " |funcs len = ", len(RegFuncs), "|netIp =", netIp, "|netPort = ", netPort)
		err := errors.New("some config is illegal,plz check !")
		return err
	}

	//配置文件内容二次处理
	version = int32(tversion)
	tgroupids := iniconf.Strings("app.groupids") //需要二次装换一下
	groupids = make([]int64, len(tgroupids))
	for k, v := range tgroupids {
		tv, _ := strconv.Atoi(v)
		groupids[k] = int64(tv)
	}

	plog.Info("Register info is:   svrtype =", svrtype, " |svrid =", svrid, " |version = ", version, " |svrname =", svrname, " |funcs = ", RegFuncs, "|netIp =", netIp, "|netPort = ", netPort)

	//client.SetLocalAddress("192.168.201.80", 55567)
	rpccall.SetSvr(svrtype, svrid, version, groupids)
	rpccall.SetSvrNames(svrname, svrname, RegFuncs) //设置本rpc服务的svrName、svcName和服务开放的方法名！（此处svrName和svcName设置为相同）
	return nil
}

/*
*设置地址 servtype:服务类型；in：服务id
 */
func (rpccall *RpcCall) SetSvr(servtype int64, id int64, version int32, group_ids []int64) {
	rpccall.rpcclient.setSvrType(servtype)
	rpccall.rpcclient.setSvrId(id)
	rpccall.rpcclient.setVersion(version)
	rpccall.rpcclient.setGroupIds(group_ids)
}

/*
*设置地址 serverName/serviceName/functionName
 */
func (rpccall *RpcCall) SetSvrNames(svrN string, svcN string, funcsN []string) {
	rpccall.rpcclient.setServerName(svrN)
	rpccall.rpcclient.setServiceName(svcN)
	rpccall.rpcclient.setFuncsName(funcsN)
}

/**
绑定本地ip_port
*/
func (rpccall *RpcCall) SetLocalAddress(localip string, localport int) {
	rpccall.rpcclient.setLocalAddress(localip, localport)
}

/*
接口：点对点有返回
serviceAndMethodName	string	eg: User.getUserInfo  表示获取User服务的getUserInfo信息
给目标[dststype,dstsvrid]调用方法(methodname),调用发送的数据是req，timeout超时时间。并返回ret
*/
func (rpccall *RpcCall) SendAndRecvRespRpcMsg(serviceAndMethodName string, req []byte, timeout int, hashkey int64) (ret *RpcResponse) { //dstsvrtype int64, dstsvrid int64,
	var servicename string
	var methodname string
	var svrid int64
	//获取服务名称
	plog.Debug("进来方法了")
	pos := strings.LastIndex(serviceAndMethodName, ".")
	if pos > 0 {
		servicename = serviceAndMethodName[0:pos]
		if pos < len(serviceAndMethodName) {
			methodname = serviceAndMethodName[pos+1:] //避免越界
		}
	}
	plog.Debug("servicename是", servicename)
	plog.Debug("methodname是", methodname)
	if len(servicename) == 0 || len(methodname) == 0 {
		//参数异常的情况下
		err := errors.New("调用参数异常！")
		return &RpcResponse{ReturnCode: RPC_RESPONSE_NULL, Err: err} //参数异常的问题
	}

	//获取svrid和svrtype var dstsvrtype, dstsvrid int64 //服务type和id
	typ, ids, err := rpccall.Szkclient.GetSerByNames(servicename, methodname)
	plog.Debug("获取zk数据开始")
	//plog.Info("+++++++++++++++++++++", typ, ids)
	if err != nil {
		plog.Fatal("get zookeeper config err ", err, servicename, methodname)
		return &RpcResponse{ReturnCode: RPC_RESPONSE_NULL, Err: err} //zk获取失败（异常）
	}
	plog.Debug("获取zk数据结束")
	if typ == 0 || len(ids) == 0 {
		plog.Info(" illegal config ", servicename, methodname)
		return &RpcResponse{ReturnCode: RPC_RESPONSE_NULL, Err: err} //zk获取成功，但是数据异常
	}

	//有状态和无状态需要区分，有状态的需要计算hash值对应的svrid，无状态的随机获取
	//plog.Debug(" ======ids: ", ids,servicename,methodname)
	_, exists := statefulMap[servicename]
	if exists {
		var ids_length int
		for _, v := range ids {
			ids_length += len(v)
		}
		//进到此处说明ids的len不为空
		//ids_length := len(ids)
		hashvalue := int64(hashkey % int64(ids_length))
		if hashvalue == 0 {
			hashvalue = int64(ids_length)
		}
		svrid = hashvalue
		//plog.Debug("the stateful len is ", ids_length, "the hashkey is ", hashkey, " the svrid is:", svrid)
	} else {
		svrid, err = getRandomId(ids)
		if err != nil {
			plog.Fatal(" get random id error", servicename, methodname)
			return &RpcResponse{ReturnCode: RPC_RESPONSE_NULL, Err: err} //zk获取成功，但是数据异常
		}
	}
	plog.Debug("已经到这来了")
	return rpccall.rpcclient.sendAndRecvRespRpcMsg(serviceAndMethodName, typ, svrid, req, timeout)
}

/*
点对点无响应(下方原注释仅供参考)
给目标[dststype,dstsvrid]调用方法(methodname),调用发送的数据是req，无需返回
这里groupids广播的时候用，groupids中的字段可以填充多个[servertype+svrsvrid]，给这组服务发送消息，注意，这里只能提供推送，并且目前这个逻辑在Agent是否实现有待确认[暂不推荐使用该groupid]。
目前groupids的设计是基于业务来设计的
*/
func (rpccall *RpcCall) SendNoRespRpcMsg(serviceAndMethodName string, req []byte, timeout int, hashkey int64) (err error) { //dstsvrtype int64, dstsvrid int64,, groupids []int64
	var groupids []int64 //先保留该字段，值给nil
	var servicename string
	var methodname string
	var svrid int64
	//获取服务名称
	pos := strings.LastIndex(serviceAndMethodName, ".")
	if pos > 0 {
		servicename = serviceAndMethodName[0:pos]
		if pos < len(serviceAndMethodName) {
			methodname = serviceAndMethodName[pos+1:] //避免越界
		}
	}
	if len(servicename) == 0 || len(methodname) == 0 {
		//参数异常的情况下
		plog.Info("调用参数异常！")
		err = errors.New("调用参数异常！")
		return //参数异常的问题
	}

	//获取svrid和svrtype var dstsvrtype, dstsvrid int64 //服务type和id
	typ, ids, err := rpccall.Szkclient.GetSerByNames(servicename, methodname)

	if err != nil {
		plog.Fatal("get zookeeper config err ", err, servicename, methodname)
		return //zk获取失败（异常）
	}
	if typ == 0 || len(ids) == 0 {
		plog.Fatal(" illegal config ", servicename, methodname)
		err = errors.New("typ is zero or length of ids is zero")
		return //zk获取成功，但是数据异常
	}

	//有状态和无状态需要区分，有状态的需要计算hash值对应的svrid，无状态的随机获取
	_, exists := statefulMap[servicename]
	if exists {
		var ids_length int
		for _, v := range ids {
			ids_length += len(v)
		}
		//进到此处说明ids的len不为空
		//ids_length := len(ids)
		hashvalue := int64(hashkey % int64(ids_length))
		if hashvalue == 0 {
			hashvalue = int64(ids_length)
		}
		svrid = hashvalue
		plog.Info("the stateful len is ", ids_length, "the hashkey is ", hashkey, " the svrid is:", svrid)
	} else {
		svrid, err = getRandomId(ids)
		if err != nil {
			plog.Fatal(" get random id error", servicename, methodname)
			return //zk获取成功，但是数据异常
		}
	}

	rpccall.rpcclient.sendNoRespRpcMsg(methodname, typ, svrid, req, groupids)
	return
}

/*
给源[ssvrtype, ssvrid]返回rpc请求的响应
注意，这里req来自于请求的rpc请求的时交付给业务的req，包含了请求的[type,id],以及请求的头和对应的mtid，不要修改，否则对方有收不到返回的情况
*/
func (rpccall *RpcCall) SendPacket(req *RpcRecvReq, data []byte) (err error) {
	return rpccall.rpcclient.sendPacket(req, data, rpccall)
}

/**
推送接口，推送接口的数据不是rpc数据，是自定义数据，这里是非rpc消息发送， 给客户端发送主要用此接口
*/
func (rpccall *RpcCall) SendNotifyMsg(dstsvrtype int64, dstsvrid int64, data []byte) {
	rpccall.rpcclient.sendNotifyMsg(dstsvrtype, dstsvrid, data)
}

/*
通过网关服务[dstsvrtype, dstsvrid]给客户端id发送数据msgdata，注意该接口暂不建议使用。
*/
func (rpccall *RpcCall) SendMsgToClient(clientid int32, dstsvrtype int64, dstsvrid int64, msgdata []byte) (err error) {
	return rpccall.rpcclient.sendMsgToClient(clientid, dstsvrtype, dstsvrid, msgdata)
}

/*
启动rpc的调用客户端，连接到Agent,Agent的真实网络地址是【remoteip,remoteport】
对收到的rpc请求通过rev来接口
*/
func (rpccall *RpcCall) LaunchRpcClient(rev Dispatch) (err error) {
	return rpccall.rpcclient.launchRpcClient(rpccall.netIp, rpccall.netPort, rev, rpccall.Szkclient)
}

/**
给指定的dstsvrtype和groupid广播（点对组且无需响应）
*/
func (rpccall *RpcCall) SendGroupNotifyMsg(servicename string, groupid int64, data []byte) { //dstsvrtype int64,
	dstsvrtype, _, err := rpccall.Szkclient.GetSertypeByServiceName(servicename)
	if err != nil {
		//plog.Debug("get zookeeper config err ", err, servicename)
		return
	}
	if dstsvrtype == 0 {
		//plog.Debug(" illegal config ", servicename)
		return
	}

	rpccall.rpcclient.sendGroupNotifyMsg(dstsvrtype, groupid, data)
}

/**
通过zmq发送日志(普通日志)
*/
func (rpccall *RpcCall) WriteNormalLog(apiName string, data string) error {
	return rpccall.Szmqpushclient.WriteNormalLog(apiName, data)
}

/**
通过zmq发送日志（实时日志）
*/
func (rpccall *RpcCall) WriteRealLog(apiName string, data string) error {
	return rpccall.Szmqpushclient.WriteRealLog(apiName, data)
}

/**
通过zmq发送日志（调试日志）
*/
func (rpccall *RpcCall) WriteDebugLog(apiName string, data string) error {
	return rpccall.Szmqpushclient.WriteDebugLog(apiName, data)
}

/*
业务使用：注册方法(反射使用，暂停)
*/
/*
func (rpccall *RpcCall) RegisterFuncs(rcvr interface{}) (err error) {
	return rpccall.rpcclient.registerFuncs(rcvr)
}
func (rpccall *RpcCall) CallFuncs(req *RpcRecvReq, body []byte) (err error) {
	return rpccall.rpcclient.callFuncs(req, body)
}
*/

/**
辅助方法：随机获取一个实例id
*/
func getRandomId(ids map[int32][]int64) (id int64, err error) {
	//备选的高版本和低版本
	var maxver, minver int32
	//可选map为空
	if len(ids) == 0 {
		err = errors.New("ids is empty")
		return
	}

	//选择的备选id
	var choosedIds []int64
	//赋初值
	for k, _ := range ids {
		maxver = k
		minver = k
	}
	//找大版本和小版本
	for k, _ := range ids {
		if k > maxver {
			maxver = k
		}
		if k < minver {
			minver = k
		}
	}
	//选择备选的实例id
	if len(ids[maxver]) > 0 {
		choosedIds = ids[maxver]
	}
	//低版本与高版本不相等，且实例id存在，且实例个数大于高版本实例个数时进入服务备选
	if minver != maxver && len(ids[minver]) > 0 && len(ids[minver]) > len(ids[maxver]) {
		choosedIds = append(choosedIds, ids[minver]...)
	}
	//	fmt.Println(maxver, len(ids[minver]), minver, len(ids[maxver]), choosedIds)
	return getId(choosedIds)
}

/**
辅助方法：从一个切片中随机选择一个id（场景：从zk中找到某种服务的实例id后随机挑选一个）
*/
func getId(choosedIds []int64) (id int64, err error) {
	if len(choosedIds) == 0 {
		err = errors.New("ids is empty")
		return
	}
	k := rand.Intn(len(choosedIds))
	id = choosedIds[k]
	return
}
