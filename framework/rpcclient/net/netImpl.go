//类似conn模块
package rpcnet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"putil/log"
	"putil/net"
	. "putil/perror"
	"putil/perrors"
	"sync"
)

const (
	NETERR_NULL    = iota
	NETERR_CONNECT //网络连接异常
	NETERR_WRITE   //网络写异常
	NETERR_READ    //网络读异常
)

const (
	NETSTATU_NIL = iota
	NETSTATU_LISTEN
	NETSTATU_CONNECTED
	NETSTATU_RW
	NETSTATUS_REXCEP
	NETSTATUS_WEXCEP
	NETSTATUS_CLOSE
)

type TcpNet struct {
	remoteaddr string       //远程ip地址
	localaddr  string       //本地ip地址
	conn       *net.TCPConn //tcp连接
	recontimes int          //重连次数
	status     int          //状态
	wMutex     sync.RWMutex //读写信号锁
}

//获取连接状态（字符串类型）
func (tcpnet *TcpNet) getStatus() string {
	switch tcpnet.status {
	case NETSTATU_NIL:
		return "NETSTATU_NIL"
	case NETSTATU_LISTEN:
		return "NETSTATU_LISTEN"
	case NETSTATU_CONNECTED:
		return "NETSTATU_CONNECTED"
	case NETSTATU_RW:
		return "NETSTATU_RW"
	case NETSTATUS_REXCEP:
		return "NETSTATUS_REXCEP"
	case NETSTATUS_WEXCEP:
		return "NETSTATUS_WEXCEP"
	case NETSTATUS_CLOSE:
		return "NETSTATUS_CLOSE"
	}

	return "NETSTATU_NIL"
}

//设置重连的参数
func (tcpnet *TcpNet) SetReconTimes(times int) {
	tcpnet.recontimes = times
}

//发起网络连接
func (tcpnet *TcpNet) Connect(localIp string, localport int, remoteIP string, remoteport int) (err error) {
	addr := netutil.GetNetAddr(remoteIP, remoteport) //ip:port 型字符串
	plog.Debug(addr)
	tcpnet.remoteaddr = addr
	if localIp == "" || localport == 0 {
		tcpnet.localaddr = ""
	} else {
		tcpnet.localaddr = netutil.GetNetAddr(localIp, localport)
	}

	var isconnect bool = false
	for i := 0; i < tcpnet.recontimes; i++ {
		err = tcpnet.Reconnect()
		if err != nil {
			plog.Fatal("tcpNet.Conect times", i, "failed")
			continue //连不上就继续连，一直到最后结束
		}
		isconnect = true
		break
	}
	//TODO:应该先判断isconnect之后再赋值吧？
	tcpnet.status = NETSTATU_CONNECTED

	if isconnect {
		err = nil
		return
	}
	return
}

//TODO 什么时候调用呢？
func (tcpnet *TcpNet) CofigConn(conn *net.TCPConn) {
	tcpnet.conn = conn
	tcpnet.status = NETSTATU_CONNECTED
}

//至少读取多少字节
func (tcpnet *TcpNet) ReadAtLeast(data []byte, size int) (perr Perror) {
	//状态检查
	if tcpnet.status != NETSTATU_CONNECTED {
		plog.Fatal("Net error!, cann't read data netstatus=", tcpnet.getStatus())

		err := errors.New("netstatus err")
		perr = perrors.New2(err)
		perr.SetCode(NETERR_CONNECT)
		perr.SetErr(err)
		return
	}
	//m := 0
	//for m < size {
	//	n, err0 := tcpnet.conn.Read(data[m:size])
	//	if err0 != nil {
	//		plog.Fatal("tcpNet.ReadAtLeast", err0)
	//		perr = perrors.New2(err0)
	//		perr.SetCode(NETERR_READ)
	//		perr.SetErr(err0)
	//	}
	//	m += n
	//}

	//字节切片装不下的时候报错
	if len(data) < size {
		return perrors.New(-1, "short write")
	}
	n := 0
	var err error = nil
	//读不满就一直读
	for n < size && err == nil {
		var nn int
		nn, err = tcpnet.conn.Read(data[n:size])
		//plog.Debug("read size is: >>>>>>>>>>>>>>>>>>>>>>>>>>>>>", nn)
		n += nn
	}
	if n == size {
		err = nil
	} else if n > 0 && err == io.EOF { //读到结尾了
		err = io.ErrUnexpectedEOF
		perr = perrors.New2(err)
		perr.SetCode(NETERR_READ)
		perr.SetErr(err)
	}
	//plog.Debug("realreadlength=", n)

	//if _, err0 := io.ReadAtLeast(tcpnet.conn, data, size); err0 != nil {
	//	tcpnet.status = NETSTATUS_REXCEP
	//	plog.Fatal("tcpNet.ReadAtLeast", err0)
	//	err = perrors.New2(err0)
	//	err.SetCode(NETERR_READ)
	//	err.SetErr(err0)
	//}
	return
}

