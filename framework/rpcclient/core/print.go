/*****
*调试使用
*
 */
package rpcclient

import (
	//	"bytes"
	"bytes"
	"fmt"
	_ "framework/rpcclient/bgf"
	"framework/rpcclient/srpc"
)

func getNetHeadString(bgfmsg *TransHead) string {

	return fmt.Sprintf("BGFMsg: cmd=%x dst[type:%d,id:%d] src[type:%d,id:%d]",
		bgfmsg.Cmd, bgfmsg.DSerType, bgfmsg.DSerId, bgfmsg.SSerType, bgfmsg.SSerId)
}

func GetRpcHeadString(rpchead *srpc.CRpcHead) string {
	return fmt.Sprintf("rpchead: len=%d seq=%d methodname=%s",
		rpchead.Size(), rpchead.Sequence, rpchead.MethodName)
}

func (rrs *RpcResponse) String() string {

	var s bytes.Buffer
	s.WriteString(fmt.Sprintf("seq=%d,returncode=%d", rrs.Seq, rrs.ReturnCode))
	if rrs.Head != nil {
		s.WriteString("\nhead = ")
		s.WriteString(GetRpcHeadString(rrs.Head))
	}
	if rrs.Body != nil {
		s.WriteString("\nbody = ")
		s.WriteString(string(rrs.Body))
	}
	return string(s.Bytes())
}
