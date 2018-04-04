package rpcclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"framework/rpcclient/bgf"
	"framework/rpcclient/net"
	"framework/rpcclient/srpc"
	"framework/rpcclient/szk"
	"framework/rpcclient/szmq"
	"putil/byteorder"
	"putil/log"
	"putil/timer"
	"sync"
	"time"
	//"sync/atomic"
	//"fmt"
	_ "putil/format"
	//"reflect"
	//"strings"
	//"unicode"
	//"unicode/utf8"
)

const (
	DIS_CLOSE     int32 = 0 //连接关闭状态
	DIS_CONNECTED int32 = 1 //连接但未注册
	DIS_REGISTE   int32 = 2 //注册成功
)

const MAXPACKEGLEN = 80960

const iNVALID_SEQUENCE = 0xffffffff

var globalTimer timer.Timer

func init() {
	globalTimer = timer.NewTimer()
}

/*反射使用，暂停
//方便业务注册和调用
type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

//方便业务注册和使用
type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	numCalls   uint
}

//方便业务注册和使用
// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
*/

type RpcClient struct {
	isexit bool        //退出标记
	net    netstreamer //网络流
	rev    Dispatch    //收到rpc请求的回调接口

	//主要用于注册
	svrtype   int64   //当前业务的typeid
	svrid     int64   //当前业务的severid
	group_ids []int64 //当前业务的groupid

	status int32 //当前运行状态

	seqMutex sync.Mutex //对序号增量加锁
	seq      uint64     //rpc请求序号
	flowid   uint64     //有可能开多个RPCClient，目前暂不使用

	reqRWMutex sync.RWMutex //对请求记录进行添加的锁
	//reqRMutex	sync.RMutex	//对请求记录进行读的锁
	//请求记录
	reqMap map[uint64]*rpcRequestRecord

	regnotify chan bool     //注册通知channel，只有注册成功后才能进行其他业务
	timeroo   timer.TimerOO //计时器

	initflag []byte //初始化标志，主要预防业务使用的人没有调用newRpcClient

	//绑定的本地ip和端口
	localIp       string
	localPort     int
	serverName    string   //服务名称
	serviceName   string   //服务名称
	version       int32    //服务的版本号
	functionsName []string //方法名
	regFuncsFlag  bool     //方法是否已经注册的标记！

	reqMap2 sync.Map

	//serviceMap sync.Map // 业务辅助map[string]*service 反射使用，暂停
}

//完整的包的组成
type WholeMsg struct {
	hd TransHead //网络层头
	bd []byte    //rpc包
}

//类似之前的BGFHead，只是将pb直接修改为了二进制
type TransHead struct {
	Cmd        uint32  //路由层命令字
	TransTypes int32   //传输类型0:p2p,1:p2g,3:byuid
	DSerId     int64   //目的ip
	DSerType   int64   //目的类型
	SSerId     int64   //源ip
	SSerType   int64   //源类型
	Uid        int64   //用户id
	Sessionid  int64   //消息的业务序列号，主要用于日志bug查找
	Extlen     int16   //扩展部分长度 0=无扩展，其它为Ext_Head的长度
	Groupcnt   int16   //扩展信息：组id数量
	Groupids   []int64 //扩展信息：组id
}

//TrandHead中必须要有的长度
var HeadRequiredLen = 58

//此结构体转成传输时字节的“字节总数”,Marshal时需要使用
func (bgfmsg *WholeMsg) Size() (hlen int32) {
	//hd部分的长度
	hlen = int32(HeadRequiredLen) //必须会有的长度
	if bgfmsg.hd.Extlen > 0 {
		hlen += int32(bgfmsg.hd.Extlen)
	}
	hlen += int32(len(bgfmsg.bd)) //加上hd的len
	return
}

//编码整包(结构体转成二进制字节)
func (bgfmsg *WholeMsg) MarshalTo(dAtA []byte) (int, error) {
	//head固定部分
	byteOrderOp.Util_hton_int32(int32(bgfmsg.hd.Cmd), dAtA[0:4])
	byteOrderOp.Util_hton_int32(int32(bgfmsg.hd.TransTypes), dAtA[4:8])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.DSerId, dAtA[8:16])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.DSerType, dAtA[16:24])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.SSerId, dAtA[24:32])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.SSerType, dAtA[32:40])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.Uid, dAtA[40:48])
	byteOrderOp.Util_hton_int64(bgfmsg.hd.Sessionid, dAtA[48:56])
	byteOrderOp.Util_hton_int16(bgfmsg.hd.Extlen, dAtA[56:58])

	//head非固定部分
	var tmped = 58
	if bgfmsg.hd.Extlen > 0 {
		tmped = 60
		byteOrderOp.Util_hton_int16(bgfmsg.hd.Groupcnt, dAtA[58:60])
		//TODO:此处可以做一个校验，避免发出去的包不对！
		if bgfmsg.hd.Groupcnt > 0 {
			for i := 0; i < int(bgfmsg.hd.Groupcnt); i++ {
				byteOrderOp.Util_hton_int64(bgfmsg.hd.Groupids[i], dAtA[tmped:tmped+8])
				tmped += 8
			}
		}
	}

	//body部分
	//dAtA[tmped:(tmped + len(bgfmsg.bd))] = bgfmsg.bd
	copy(dAtA[tmped:], bgfmsg.bd)
	return tmped + len(bgfmsg.bd), nil
}

