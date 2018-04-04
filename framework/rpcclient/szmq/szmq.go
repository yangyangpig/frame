package szmq

import (
	"errors"
	"fmt"
	"putil/rlog"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	zmq "github.com/pebbe/zmq4"
)

var (
	Nmutex sync.Mutex //zmq的api非线程安全，需要锁控制
	Rmutex sync.Mutex //zmq的api非线程安全，需要锁控制
	Dmutex sync.Mutex //zmq的api非线程安全，需要锁控制

	//	Npull        = []string{"tcp://192.168.96.156:9898"} //普通日志pull服务端
	//	Rpull        = []string{"tcp://192.168.96.156:9897"} //实时日志pull服务端
	//	Dpull        = []string{"tcp://192.168.96.156:9899"} //调试日志pull服务端
	//	sndHighWater = 10000                                 //发送队列长度限制
	//	rcvHighWater = 10000                                 //接收队列长度限制
	//	appId        = "dfqp"                                //上报appid
	//	appKey       = "ff3d48fea9adde39"                    //上报appkey
	//	separator    = "\t"                                  //上报时数据分隔符
	//	Logger       *SzmqPushClient
	Npull        []string //普通日志pull服务端
	Rpull        []string //实时日志pull服务端
	Dpull        []string //调试日志pull服务端
	sndHighWater int      //发送队列长度限制
	rcvHighWater int      //接收队列长度限制
	appId        string   //上报appid
	appKey       string   //上报appkey
	separator    = "	"    //上报时数据分隔符 制表符
	Logger       *SzmqPushClient
)

const (
	ZMQCONN_CLOSE     int32 = 0 //连接关闭状态
	ZMQCONN_CONNECTED int32 = 1 //已连接
)

//push客户端结构体
type SzmqPushClient struct {
	Nsoc *zmq.Socket //普通日志socket
	Rsoc *zmq.Socket //实时日志soc
	Dsoc *zmq.Socket //调试日志soc

	Ncstatus int32 //普通soc连接状态  0表示未连接、1表示连接成功
	Rcstatus int32 //普通soc连接状态
	Dcstatus int32 //普通soc连接状态

	//conMutex sync.Mutex //连接锁
	conMutex sync.RWMutex //连接锁
}

//配置初始化
func init() {
	iniconf, err := config.NewConfig("ini", "./conf/zmq.conf")
	if err != nil {
		rplog.Debug("new config adapter err: ", err)
	}

	Npull = iniconf.Strings("zmq.Npull")
	Rpull = iniconf.Strings("zmq.Rpull")
	Dpull = iniconf.Strings("zmq.Dpull")
	sndHighWater, _ = iniconf.Int("zmq.sndHighWater")
	rcvHighWater, _ = iniconf.Int("zmq.rcvHighWater")
	appId = iniconf.String("zmq.appId")
	appKey = iniconf.String("zmq.appKey")
	//separator = iniconf.String("zmq.separator")

	//plog.Debug("ssssssssssssssssssssssssssssssszmq ", sndHighWater)
}

//实例化
func NewSzmqPushClient() (szmqpushclient *SzmqPushClient, err error) {
	szmqpushclient = new(SzmqPushClient)
	if szmqpushclient == nil {
		return nil, errors.New("alloc memory failed")
	}
	Logger = szmqpushclient
	return
}

//连接pull端
//type为N/R/D中的某一个,pull为连接的pull的地址切片
func (szmqpushclient *SzmqPushClient) setConnect(typ string, pull []string) (err error) {
	if typ != "N" && typ != "R" && typ != "D" {
		err = errors.New("zmq connect err! illegal type!")
		return
	}

	szmqpushclient.conMutex.Lock() //锁
	defer szmqpushclient.conMutex.Unlock()

	//连接的状态指针
	var status *int32
	var soc *zmq.Socket
	switch typ {
	case "N":
		status = &szmqpushclient.Ncstatus
	case "R":
		status = &szmqpushclient.Rcstatus
	case "D":
		status = &szmqpushclient.Dcstatus
	}

	//未连接的情况下，进行连接！
	if *status == 0 {
		bind_to := pull[0] //TODO：暂用第一个（按计划每台机器上只起一个zmq的服务端）
		soc, err = zmq.NewSocket(zmq.PUSH)
		if err != nil {
			rplog.Fatal("NewSocket 2:", err) //socket创建失败
			return
		}

		soc.SetSndhwm(sndHighWater) //设置发送最大队列
		soc.SetRcvhwm(rcvHighWater) //设置接收最大队列
		soc.SetConflate(false)
		err = soc.Connect(bind_to)
		if err != nil {
			rplog.Fatal("s_out.Connect:", err) //连接失败
			return
		}

		//设置socket
		switch typ {
		case "N":
			szmqpushclient.Nsoc = soc
		case "R":
			szmqpushclient.Rsoc = soc
		case "D":
			szmqpushclient.Dsoc = soc
		}
		*status = ZMQCONN_CONNECTED
	}
	return
}

