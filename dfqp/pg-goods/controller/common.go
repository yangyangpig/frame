package controller

import (
	"framework/rpcclient/core"
	"dfqp/proto/online"
	"dfqp/proto/goods"
	"dfqp/lang"
)

var rpcClient *rpcclient.RpcCall

// SetRPCClient 设置rpc客户端
func SetRPCClient(client *rpcclient.RpcCall) {
	rpcClient  = client
}


// GetErrorResponseBytes generate a rpc error response
func GetErrorResponseBytes(code, langType int)[]byte {
	response := new(ptGoods.RPCErrorResponse)
	response.Status = int32(code)
	response.Msg = lang.Msg(code, langType)
	bytes, _ := response.Marshal()
	return bytes
}

// GetUserOnlineInfo 获取用户在线信息
func GetUserOnlineInfo(mid int64) (*pgOnline.GetOnlineResponse, error) {
	onlineInfo := new(pgOnline.GetOnlineResponse)
	request := new(pgOnline.GetOnlineRequest)
	request.Uid = mid
	reqBytes, err := request.Marshal()
	if err != nil {
		return onlineInfo, err
	}
	resp := rpcClient.SendAndRecvRespRpcMsg("online.OnlineService.Get", reqBytes, 1000, 0)
	if resp.ReturnCode != 0 {
		return onlineInfo, resp.Err
	}
	onlineInfo.Unmarshal(resp.Body)
	return onlineInfo, nil
}



// CheckUserOnLine user online info
func CheckUserOnLine(next rpcclient.RpcHandler ) rpcclient.RpcHandler {
	return rpcclient.RpcHandleFunc(func(ctx *rpcclient.Context) {
		println("用户在线信息")

		bytes := ctx.GetRpcRequestBody()
		uID := new(ptGoods.Mid)
		err := uID.Unmarshal(bytes)
		if err != nil {
			bytes := GetErrorResponseBytes(2, lang.ZH)
			ctx.Write(bytes)
			return
		}
		request := new(pgOnline.GetOnlineRequest)
		request.Uid = uID.Id
		reqBytes, err := request.Marshal()
		if err != nil {
			bytes := GetErrorResponseBytes(2, lang.ZH)
			ctx.Write(bytes)
		}
		resp := rpcClient.SendAndRecvRespRpcMsg("online.OnlineService.Get", reqBytes, 1000, 0)
		if resp.ReturnCode != 0 {
			bytes := GetErrorResponseBytes(2001, lang.ZH)
			ctx.Write(bytes)
		}
		onlineInfo := new(pgOnline.GetOnlineResponse)
		onlineInfo.Unmarshal(resp.Body)
		println("online info:", onlineInfo.Ssid)
		if onlineInfo.Ssid == "" {
			bytes := GetErrorResponseBytes(2002, lang.ZH)
			ctx.Write(bytes)
		}
		next.ServeRPC(ctx)
	})
}