//解码整包(二进制字节切片转成结构体)
func (bgfmsg *WholeMsg) Unmarshal(dAtA []byte) error {
	msglen := len(dAtA)
	plog.Debug("here Unmarshal length is:", msglen)
	//plog.Debug("here Unmarshal data is:", dAtA)
	if msglen < HeadRequiredLen {
		plog.Fatal("the whole length is less than the head length")
		//return fmt.Errorf("the whole length is less than the head length")
	}

	//Cmd
	tCmd, err := byteOrderOp.Util_ntoh_int32(dAtA[0:4])
	if err != nil {
		plog.Fatal("decode Cmd err", err)
	}

	bgfmsg.hd.Cmd = uint32(tCmd)
	//TransTypes
	tTransTypes, err := byteOrderOp.Util_ntoh_int32(dAtA[4:8])
	if err != nil {
		plog.Fatal("decode TransTypes err", err)
	}
	bgfmsg.hd.TransTypes = tTransTypes
	//DSerId
	tDSerId, err := byteOrderOp.Util_ntoh_int64(dAtA[8:16])
	if err != nil {
		plog.Fatal("decode DSerId err", err)
	}
	bgfmsg.hd.DSerId = tDSerId
	//DSerType
	tDSerType, err := byteOrderOp.Util_ntoh_int64(dAtA[16:24])
	if err != nil {
		plog.Fatal("decode DSerType err", err)
	}
	bgfmsg.hd.DSerType = tDSerType
	//SSerId
	tSSerId, err := byteOrderOp.Util_ntoh_int64(dAtA[24:32])
	if err != nil {
		plog.Fatal("decode SSerId err", err)
	}
	bgfmsg.hd.SSerId = tSSerId
	//SSerType
	tSSerType, err := byteOrderOp.Util_ntoh_int64(dAtA[32:40])
	if err != nil {
		plog.Fatal("decode SSerType err", err)
	}
	bgfmsg.hd.SSerType = tSSerType
	//Uid
	tUid, err := byteOrderOp.Util_ntoh_int64(dAtA[40:48])
	if err != nil {
		plog.Fatal("decode Uid err", err)
	}
	bgfmsg.hd.Uid = tUid
	//Sessionid
	tSessionid, err := byteOrderOp.Util_ntoh_int64(dAtA[48:56])
	if err != nil {
		plog.Fatal("decode Sessionid err", err)
	}
	bgfmsg.hd.Sessionid = tSessionid
	//Extlen
	tExtlen, err := byteOrderOp.Util_ntoh_int16(dAtA[56:58])
	if err != nil {
		plog.Fatal("decode Extlen err", err)
	}
	bgfmsg.hd.Extlen = tExtlen
	tmpet := 58 //正常终止POS
	//plog.Debug("the msg is", bgfmsg.hd)
	//存在扩展字段：
	if bgfmsg.hd.Extlen > 0 {
		//Groupcnt
		tGroupcnt, err := byteOrderOp.Util_ntoh_int16(dAtA[58:60])
		if err != nil {
			plog.Fatal("decode Groupcnt err", err)
		}
		bgfmsg.hd.Groupcnt = tGroupcnt
		tmpet = 60
		//存在几个组id
		if bgfmsg.hd.Groupcnt > 0 {
			tmpst := 60
			//tmpet = 68
			for i := 0; i < int(bgfmsg.hd.Groupcnt); i++ {
				gpid, err := byteOrderOp.Util_ntoh_int64(dAtA[tmpst:(tmpst + 8)])
				if err != nil {
					plog.Fatal("decode Groupids err", err)
				}
				tmpst = tmpst + 8
				tmpet = tmpst
				bgfmsg.hd.Groupids = append(bgfmsg.hd.Groupids, gpid) //添加数据
			}
		}
	}
	//plog.Debug("the unmarshal body flag111111:", tmpet, msglen, "======>", bgfmsg.hd)
	bgfmsg.bd = dAtA[tmpet:msglen]
	//plog.Debug("the unmarshal body flag222222:", tmpet, msglen)

	return nil
}

//生成rpc序列号并返回（为防止序列号冲突，需要加锁）
func (rpcclient *RpcClient) getSequence() uint64 {
	rpcclient.seqMutex.Lock()
	//
	defer rpcclient.seqMutex.Unlock()
	//atomic.AddUint64(&rpcclient.seq, 1)	//用原子操作代替同步操作
	rpcclient.seq++
	return rpcclient.seq
}

//从请求表中剔除某个序列号
func (rpcclient *RpcClient) deleteReqRecord(seq uint64) {

	rpcclient.reqRWMutex.Lock()
	defer rpcclient.reqRWMutex.Unlock()

	_, isfind := rpcclient.reqMap[seq]
	if !isfind {
		plog.Debug("deleteReqRecord cann't find reqrecord seq:", seq)
	}
	plog.Debug("deleteReqRecord, seq = ", seq)
	delete(rpcclient.reqMap, seq)
}

//存储请求队列
func (rpcclient *RpcClient) addReqRecord(seq uint64, rpcReqRecd *rpcRequestRecord) {
	//都需要加锁操作
	rpcclient.reqRWMutex.Lock()
	defer rpcclient.reqRWMutex.Unlock()

	old, isfind := rpcclient.reqMap[seq]
	if isfind {
		plog.Debug("addReqRecord find a repeated reqrecord seq:", seq)
		old.close()
	}
	rpcclient.reqMap[seq] = rpcReqRecd
	plog.Debug("client发送请求时的reqMap为:", rpcclient.reqMap)
	return
}

//获取指定序列号的追踪器
func (rpcclient *RpcClient) getReqRecord(seq uint64) (rpcReqRecd *rpcRequestRecord, isfind bool) {
	rpcclient.reqRWMutex.RLock()
	defer rpcclient.reqRWMutex.RUnlock()
	plog.Debug("client响应时的reqMap为", rpcclient.reqMap)
	rpcReqRecd, isfind = rpcclient.reqMap[seq]
	plog.Debug("获取指定序列号的追踪器结果为", isfind)
	return
}

