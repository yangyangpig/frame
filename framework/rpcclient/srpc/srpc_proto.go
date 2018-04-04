package srpc

import (
	"bytes"
	"errors"
	"putil/byteorder"
	"putil/log"
)

const (
	PROTO_RPC_STX         = 0x28 //表示rpc包的开头(一个字节)
	PROTO_RPC_ETX         = 0x29 //表示rpc包的结尾(一个字节)
	PROTO_HEAD_MIN_LEN    = 10
	PROTO_HEAD_OFFSET     = 9
	PROTO_BODY_LEN_OFFSET = 5
	PROTO_HEAD_LEN_OFFSET = 1
)

//组装报文
//rpc报文格式：PROTO_RPC_STX head_len body_len head body PROTO_RPC_ETX
func SrpcPackPkg(head *CRpcHead, body []byte) (rpcreq []byte, err error) {

	if nil == head || nil == body {
		err = errors.New("SrpcPackPkg parameter err")
	}
	headlen := int32(head.Size())
	bodylen := int32(len(body))

	var buf bytes.Buffer

	buf.WriteByte(PROTO_RPC_STX)
	//lenbuf := make([]byte, 4)
	var lenbuf [4]byte
	byteOrderOp.Util_hton_int32(headlen, lenbuf[:])
	buf.WriteString(string(lenbuf[:]))

	byteOrderOp.Util_hton_int32(bodylen, lenbuf[:])
	buf.WriteString(string(lenbuf[:]))

	headbuf, err := head.Marshal()
	if err != nil {
		plog.Fatal("SrpcPackPkg parameters head invalid", err)
		return
	}
	buf.WriteString(string(headbuf))
	buf.WriteString(string(body))
	buf.WriteByte(byte(PROTO_RPC_ETX))

	rpcreq = buf.Bytes()
	return
}

//没有body的rpc报文格式，最小长度是PROTO_HEAD_MIN_LEN + rpcheadlen
//PROTO_RPC_STX head_len body_len head PROTO_RPC_ETX
func SrpcPackPkgNoBody(head *CRpcHead) (rpcreq []byte, err error) {
	if nil == head {
		err = errors.New("SrpcPackPkgNoBody parameter err")
	}

	headlen := int32(head.Size())
	msglen := PROTO_HEAD_MIN_LEN + headlen

	msgbuf := make([]byte, msglen)
	//var msgbuf [msglen]byte
	buf := bytes.NewBuffer(msgbuf)
	buf.WriteByte(PROTO_RPC_STX)

	//lenbuf := make([]byte, 4)
	var lenbuf [4]byte

	byteOrderOp.Util_hton_int32(headlen, lenbuf[:])
	buf.WriteString(string(lenbuf[:]))
	buf.WriteRune(rune(0))
	headbuf, err := head.Marshal()
	if err != nil {
		plog.Fatal("SrpcPackPkg parameters head invalid", err)
		return
	}
	buf.WriteString(string(headbuf))
	buf.WriteByte(byte(PROTO_RPC_ETX))

	rpcreq = buf.Bytes()

	return
}

//将收到的字节切片（包）解析到head和body中
func SrpcUnpackPkgHeadandBody(buff []byte) (head *CRpcHead, body []byte, err error) {

	if buff == nil {
		err = errors.New("SrpcUnpackPkgHead parameters invalid!")
		return
	}

	buff0 := buff[:]
	buflen := len(buff)

	if buflen < 0 {
		plog.Fatal("SrpcUnpackPkgHeadandBody parameter buflen invalid")
		err = errors.New("SrpcUnpackPkgHeadandBody parameter buflen invalid")
		return
	}
	if buflen <= PROTO_HEAD_MIN_LEN || buff[0] != PROTO_RPC_STX || buff[buflen-1] != PROTO_RPC_ETX {
		plog.Fatal("SrpcUnpackPkgHeadandBody format invalid")
		err = errors.New("SrpcUnpackPkgHeadandBody formate invalid")
		return
	}
	headlenslice := buff[PROTO_HEAD_LEN_OFFSET:PROTO_BODY_LEN_OFFSET]
	headlen, err := byteOrderOp.Util_ntoh_int32(headlenslice)
	if err != nil {
		plog.Fatal("SrpcUnpackPkgHeadandBody headlen ntoh ", err)
		return
	}

	//读取包体的长度
	bodylen, err := byteOrderOp.Util_ntoh_int32(buff[PROTO_BODY_LEN_OFFSET:(PROTO_BODY_LEN_OFFSET + 4)])
	if err != nil {
		plog.Fatal("SrpcUnpackPkgHeadandBody bodylen ntoh ", err)
		return
	}

	//校验
	if headlen < 0 || bodylen < 0 || buflen != int(headlen+bodylen+PROTO_HEAD_MIN_LEN) {
		plog.Fatal("SrpcUnpackPkgHeadandBody parameter hl or bl invalid")
		err = errors.New("SrpcUnpackPkgHeadandBody parameter hl or bl invalid")
		return
	}

	head = new(CRpcHead)
	err = head.Unmarshal(buff0[PROTO_HEAD_OFFSET:int(PROTO_HEAD_OFFSET+headlen)])
	if err != nil {
		plog.Fatal("head.Unmarshal err =", err)
	}

	body = buff0[PROTO_HEAD_OFFSET+headlen : (buflen - 1)] //留一个结尾符号
	return
}
