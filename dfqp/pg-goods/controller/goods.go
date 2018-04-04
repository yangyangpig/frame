package controller

import (
	"dfqp/lang"
	"dfqp/pg-goods/entity"
	"dfqp/pg-goods/logic"
	"dfqp/pg-goods/model"
	"dfqp/proto/goods"
	"framework/rpcclient/core"
	"putil/log"
	"strconv"
)

func init() {
	println("controller/goods.init")
}

// CreateGoods Create a goods
func CreateGoods(ctx *rpcclient.Context) {
	plog.Debug("-----------pgGoods.create ... ----------------")
	body := ctx.GetRpcRequestBody()
	cq := new(ptGoods.CreateGoodsRequest)
	cq.Unmarshal(body)
	var item entity.PgGoods
	item.Name = cq.Name
	item.Desc = cq.Desc
	item.Type = cq.Type
	item.Img = cq.Img
	item.Conditions = cq.Conditions
	item.Price = cq.Price
	item.PriceOrg = cq.PriceOrg
	item.GoodsExt = cq.GoodsExt
	item.ExpireTime = cq.ExpireTime
	status, _ := strconv.Atoi(ptGoods.UNAVAILABLE.String())
	item.Status = int32(status)
	item.CreateTime = cq.CreateTime
	item.CreateBy = cq.CreateBy
	// insert to mysql
	_, err := model.Goods.Create(&item)
	if err != nil {
		plog.Debug("create goods failed:", err.Error())
	}

	ctx.Write([]byte("Hello"))

}

// All Get all goods info
func All(ctx *rpcclient.Context) {
	// body := ctx.GetRpcRequestBody()
	println("pgGoods.all")
}

// GetGoodsInfo Get a goods info
func GetGoodsInfo(ctx *rpcclient.Context) {
	body := ctx.GetRpcRequestBody()
	request := new(ptGoods.GoodsListRequest)
	response := new(ptGoods.GoodsListResponse)
	err := request.Unmarshal(body)
	// 请求参数错误
	if err != nil {
		ctx.Write(GetErrorResponseBytes(4, lang.ZH))
		return
	}

	goodIds := request.GoodsId
	data, code := logic.GoodsLogic.GetItemsBygoodIds(goodIds)
	// 请求失败
	if code != 0 {
		ctx.Write(GetErrorResponseBytes(code, lang.ZH))
		return
	} else {
		response.Goods = data
	}
	response.Status = 0
	response.Msg = lang.Msg(0, lang.ZH)
	bytes, err := response.Marshal()
	if err != nil {
		ctx.Write(GetErrorResponseBytes(3, lang.ZH))
		return
	}
	ctx.Write(bytes)
}

// ExchangeTypeInfo Goods exchange type info
func ExchangeTypeInfo(ctx *rpcclient.Context) {
	println("pgGoods.exchangeTypeInfo")
}