//这里应该说是打印reqMap
func (rpcclient *RpcClient) printRecMap() {
	rpcclient.reqRWMutex.RLock()
	defer rpcclient.reqRWMutex.RUnlock()
	for k, v := range rpcclient.reqMap {
		plog.Debug("k:", k, "rec:status ", v.status, "seq: ", v.seq)
	}
}

//清空reqMap
func (rpcclient *RpcClient) cleanRequestRecords() {
	plog.Debug("清空reqMap")
	rpcclient.reqRWMutex.Lock()
	defer rpcclient.reqRWMutex.Unlock()
	for key, v := range rpcclient.reqMap {
		if v.status == RPCREQ_DISCARD { //清理被抛弃的
			rpcclient.deleteReqRecord(key)
		}
	}
}

//设置svrType
func (rpcclient *RpcClient) setSvrType(serv_type int64) {
	_ = rpcclient.initflag[0]
	rpcclient.svrtype = serv_type
}

//设置svrId
func (rpcclient *RpcClient) setSvrId(id int64) {
	_ = rpcclient.initflag[0]
	rpcclient.svrid = id
}

//设置该服务的版本
func (rpcclient *RpcClient) setVersion(version int32) {
	_ = rpcclient.initflag[0]
	rpcclient.version = version
}

//设置groupId
func (rpcclient *RpcClient) setGroupIds(group_ids []int64) {
	_ = rpcclient.initflag[0]
	rpcclient.group_ids = group_ids
}

//设置serverName
func (rpcclient *RpcClient) setServerName(svrN string) {
	_ = rpcclient.initflag[0]
	rpcclient.serverName = svrN
}

//设置serviceName
func (rpcclient *RpcClient) setServiceName(svcN string) {
	_ = rpcclient.initflag[0]
	rpcclient.serviceName = svcN
}

//设置functionsName
func (rpcclient *RpcClient) setFuncsName(funcsN []string) {
	_ = rpcclient.initflag[0]
	rpcclient.functionsName = funcsN
}

//实例化rpcclient map表预留100万
func newRpcClient() (rpcclient *RpcClient, err error) {
	rpcclient = new(RpcClient)
	if rpcclient == nil {
		return nil, errors.New("alloc memory failed")
	}
	rpcclient.flowid = 0
	rpcclient.reqMap = make(map[uint64]*rpcRequestRecord, 1000000)
	//rpcclient.timer = timer.NewTimer()
	rpcclient.initflag = make([]byte, 1)
	return
}

//设置本地ip、port
func (rpcclient *RpcClient) setLocalAddress(localip string, localport int) {
	rpcclient.localIp = localip
	rpcclient.localPort = localport
}

//注册服务
func (rpcclient *RpcClient) registerSvr() {
	var head TransHead
	head.Uid = 1                     //user_id
	head.Cmd = srpc.DISPATCH_REG_SER //表示向Net层注册
	head.DSerId = -1
	head.DSerType = srpc.DISPATCH
	head.TransTypes = srpc.TRANS_P2P
	//head.MsgType = 1
	head.Groupids = rpcclient.group_ids
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype

	//注册的消息内容(不使用CDispRegMsg)
	var regbody RPCProto.RegisterDispatchReq
	cantreplace := int32(0)
	regbody.CantReplace = &cantreplace    //本服务是否可替换
	regbody.Groupid = rpcclient.group_ids //本服务所属的组

	plog.Debug("server register body is ", regbody)

	regbuf, err := regbody.Marshal()
	if err != nil {
		plog.Fatal("Marshal Regbody err", err)
	}
	/*
		regbuf := make([]byte) //, MAXPACKEGLEN
		_, err := regbody.MarshalTo(regbuf)
		if err != nil {
			plog.Fatal("Marshal Regbody err:", err)
		}
	*/
	plog.Debug("server register head is ", head)
	plog.Debug("server register body is ", regbuf)

	rpcclient.send(&head, regbuf, nil)

	return
}

//注册serverName、serviceName、funcs
func (rpcclient *RpcClient) registerFuncs() {
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.DISPATCH_REG_METHOD_REQ //表示向Net层注册方法
	head.DSerId = -1
	head.DSerType = srpc.DISPATCH
	head.TransTypes = srpc.TRANS_P2P
	//head.MsgType = 1
	head.Groupids = rpcclient.group_ids
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	if rpcclient.svrtype < 1000 {
		plog.Fatal("svrtype is less than 1000 ERR!")
	}

	//注册的消息内容
	var regbody RPCProto.DispatchRegMethodReq

	plog.Debug("&rpcclient :", &rpcclient)
	//赋值由业务端调用操作,调试阶段先手动赋值
	regbody.ServerName = &rpcclient.serverName
	regbody.ServiceName = &rpcclient.serviceName
	regbody.MethodName = rpcclient.functionsName
	plog.Debug("registerFuncs send info :", head, regbody)
	regbuf, err := regbody.Marshal()
	if err != nil {
		plog.Fatal("Marshal registerFuncs info err", err)
	}
	rpcclient.send(&head, regbuf, nil)

	return
}

