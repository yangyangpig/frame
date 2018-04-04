//保存从本地发起的request记录，在收到其他rpc服务的响应之后，从这里剔除！（它是本地request的追踪者！）
package rpcclient

import (
	"errors"
	_ "framework/rpcclient/bgf"
	"framework/rpcclient/srpc"
	"putil/log"
	"putil/timer"
	"sync"
	"time"
)

//请求的状态
const (
	RPCREQ_NULL = iota
	RPCREQ_WAIT
	RPCREQ_COMPLETE
	RPCREQ_DISCARD
)

var iCloseREC bool = false

type rpcRequestRecord struct {
	bghead      *TransHead
	rpchead     *srpc.CRpcHead
	seq         uint64
	rpcbody     []byte
	timeout     int
	statusMutex sync.Mutex
	status      int
	respchan    chan *RpcResponse
	rpcclient   *RpcClient
	timeroo     timer.TimerOO
}

func (rrr *rpcRequestRecord) init(bghead *TransHead, rpchead *srpc.CRpcHead, seq uint64, rpcbody []byte, timeout int, rpcclient *RpcClient) {
	if iCloseREC {
		return
	}
	rrr.bghead = bghead
	rrr.rpchead = rpchead
	rrr.timeout = timeout
	rrr.seq = seq
	rrr.rpcbody = rpcbody
	rrr.status = RPCREQ_WAIT
	rrr.respchan = make(chan *RpcResponse)
	rrr.rpcclient = rpcclient

}

func (rrr *rpcRequestRecord) wait() *RpcResponse {
	if iCloseREC {
		return &RpcResponse{}
	}
	return <-rrr.respchan
}
func (rrr *rpcRequestRecord) TimeNotify() {
	switch rrr.status {
	case RPCREQ_WAIT:

		_, isfind := rrr.rpcclient.reqMap[rrr.seq]
		if !isfind {
			//plog.Debug("TIMEOUT, not find request record,maybe have complete")
			return
		}
		//通知调用方，已经超时
		plog.Debug("occur time out seq=", rrr.seq)
		rpcResp := new(RpcResponse)
		rpcResp.ReturnCode = RPC_RESPONSE_TIMEOUT
		rpcResp.Seq = rrr.seq
		rpcResp.Head = rrr.rpchead
		rpcResp.Err = errors.New("ocur timeout")
		rrr.respchan <- rpcResp
		close(rrr.respchan)
		rrr.rpcclient.deleteReqRecord(rrr.seq)
	case RPCREQ_COMPLETE:
		//调用已经完成,转入废弃状态，等待上层清理
		//rrr.rpcclient.deleteReqRecord(rrr.seq)
		//rrr.setStatus(RPCREQ_DISCARD)

	}

}

func (rrr *rpcRequestRecord) setStatus(status int) {
	//rrr.statusMutex.Lock()
	//defer rrr.statusMutex.Unlock()
	//rrr.status = status
}
func (rrr *rpcRequestRecord) startTimer() {
	if iCloseREC {
		return
	}
	rrr.timeroo = globalTimer.StartTimer(time.Duration(rrr.timeout)*time.Millisecond, rrr)
}

/*
当code是一个网络错误(这里只指写错误)或者超时错误的时候，rpchead指的是发送时的RPCHEAD， body指发送时的数据
当code是一个完成状态，会填充完成状态的值
*/
func (rrr *rpcRequestRecord) returnResponse(rpchead *srpc.CRpcHead, body []byte, code int) {
	if iCloseREC {
		return
	}
	recresp := new(RpcResponse)
	if rpchead == nil {
		recresp.Head = rrr.rpchead
	} else {
		recresp.Head = rpchead
	}

	if body == nil {
		recresp.Body = rrr.rpcbody
	} else {
		recresp.Body = body
	}
	recresp.Seq = rpchead.Sequence
	recresp.ReturnCode = code
	//plog.Debug(rrr.rpcclient),由于在rpcclient中有map，map中进行了读写锁，因此这里调用非常不安全
	rrr.rpcclient.deleteReqRecord(rrr.seq) //已经收到rpc响应了，这里可以剔除了。
	plog.Info("rrr.respchan <- recresp")
	rrr.respchan <- recresp
	plog.Info("rrr.timeroo.Close()")
	rrr.timeroo.Close()
}

func (rrr *rpcRequestRecord) close() {
	if iCloseREC {
		return
	}
	rrr.timeroo.Close()
}
