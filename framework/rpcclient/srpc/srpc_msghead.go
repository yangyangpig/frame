package srpc

import (
	_ "framework/rpcclient/bgf"
	_ "putil/log"
)

/*
//向Net层注册消息返回的状态
type CDispRegMsg struct {
	RetCode int32
}

//编码
func (cdrm *CDispRegMsg) Encode() (body []byte, err error) {
	var msg RPCProto.RegisterDispatchResp
	msg.RetCode = cdrm.RetCode
	body, err = msg.Marshal()
	if err != nil {
		plog.Fatal("CDispRespMsg", err)
	}
	return
}

//解码
func (cdrm *CDispRegMsg) Decode(body []byte) (err error) {
	var msg RPCProto.RegisterDispatchResp
	err = msg.Unmarshal(body)
	if err != nil {
		plog.Fatal("CDispRespMsg", err)
	}
	return
}
*/