//（important-接收rpc请求的地方）接收Net发过来的数据，然后配合for循环无限开启goroutine来处理收到的（完整的）包！
func (rpcclient *RpcClient) clientInput() {
	//预分配，防止GC TODO:为何要预分配？
	protolen := make([]byte, 4)
	data := make([]byte, MAXPACKEGLEN)

	rpcclient.isexit = false
	//plog.Debug("rpcclient is ", rpcclient)
	plog.Debug("rpcclient.isexit is ", rpcclient.isexit)
	for !rpcclient.isexit {
		//获取整包长度datalen
		//plog.Debug("rpcclient.isexit = ", rpcclient.isexit)
		if err := rpcclient.net.ReadAtLeast(protolen, 4); err != nil {
			plog.Fatal("io exception  ", err)
			rpcclient.close() //整个服务进程都要结束,是不是太狠了？
			return
		}

		var datalen int32
		datalen, err := byteOrderOp.Util_ntoh_int32(protolen)
		if err != nil {
			plog.Fatal("decode datapackage len err", err)
		}

		//包太大了，丢弃！然后一直在循环，知道超时
		if datalen > MAXPACKEGLEN {
			extdata := make([]byte, datalen)
			//data = extdata
			rpcclient.net.ReadAtLeast(extdata, int(datalen-4))
			plog.Fatal("input.......................package exceed MAXlen")
			continue
		}

		//读取剩余的包
		if err := rpcclient.net.ReadAtLeast(data, int(datalen-4)); err != nil {
			plog.Fatal("io excepton", err)
			rpcclient.close()
			return
		}

		//bgfmsg := new(RPCProto.BGFMsg)
		var bgfmsg WholeMsg
		decodedata := data[:datalen-4]
		//plog.Debug("recv msg size :", datalen, "data:", format.ByteSlice2HexString(decodedata))
		err = bgfmsg.Unmarshal(decodedata)
		if err != nil {
			plog.Fatal("occure decode err!")
			continue
		}

		plog.Debug("iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiinput head: ", getNetHeadString(&bgfmsg.hd))
		go rpcclient.processPacket(&bgfmsg) //收包之后，放到另外的线程去处理！
	}
}

//important 分析收到的各种各样的包，然后分发到各自的处理器上处理！（多线程处理了！）
func (tcpclient *RpcClient) processPacket(bgfmsg *WholeMsg) {
	plog.Debug("processPacket", bgfmsg.hd.Cmd)
	switch bgfmsg.hd.Cmd {
	case srpc.DISPATCH_REG_SER:
		fallthrough
	case srpc.DISPATCH_REG_RSP: //0x106
		tcpclient.processRegisterResp(bgfmsg) //处理服务注册响应包
	case srpc.DISPATCH_SENDINFO_SER: //0x104
		tcpclient.processServiceDelete(bgfmsg) //处理服务冲突，删除包
	case srpc.DISPATCH_REG_METHOD_RSP: //0x202
		tcpclient.processRegisterFuncsResp(bgfmsg) //处理方法注册响应包
	case srpc.RESPONSE_DISPATCH_KEEP_ALIVE_SER:
		tcpclient.processKeepAliveResp(bgfmsg) //处理心跳响应包
	case srpc.RPC_REQUEST_PACKAGE:
		tcpclient.processRpcRequest(bgfmsg) //处理rpc请求包：别的机器向我这里发rpc请求(我需要响应)
	case srpc.RPC_RESPONSE_PACKET:
		tcpclient.processRpcResponse(bgfmsg) //处理rpc响应包：我本地的客户端之前发了个请求出去，现在要在这里处理收到的响应包
	case srpc.DISPATCH_SER_NOT_FIND:
		tcpclient.ProcessSvrNotFind(bgfmsg) //dispatch(Net)回应服务不存在的消息
	case srpc.RPC_NOTIFY_PACKET:
		//return ProcessRpcNotify(data, len)	//处理收到的点对组的无需响应的包(收到点对点无需响应包)
	case srpc.RPC_METHOD_NOTIFY_PACKET:
		//return ProcessRpcMethodNotify(&msg)//处理rpc请求包：别的机器向我这里发rpc请求(收到点对点无需响应包) TODO:模拟正常的点对点有相应的请求即可！
	default:
		//return ProcessQEPacket(data, len)
	}
	//  return 1
}

//接收其他rpc服务发过来的请求数据，然后转给RpcRequest来处理
func (rpcclient *RpcClient) processRpcRequest(bgfmsg *WholeMsg) {
	pRpcHead, rpcBody, err := srpc.SrpcUnpackPkgHeadandBody(bgfmsg.bd)
	if err != nil {
		plog.Fatal(err)
		return
	}
	var req RpcRecvReq = RpcRecvReq{
		Rpchead:   pRpcHead,
		Ssertype:  bgfmsg.hd.SSerType,
		Sserid:    bgfmsg.hd.SSerId,
		smtid:     pRpcHead.MtId, //TODO：MtId bgfmsg.bd.MtId
		Inputtime: time.Now().UnixNano()}
	plog.Debug("server接受请求序列号", req.Rpchead.Sequence)
	rpcclient.rev.RpcRequest(&req, rpcBody) //进入业务
}

//处理dispatch(Net)回应服务不存在的消息
func (rpcclient *RpcClient) ProcessSvrNotFind(bgfmsg *WholeMsg) {

	pRpcHead, rpcBody, err := srpc.SrpcUnpackPkgHeadandBody(bgfmsg.bd)
	if err != nil {
		plog.Fatal(err)
		return
	}
	plog.Debug(GetRpcHeadString(pRpcHead))

	v, isfind := rpcclient.getReqRecord(pRpcHead.Sequence)
	if !isfind {
		plog.Fatal("cann't find request record")
		return
	}
	//rpcclient.deleteReqRecord(pRpcHead.Sequence)
	plog.Debug("返回的状态码说：Net层没找到dst")
	v.returnResponse(pRpcHead, rpcBody, RPC_RESPONSE_TARGET_NOTFOUND) //返回的状态码说：Net层没找到dst
}