//通过zmq写普通日志
func (szmqpushclient *SzmqPushClient) WriteNormalLog(apiName string, data string) (err error) {
	if apiName == "" {
		return errors.New("apiName is empty")
	}
	//连接pull
	szmqpushclient.conMutex.RLock()
	st := szmqpushclient.Ncstatus
	szmqpushclient.conMutex.RUnlock()
	if st == 0 {
		rplog.Debug("szmqpushclient.Ncstatus = ", st)
		terr := szmqpushclient.setConnect("N", Npull)
		if terr != nil {
			err = terr
			return
		}
	}

	var output string = "N"
	output = output + apiName + "|" + appId + "|" + appKey + "|" + fmt.Sprint(time.Now().UnixNano()/1e6) + separator + data
	Nmutex.Lock()                                           //zmq非线程安全！
	_, err = szmqpushclient.Nsoc.Send(output, zmq.DONTWAIT) //重要设置：zmq.DONTWAIT 表示非阻塞！
	Nmutex.Unlock()                                         //zmq非线程安全！
	if err != nil {
		rplog.Fatal("zmq s_out.Send %v: %v", output, err)
	}
	return
}

//通过zmq写实时日志
func (szmqpushclient *SzmqPushClient) WriteRealLog(apiName string, data string) (err error) {
	if apiName == "" {
		return errors.New("apiName is empty")
	}
	//连接pull
	szmqpushclient.conMutex.RLock()
	st := szmqpushclient.Rcstatus
	szmqpushclient.conMutex.RUnlock()
	if st == 0 {
		//rplog.Debug("szmqpushclient.Rcstatus = ", szmqpushclient.Rcstatus)
		terr := szmqpushclient.setConnect("R", Rpull)
		if terr != nil {
			err = terr
			return
		}
	}
	//rplog.Debug("==============================> ", szmqpushclient.Ncstatus, szmqpushclient.Nsoc)
	var output string = "R"
	output = output + apiName + "|" + appId + "|" + appKey + "|" + fmt.Sprint(time.Now().UnixNano()/1e6) + separator + data
	Rmutex.Lock()
	_, err = szmqpushclient.Rsoc.Send(output, zmq.DONTWAIT) //重要设置：zmq.DONTWAIT 表示非阻塞！
	Rmutex.Unlock()
	if err != nil {
		rplog.Fatal("zmq s_out.Send %v: %v", output, err)
	}
	return
}

//通过zmq写调试日志
func (szmqpushclient *SzmqPushClient) WriteDebugLog(apiName string, data string) (err error) {
	if apiName == "" {
		return errors.New("apiName is empty")
	}
	//连接pull
	szmqpushclient.conMutex.RLock()
	st := szmqpushclient.Dcstatus
	szmqpushclient.conMutex.RUnlock()
	if st == 0 {
		//rplog.Debug("szmqpushclient.Dcstatus = ", szmqpushclient.Dcstatus)
		terr := szmqpushclient.setConnect("D", Dpull)
		if terr != nil {
			err = terr
			return
		}
	}
	var output string = "D"
	output = output + apiName + "|" + appId + "|" + appKey + "|" + fmt.Sprint(time.Now().UnixNano()/1e6) + separator + data
	Dmutex.Lock()
	res, err := szmqpushclient.Dsoc.Send(output, zmq.DONTWAIT) //重要设置：zmq.DONTWAIT 表示非阻塞！
	rplog.Debug("res", res)
	Dmutex.Unlock()
	if err != nil {
		rplog.Fatal("zmq s_out.Send %v: %v", output, err)
	}
	return
}

/*
func main() {
	bind_to := "tcp://192.168.96.156:1897"
	s_out, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		rplog.Fatal("NewSocket 2:", err)
	}

	s_out.SetSndhwm(10) //设置发送最大队列
	s_out.SetRcvhwm(10) //设置接收最大队列
	s_out.SetConflate(false)
	err = s_out.Connect(bind_to)
	if err != nil {
		rplog.Fatal("s_out.Connect:", err)
	}

	message_count := 20

	for j := 0; j < message_count; j++ {
		_, err = s_out.Send("Rhelloworld", zmq.DONTWAIT) //重要设置：zmq.DONTWAIT 表示非阻塞！
		if err != nil {
			rplog.Fatal("s_out.Send %d: %v", j, err)
		}
		//time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 50000)
}
*/
