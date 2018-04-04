package controller

import (
	"framework/rpcclient/core"
	"dfqp/proto/goods"
	"dfqp/pg-goods/logic"
	"dfqp/pg-goods/model"
	"dfqp/lang"
)

// UserBag 用户背包列表
func UserBag(ctx *rpcclient.Context) {

	println("pgGoods.userBag")
	// 解析 rpc-client 请求参数
	params := ctx.GetRpcRequestBody()
	mid := new(ptGoods.Mid)
	err := mid.Unmarshal(params)
	// 请求参数错误
	if err != nil {
		ctx.Write(GetErrorResponseBytes(4, lang.ZH))
		return
	}
	uID := mid.Id
	if uID <= 0 {
		ctx.Write(GetErrorResponseBytes(2, lang.ZH))
		return
	}

	// 背包标签
	bagTabs, err := model.GoodsLabel.All()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(2022, lang.ZH))
		return
	}

	// 获取最新物品
	newGoods, err := logic.BagLogic.GetNewGoods(mid.GetId())
	if err != nil {
		ctx.Write(GetErrorResponseBytes(2023, lang.ZH))
		return
	}

	// 背包所有物品
	bagItems, err := logic.BagLogic.GetAll(uID)
	if err != nil {
		ctx.Write(GetErrorResponseBytes(2024, lang.ZH))
		return
	}
	// 实例化一个rpc返回对象
	response := new(ptGoods.BagItemsResponse)
	response.Status = 0
	response.Msg = "success"
	bag := new(ptGoods.Bag)
	bag.Tabs = bagTabs
	bag.NewGoods = newGoods
	bag.AllGoods = bagItems
	response.Data = bag
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3, lang.ZH))
		return
	}
	ctx.Write(bytes)

}

// Use user use a goods from bag
func Use(ctx *rpcclient.Context) {
	// 校验参数
	request := new(ptGoods.GoodsUseRequest)
	err := request.Unmarshal(ctx.GetRpcRequestBody())
	if err != nil {
		ctx.Write(GetErrorResponseBytes(4, lang.ZH))
		return
	}
	// 检验用户是否在线
	onLineInfo, err := GetUserOnlineInfo(request.Mid)
	println("user session ID: ", onLineInfo.Ssid)
	/*
	if onLineInfo.Ssid == "" {
		ctx.Write(GetErrorResponseBytes(2002, lang.ZH))
		return
	}
	*/

	leftNum, errCode := logic.BagLogic.Use(request)
	response := new(ptGoods.GoodsUseResponse)
	response.Num = leftNum
	response.Msg = lang.Msg(errCode, lang.ZH)
	response.Status = int64(errCode)
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3,lang.ZH))
		return
	}
	ctx.Write(bytes)

}

// Synthesis 合成碎片
func Synthesis(ctx *rpcclient.Context) {
	request := new(ptGoods.SynthesisRequest)
	err := request.Unmarshal(ctx.GetRpcRequestBody())
	if err != nil {
		ctx.Write(GetErrorResponseBytes(4,lang.ZH))
		return
	}
	// 用户是否登录
	// onLineInfo, err := GetUserOnlineInfo(request.Mid)
	// if onLineInfo.Ssid == "" {
	// 	ctx.Write(GetErrorResponseBytes(2002, lang.ZH))
	// 	return
	// }

	leftNum, errCode := logic.BagLogic.Synthesis(request)
	response := new(ptGoods.GoodsUseResponse)
	response.Status = int64(errCode)
	response.Msg = lang.Msg(errCode,lang.ZH)
	response.Num = leftNum
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3,lang.ZH))
		return
	}
	ctx.Write(bytes)

}

// ExchangeRealGoods 兑换实物
func ExchangeRealGoods(ctx *rpcclient.Context) {
	request := new(ptGoods.ExchangeRealGoodsRequest)
	err := request.Unmarshal(ctx.GetRpcRequestBody())
	if err != nil {
		println("_-----------",err.Error())
		ctx.Write(GetErrorResponseBytes(4,lang.ZH))
		return
	}
	// 用户是否登录
	// onLineInfo, err := GetUserOnlineInfo(request.Mid)
	// if onLineInfo.Ssid == "" {
	// 	ctx.Write(GetErrorResponseBytes(2002, lang.ZH))
	// 	return
	// }

	leftNum, errCode := logic.BagLogic.ExchangeRealGoods(request)
	response := new(ptGoods.GoodsUseResponse)
	response.Status = 0
	response.Msg = lang.Msg(errCode,lang.ZH)
	response.Num = leftNum
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3,lang.ZH))
		return
	}
	ctx.Write(bytes)
}

// ExHistory 兑换历史纪录
func ExHistory(ctx *rpcclient.Context) {
	r := new(ptGoods.ExHistoryRequest)
	err := r.Unmarshal(ctx.GetRpcRequestBody())
	if err != nil {
		ctx.Write(GetErrorResponseBytes(4,lang.ZH))
		return
	}
	// 用户是否登录
	// onLineInfo, err := GetUserOnlineInfo(request.Mid)
	// if onLineInfo.Ssid == "" {
	// 	ctx.Write(GetErrorResponseBytes(2002, lang.ZH))
	// 	return
	// }

	list, err := logic.BagLogic.ExHistory(r.Mid, r.PreIndex,r.New, r.PageSize)
	response := new(ptGoods.ExHistoryResponse)
	response.Status = 0
	response.Msg = lang.Msg(0,lang.ZH)
	response.Data = list
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3,lang.ZH))
		return
	}
	ctx.Write(bytes)
}

// ExchangeTelFee 兑换话费
func ExchangeTelFee(ctx *rpcclient.Context) {
	request := new(ptGoods.ExChangeTelFeeRequest)
	request.Unmarshal(ctx.GetRpcRequestBody())
	// 用户是否登录
	// onLineInfo, err := GetUserOnlineInfo(request.Mid)
	// if onLineInfo.Ssid == "" {
	// 	ctx.Write(GetErrorResponseBytes(2002, lang.ZH))
	// 	return
	// }

	num, code := logic.BagLogic.ExchangeTelFee(request)
	if code != 0 {
		ctx.Write(GetErrorResponseBytes(code,lang.ZH))
		return
	}
	response := new (ptGoods.GoodsUseResponse)
	response.Status = 0
	response.Msg = lang.Msg(0, lang.ZH)
	response.Num = num
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3,lang.ZH))
		return
	}
	ctx.Write(bytes)
}