//接收其他rpc服务发过来的响应数据！
func (rpcclient *RpcClient) processRpcResponse(bgfmsg *WholeMsg) {
	pRpcHead, rpcBody, err := srpc.SrpcUnpackPkgHeadandBody(bgfmsg.bd)
	if err != nil {
		plog.Fatal(err)
		return
	}
	plog.Debug(GetRpcHeadString(pRpcHead))

	plog.Debug("client响应的请求序列号", pRpcHead.Sequence)
	v, isfind := rpcclient.getReqRecord(pRpcHead.Sequence)
	if !isfind {
		plog.Fatal("cann't find request record")
		return
	}
	plog.Debug("client接受来自服务发过来的响应数据头部为", pRpcHead)
	plog.Debug("client接受来自服务发过来的响应数据body为", rpcBody)
	plog.Debug("标识符为", RPC_RESPONSE_COMPLET)
	v.returnResponse(pRpcHead, rpcBody, RPC_RESPONSE_COMPLET)
}

//Net层说：有一个新的服务，type和id跟你一样，请你退出吧。
func (rpcclient *RpcClient) processServiceDelete(bgfmsg *WholeMsg) {
	//响应内容
	var resp RPCProto.RegisterDispatchResp

	err := resp.Unmarshal(bgfmsg.bd)
	if err != nil {
		plog.Debug("ServiceDelete body err!")
		return
	}
	if 8 != *resp.RetCode {
		plog.Fatal("Register Dis failed ret", *resp.RetCode)
	}
	plog.Debug("Service need deleted!!! returncode is................. ", *resp.RetCode)
	rpcclient.close() //关闭

	return
}

//处理服务注册的响应包！(通过input接收转过来的)
func (rpcclient *RpcClient) processRegisterResp(bgfmsg *WholeMsg) {

	//响应内容
	var resp RPCProto.RegisterDispatchResp

	err := resp.Unmarshal(bgfmsg.bd)
	if err != nil {
		plog.Debug("RegisterResp err!")
		return
	}
	if 0 != *resp.RetCode {
		plog.Fatal("Register Dis failed ret", resp.RetCode)
	}
	plog.Debug("Register succeed!!! returncode is................. ", *resp.RetCode)
	if rpcclient.status == DIS_CONNECTED { //连接成功
		rpcclient.regnotify <- true //通知luanch，服务注册成功了，主线程不用一直等着了！
		close(rpcclient.regnotify)
		rpcclient.status = DIS_REGISTE
	}

	//rpcclient.last_heart_time = time.Now()

	rpcclient.timeroo = globalTimer.StartTimer(30*time.Second, rpcclient)

	return
}

//处理方法注册的响应包！(通过input接收转过来的)
func (rpcclient *RpcClient) processRegisterFuncsResp(bgfmsg *WholeMsg) {
	//响应内容
	var resp RPCProto.DispatchRegMethodResp
	err := resp.Unmarshal(bgfmsg.bd)
	if err != nil {
		plog.Debug("processRegisterFuncsResp Unmarshal err!")
		return
	}
	if 0 != *resp.Result {
		plog.Fatal("Register Dis failed ret", resp.Result)
	}
	plog.Debug("processRegisterFuncsResp succeed!!! returncode is................. ", *resp.Result)
	rpcclient.regFuncsFlag = true //方法注册成功之后，修改一下标志位！

	//rpcclient.timeroo = globalTimer.StartTimer(30*time.Second, rpcclient)

	return
}

//处理Net发过来的保持心跳的操作（收到了就表示保持了）
func (rpcclient *RpcClient) processKeepAliveResp(bgfmsg *WholeMsg) {

	//	rpcclient.last_heart_time = time.Now()
	plog.Debug("receive heart beat response")
}

func (rpcclient *RpcClient) processKeepAlive() {

	plog.Debug("ProcessKeepAlive send")
	//	if time.Since(rpcclient.last_heart_time).Seconds() >= 120 {
	//		plog.Fatal("heart beat lost, last time :", rpcclient.last_heart_time.Format(time.UnixDate))
	//	}
	var msg TransHead

	//非关键信息可以在send中去默认填充
	msg.Uid = 1
	msg.Cmd = srpc.DISPATCH_KEEP_ALIVE_SER
	msg.DSerId = -1
	msg.DSerType = srpc.DISPATCH
	msg.SSerId = rpcclient.svrid
	msg.SSerType = rpcclient.svrtype
	msg.TransTypes = srpc.TRANS_P2P
	//msg.MsgType = 0

	rpcclient.send(&msg, nil, nil)
	plog.Debug("send heat beat")
	rpcclient.timeroo = globalTimer.StartTimer(30*time.Second, rpcclient)
}

//TODO:调用场景？
func (rpcclient *RpcClient) TimeNotify() {

	plog.Debug("receive time notify")
	plog.Debug("tcpclient.status =", rpcclient.status)
	switch rpcclient.status {
	case DIS_CLOSE: //需要重连
		//NGLOG_ERROR("start reconnect dispatch");
		//spp::rpcmodule::CRPCModule::Instance()->ProcCmdWithRpcCall(ReconnDispatch, (void *)this);
		break
	case DIS_CONNECTED: //需要注册未及时回包
		//直接再次发送一次注册包
		//Register_Svr();
		break
	case DIS_REGISTE:
		//暂时关闭心跳
		rpcclient.processKeepAlive()
		break
	}
}

