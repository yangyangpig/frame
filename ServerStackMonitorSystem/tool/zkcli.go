package tool

import (
	_ "fmt"

	"errors"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	ZKCONN_CONNECTED = 1
	ZKCONN_CLOSE     = 0
	SESSIONTIMEOUT   = time.Second * 2
)

type Zookeeper struct {
	conn       *zk.Conn //链接句柄
	connstatus int32
}

var (
	registerRetryTimes int
	zkIpPorts          []string
	zkcli              *Zookeeper
)

func init() {
	//获取zk的配置信息
	iniconf, err := config.NewConfig("ini", "./conf/zk.conf")
	if err != nil {
	}
	registerRetryTimes, _ = iniconf.Int("zk.retrytimes") //读取zk注册重试次数的配置
	zkIpPorts = iniconf.Strings("zk.zkipport")           //读取zk的ip端口配置
}

//实例化zk对象
func NewZookeeper() (zkcli *Zookeeper, err error) {
	zkcli = new(Zookeeper)
	if zkcli == nil {
		return nil, errors.New("alloc memory failed")
	}
	return

}

//zk的链接
func (this *Zookeeper) ZookeeperConnect() (res *Zookeeper, err error) {
	if this.connstatus == ZKCONN_CLOSE {
		//这里面可以加上重连机制进去
		conn, _, cerr := zk.Connect(zkIpPorts, SESSIONTIMEOUT)
		if cerr != nil {
			err = cerr
			return //连接出错
		}
		this.connstatus = ZKCONN_CONNECTED //存储连接状态
		this.conn = conn                   //存储连接
	}
	res = this
	return
}

//增加权限
//func (this *Zookeeper) AddAuth() (err error){
//	return _
//}

//create方法
//func (this *Zookeeper) CreateNode(path string, data []byte)  {

//)

//get方法
func (this *Zookeeper) GetData(path string) (res []byte, err error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	res, _, _ = this.conn.Get(path)
	return
}

//exist方法
func (this *Zookeeper) ExistPath(path string) bool {
	if path == "" {
		return false
	}
	res, _, _ := this.conn.Exists(path)
	return res
}

func (this *Zookeeper) GetChildren(path string) (res []string, err error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	res, _, err = this.conn.Children(path)
	if err != nil {
		return
	}
	return
}