//往连接中写入数据
func (tcpnet *TcpNet) Write(data []byte) (err Perror) {
	if tcpnet.status != NETSTATU_CONNECTED {
		plog.Fatal("Net error!, cann't write data netstatus=", tcpnet.getStatus())
		return
	}

	if data == nil {
		plog.Fatal("TcpNet.Write paramete invalid")
		err = perrors.New(0, "TcpNet.Write paramete invalid")
		return
	}
	dlen := len(data)

	//写数据时加锁控制
	tcpnet.wMutex.Lock()
	defer tcpnet.wMutex.Unlock()

	wlen, err0 := tcpnet.conn.Write(data)
	if wlen < dlen {
		//写异常的处理,原本要写入的和实际写入的不一致
		tcpnet.status = NETSTATUS_WEXCEP
		fmt.Println("net err happened :", err, "wlen is :", wlen, "dlen is: ", dlen)
		//通过防火墙控制测试：当reject掉rpc与Net层的收发后，再恢复时，时而会出现wlen小于dlen的情况，且此时err0为nil
		plog.Fatal("tcpNet.Write", err)
		err.SetErr(err0)
		err.SetCode(NETERR_WRITE)
		return
	}
	//	for curlen < dlen {
	//		wlen, err0 := tcpnet.conn.Write(data[curlen:])
	//		if err0 != nil {
	//			tcpnet.status = NETSTATUS_WEXCEP
	//			plog.Fatal("tcpNet.Write", err)
	//			err.SetErr(err0)
	//			err.SetCode(NETERR_WRITE)
	//			return
	//		}
	//		curlen += wlen
	//	}

	//plog.Debug("flush data happened, the length is :", wlen)
	bufio.NewWriter(tcpnet.conn).Flush() //// Flush 将缓存中的数据提交到底层的 io.Writer 中
	return
}

//TCP真正的建立连接的地方
func (tcpnet *TcpNet) Reconnect() (err error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", tcpnet.remoteaddr)
	if err != nil {
		plog.Fatal("tcpNet.Reconnect:ResolveTCPAddr retmote addr err", err)
		return
	}
	tcplocaladdr, err := net.ResolveTCPAddr("tcp", tcpnet.localaddr)
	if err != nil {
		plog.Fatal("tcpNet.Reconnect:ResolveTCPAddr localaddress err！", err)
		tcplocaladdr = nil
	}

	//真正的发起TCP连接
	tcpconn, err := net.DialTCP("tcp", tcplocaladdr, tcpaddr)
	if err != nil {
		plog.Fatal("tcpNet.Reconnect:DialTCP", err)
		return
	}
	tcpconn.SetReadBuffer(8096)  //设置该连接的接收缓冲
	tcpconn.SetWriteBuffer(8096) //设置发送缓冲

	tcpnet.conn = tcpconn
	tcpnet.localaddr = tcpconn.LocalAddr().String() //LocalAddr返回本地网络地址（返回的是一个结构体，需要经过String转换成字符串）
	//plog.Debug("local addr:", tcpnet.localaddr)

	return
}

//获取连接中远程的IP、port
func (tcpnet *TcpNet) GetRomteAddr() (ip string, port int, err error) {
	if tcpnet.remoteaddr == "" {
		err = errors.New("No connect and not remoteaddr")
	}
	ip, port, err = netutil.GetIPandPortByAddr(tcpnet.remoteaddr)
	return
}

//获取连接中本地的IP、port
func (tcpnet *TcpNet) GetLocalAddr() (ip string, port int, err error) {
	if tcpnet.localaddr == "" {
		err = errors.New("No connect and not localaddr")
	}
	ip, port, err = netutil.GetIPandPortByAddr(tcpnet.localaddr)
	return
}

//关闭连接
func (tcpnet *TcpNet) Close() {
	tcpnet.conn.Close()
	tcpnet.status = NETSTATUS_CLOSE
	return
}