//(important)作为客户端，发送RPC请求，并等待响应
func (rpcclient *RpcClient) sendAndRecvRespRpcMsg(serviceAndMethodName string, dstsvrtype int64, dstsvrid int64, req []byte, timeout int) (ret *RpcResponse) {
	_ = rpcclient.initflag[0]
	//请求头消息
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.RPC_REQUEST_PACKAGE //表示这是一个rpc请求包
	head.DSerId = dstsvrid              //请求目标是谁
	head.DSerType = dstsvrtype          //请求目标是谁
	head.SSerId = rpcclient.svrid       //请求者（我）是谁
	head.SSerType = rpcclient.svrtype   //请求者（我）是谁
	head.TransTypes = srpc.TRANS_P2P    //点对点
	//head.MsgType = 0
	//head.MtId = 0

	//rpchead := new(srpc.CRpcHead)
	//rpchead消息
	var rpchead srpc.CRpcHead
	rpchead.MethodName = []byte(serviceAndMethodName) //请求哪个方法
	rpchead.Sequence = rpcclient.getSequence()        //本次请求的序列号
	plog.Debug("client开始请求的序列号为", rpchead.Sequence)

	//rpchead.FlowId = rpchead.Sequence

	//rpchead.MtId = int32(rpchead.Sequence)
	//head.MtId = int32(rpchead.Sequence)

	//建立追踪器(追踪器主要保持当前的各种状态)
	rpcReqRecd := new(rpcRequestRecord)
	rpcReqRecd.init(&head, &rpchead, rpchead.Sequence, req, timeout, rpcclient)
	//追踪器进入追踪表
	rpcclient.addReqRecord(rpchead.Sequence, rpcReqRecd)
	//开启超时器
	rpcReqRecd.startTimer()

	rpcpack, err := srpc.SrpcPackPkg(&rpchead, req) //组织报文
	if err != nil {
		plog.Fatal(err)
	}

	rpcclient.send(&head, rpcpack, rpcReqRecd) //发送

	ret = rpcReqRecd.wait() //阻塞在Wait里面，一直到数据返回
	plog.Debug("阻塞等待请求序列号为", rpchead.Sequence)
	return
}

//发送不需要回应的rpc请求
func (rpcclient *RpcClient) sendNoRespRpcMsg(method_name string, dstsvrtype int64, dstsvrid int64, req []byte, groupids []int64) {
	_ = rpcclient.initflag[0]
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.RPC_METHOD_NOTIFY_PACKET
	head.DSerId = dstsvrid
	head.DSerType = dstsvrtype
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	head.TransTypes = srpc.TRANS_P2P //点对点
	if groupids != nil {
		head.Groupids = groupids
		head.Groupcnt = int16(len(groupids))
		head.Extlen = int16(head.Groupcnt*8 + 2) //乘以8是一个groupid为8字节，+2是加上Groupcnt的字节数
	}
	//head.MsgType = 0

	var rpchead srpc.CRpcHead
	rpchead.MethodName = []byte(method_name)
	rpchead.Sequence = rpcclient.getSequence()
	//rpcclient.flowid = rpchead.Sequence

	//启动数据追踪，但无需设置定时器
	var rpcReqRecd rpcRequestRecord
	rpcReqRecd.init(&head, &rpchead, rpchead.Sequence, req, 0, rpcclient)

	rpcpack, err := srpc.SrpcPackPkg(&rpchead, req)
	if err != nil {
		plog.Fatal(err)
	}

	rpcclient.send(&head, rpcpack, &rpcReqRecd)
}

//提供给rpc请求返回的接口(回个包就行了)
func (rpcclient *RpcClient) sendPacket(req *RpcRecvReq, data []byte, rpccall *RpcCall) (err error) {
	_ = rpcclient.initflag[0]
	//var head RPCProto.BGFHead
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.RPC_RESPONSE_PACKET
	head.DSerId = req.Sserid
	head.DSerType = req.Ssertype
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	head.TransTypes = srpc.TRANS_P2P
	//head.MsgType = 0
	//head.MtId = req.smtid //为了兼容当前C++的rpc，这里表示c++中实现协成id

	//rpchead.Sequence = rpcclient.getSequence() //对返回添加一个序号
	//rpchead.FlowId = rpchead.Sequence

	rpcdata, err := srpc.SrpcPackPkg(req.Rpchead, data) //组装报文
	if err != nil {
		plog.Debug("SendPacket err1", err)
	}
	//启动追踪
	var rpcReqRecd rpcRequestRecord
	rpcReqRecd.init(&head, req.Rpchead, iNVALID_SEQUENCE, data, 0, rpcclient)
	rpcclient.send(&head, rpcdata, &rpcReqRecd)
	plog.Debug("回包请求序列号", rpcReqRecd.rpchead.Sequence)
	//rpc响应时间上报zmq(此处不放goroutine  1、业务线程本身就会有有多个；2、zmq那边是非阻塞的！)
	var dcrpcmonitor szmq.DcRpcMonitor
	nowtime := time.Now().UnixNano()
	var diffMicroSeconds int64 = (nowtime - req.Inputtime)
	diffMicroSeconds = diffMicroSeconds / 1e3 //此处的结果就是整数了！
	dcrpcmonitor.Channel_id = ""
	dcrpcmonitor.Act_time = time.Now().Unix()
	dcrpcmonitor.Server_name = fmt.Sprint(rpcclient.svrtype)
	dcrpcmonitor.Server_id = fmt.Sprint(rpcclient.svrid)
	dcrpcmonitor.Func_name = string(req.Rpchead.MethodName)
	dcrpcmonitor.Typ = ""
	dcrpcmonitor.Exec_microsecond = diffMicroSeconds
	dcrpcmonitor.Ext_data = ""
	databytes, _ := json.Marshal(dcrpcmonitor)

	//for i := 0; i < 100000; i++ {
	go rpccall.WriteNormalLog("monitor_rpc", string(databytes))
	//}

	return
}

func (rpcclient *RpcClient) sendNotifyMsg(dstsvrtype int64, dstsvrid int64, data []byte) {
	_ = rpcclient.initflag[0]
	//	var head RPCProto.BGFHead
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.RPC_NOTIFY_PACKET
	head.DSerId = dstsvrid
	head.DSerType = dstsvrtype
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	head.TransTypes = srpc.TRANS_P2P
	//head.MsgType = 0
	//head.MtId = 0

	//启动追踪
	var rpcReqRecd rpcRequestRecord
	rpcReqRecd.init(&head, nil, iNVALID_SEQUENCE, data, 0, rpcclient)

	rpcclient.send(&head, data, &rpcReqRecd)
}

//发送消息给客户端，目标地址必须是NAT
func (rpcclient *RpcClient) sendMsgToClient(clientid int32, dstsvrtype int64, dstsvrid int64, msgdata []byte) (err error) {
	_ = rpcclient.initflag[0]
	//	var head RPCProto.BGFHead
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.CLIENT_TRANSFORM_PACKET
	head.DSerId = dstsvrid
	head.DSerType = dstsvrtype
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	head.TransTypes = srpc.TRANS_P2P
	//head.MsgType = 0
	//head.MtId = 0

	var client_msg RPCProto.ClientTransMsg
	client_msg.ClientId = clientid
	client_msg.ClientMsg = msgdata
	msg, err := client_msg.Marshal()
	if err != nil {
		plog.Fatal("SendMsgToClient err", err.Error())

	}

	//启动追踪
	var rpcReqRecd rpcRequestRecord
	rpcReqRecd.init(&head, nil, iNVALID_SEQUENCE, msgdata, 0, rpcclient)

	rpcclient.send(&head, msg, &rpcReqRecd)
	return
}

//发送点对组广播消息
func (rpcclient *RpcClient) sendGroupNotifyMsg(dstsvrtype int64, groupid int64, data []byte) {
	_ = rpcclient.initflag[0]
	//	var head RPCProto.BGFHead
	var head TransHead
	head.Uid = 1
	head.Cmd = srpc.RPC_NOTIFY_PACKET
	head.DSerId = 0
	head.DSerType = dstsvrtype
	head.SSerId = rpcclient.svrid
	head.SSerType = rpcclient.svrtype
	head.Groupids = make([]int64, 1)
	if groupid != 0 {
		head.Groupids[0] = groupid
		head.Groupcnt = 1
		head.Extlen = int16(head.Groupcnt*8 + 2) //乘以8是一个groupid为8字节，+2是加上Groupcnt的字节数
	}

	head.TransTypes = srpc.TRANS_P2G
	//head.MsgType = 0
	//head.MtId = 0

	//启动追踪
	var rpcReqRecd rpcRequestRecord
	rpcReqRecd.init(&head, nil, iNVALID_SEQUENCE, data, 0, rpcclient)

	rpcclient.send(&head, data, &rpcReqRecd)
}

//发送数据
func (rpcclient *RpcClient) send(head *TransHead, body []byte, rpcReqRecd *rpcRequestRecord) {
	bfgmsg := new(WholeMsg)
	bfgmsg.hd = *head
	bfgmsg.bd = body
	//plog.Debug("bfgmsg.hd", format.ByteSlice2HexString(bfgmsg.hd))
	//plog.Debug("ooooooooooooooooooooooooooooooooooutput head is: ", getNetHeadString(&bfgmsg.hd))
	//plog.Debug("ooooooooooooooooooooooooooooooooooutput : ", bfgmsg)
	go rpcclient.clientOutput(bfgmsg, rpcReqRecd)
}

//rpcclient出数据的地方
func (rpcclient *RpcClient) clientOutput(bfgmsg *WholeMsg, rpcReqRecd *rpcRequestRecord) {
	var lenbuf []byte
	//var lenbuf [4] byte

	//WholeMsg对应传输时的字节大小
	datalen := bfgmsg.Size()
	//plog.Debug("sendrpcclient", rpcReqRecd.rpchead.Sequence)
	if datalen > int32(len(lenbuf)) { //要发出的字节数大于MAXPACKEGLEN
		lenbuf = make([]byte, 4, datalen+4)
	} else {
		lenbuf = make([]byte, 4, MAXPACKEGLEN)
	}
	//plog.Debug("output：..........................the begin lenbuf is: ", lenbuf)
	//plog.Debug("output：..........................the cap is: ", cap(lenbuf))
	byteOrderOp.Util_hton_int32(int32(datalen+4), lenbuf) //第一个4字节存入 TODO?
	//plog.Debug("output：..........................the lenbuf is: ", lenbuf)
	lenbuf1 := lenbuf[4 : datalen+4]
	//plog.Debug("lenbuf1 is:", lenbuf1)
	//plog.Debug("WholeMsg before is:", bfgmsg)
	_, err := bfgmsg.MarshalTo(lenbuf1) //把bfgmsg编码到lenbuf1中
	if err != nil {
		plog.Fatal("Marshall err :", err)
	}
	lenbuf = lenbuf[:datalen+4]
	//plog.Debug("WholeMsg is:", bfgmsg)
	//plog.Debug("lenbuf1 is:", lenbuf1)
	//plog.Debug("send msg size :", len(lenbuf), "data:", lenbuf)
	//plog.Debug("send msg size :", len(lenbuf), "data:", format.ByteSlice2HexString(lenbuf))
	//plog.Debug("output：..........................the lastest lenbuf is:", bfgmsg)
	perr := rpcclient.net.Write(lenbuf)
	//plog.Debug("对Rpc请求需要返回的标识为", perr)
	if perr != nil {
		plog.Debug("NET ERR", getNetHeadString(rpcReqRecd.bghead))
		plog.Fatal("NET ERR,", perr.Code(), perr.Err())
		if rpcReqRecd.bghead.Cmd == srpc.RPC_REQUEST_PACKAGE {
			//对Rpc请求需要返回的 当出现网络错误时 及时返回
			plog.Debug("对Rpc请求需要返回的 当出现网络错误时 及时返回")
			rpcReqRecd.returnResponse(nil, nil, RPC_RESPONSE_NETERR)
		}
		if rpcReqRecd.rpchead != nil {
			plog.Debug("NET ERR", GetRpcHeadString(rpcReqRecd.rpchead))
		}

		rpcReqRecd = nil //等待垃圾回收处理
	}
}

func (tcpclient *RpcClient) exit() {
	tcpclient.isexit = true
}

func (tcpclient *RpcClient) reconnect() {

}

//建立与Net层的连接并发起注册请求（包括服务注册和方法注册）
func (rpcclient *RpcClient) creatNet(remoteip string, remoteport int) (err error) {
	if remoteip == "" || remoteport < 1024 {
		err = errors.New("invalid address")
	}

	rpcclient.net = new(rpcnet.TcpNet)
	rpcclient.net.SetReconTimes(2)
	err = rpcclient.net.Connect(rpcclient.localIp, rpcclient.localPort, remoteip, remoteport)
	plog.Debug("beginn connetced", err)
	if err != nil {
		plog.Fatal("RpcClient.createNet", err)
		return
	}
	rpcclient.regnotify = make(chan bool)
	plog.Debug("beginn register server")
	plog.Debug("rpcclient is", rpcclient)
	rpcclient.registerSvr()

	return
}

//func sendData(data []byte) {

//}

//发起连接(来自rpccall.go中的LaunchRpcClient的调用)
func (rpcclient *RpcClient) launchRpcClient(remoteip string, remoteport int, rev Dispatch, szkclient *szk.SzkClient) (err error) {
	_ = rpcclient.initflag[0]
	err = rpcclient.creatNet(remoteip, remoteport)
	plog.Debug("it has connected and regisert server", err)
	if err != nil {
		return
	}
	rpcclient.rev = rev
	plog.Debug("rpcclient.rev is ", rpcclient.rev)
	plog.Debug("接收数据 ")
	go rpcclient.clientInput() //接收数据
	//	sendChannel := make(chan *BGFMsg)
	//	tcpclient.sendChannel = sendChannel

	rpcclient.status = DIS_CONNECTED //置为连接成功状态
	_ = <-rpcclient.regnotify        //等待注册成功的通知，收不到通知就一直堵塞在这里！
	//plog.Debug("服务注册begin", rpcclient)
	rpcclient.registerFuncs() //0x101服务注册成功之后，开始0x201方法注册
	//plog.Debug("方法注册成功", rpcclient)
	//plog.Debug("++++++++++++++++++++++++++++", rpcclient)
	//plog.Debug("++++++++++++++++++++++++++++", rpcclient.version)
	go szkclient.NodeRegister(rpcclient.serverName, rpcclient.functionsName, rpcclient.svrtype, rpcclient.svrid, rpcclient.version) //向zk注册数据
	//plog.Debug("szkclient has been", rpcclient)
	return
}

//关闭服务
func (rpcclient *RpcClient) close() {
	rpcclient.exit()
	rpcclient.net.Close()
	rpcclient.timeroo.Close()
	rpcclient.initflag = nil

	//关闭残存的计时器
	rpcclient.reqRWMutex.RLock()
	defer rpcclient.reqRWMutex.RUnlock()
	for k, v := range rpcclient.reqMap {
		v.close()
		delete(rpcclient.reqMap, k)
	}
	//==============
	plog.Debug("rpcclient Close")
}

/******************为Rpcrequest服务，方便注册方法的反射调用和处理**********************/
/* 暂停反射的使用
func (rpcclient *RpcClient) registerFuncs(rcvr interface{}) error {
	return rpcclient.register(rcvr, "", false)
}

func (rpcclient *RpcClient) register(rcvr interface{}, name string, useName bool) error {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		plog.Debug(s)
		return errors.New(s)
	}
	if !isExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		plog.Debug(s)
		return errors.New(s)
	}
	s.name = sname

	// Install the methods
	s.method = suitableMethods(s.typ, true)

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		plog.Debug(str)
		return errors.New(str)
	}
	plog.Debug("register service:", sname)
	if _, dup := rpcclient.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			if reportErr {
				plog.Debug("method", mname, "has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				plog.Debug(mname, "argument type not exported:", argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				plog.Debug("method", mname, "reply type not a pointer:", replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				plog.Debug("method", mname, "reply type not exported:", replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				plog.Debug("method", mname, "has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				plog.Debug("method", mname, "returns", returnType.String(), "not error")
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

//处理实例、方法，准备调用
func (rpcclient *RpcClient) callFuncs(req *RpcRecvReq, body []byte) (err error) {
	dot := strings.LastIndex(req.Rpchead.MethodName, ".") //req.Rpchead.MethodName 为一个string，存调用的方法
	if dot < 0 {
		err := errors.New("rpc: service/method request ill-formed: " + req.Rpchead.MethodName)
		plog.Debug(err)
		return nil
	}
	serviceName := req.Rpchead.MethodName[:dot]
	methodName := req.Rpchead.MethodName[dot+1:]
	plog.Debug("service name is:", serviceName, "method name is:", methodName)
	// Look up the request.
	svci, ok := rpcclient.serviceMap.Load(serviceName)
	if !ok {
		err := errors.New("rpc: can't find service " + serviceName)
		plog.Debug(err)
		return nil
	}
	svc := svci.(*service)          //获取service实例
	mtype := svc.method[methodName] //获取方法
	//对参数进行pb解码
	argv := reflect.New(mtype.ArgType) //TODO:获取参数
	mv := argv.MethodByName("Unmarshal")
	mv.Call([]reflect.Value{reflect.ValueOf(body)})
	//argv.Interface().Unmarshal(body)

	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{argv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()

	plog.Debug(errInter)
	return nil
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
*/
