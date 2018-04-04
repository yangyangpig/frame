package logic

import (
	"dfqp/pg-goods/model"
	"dfqp/proto/goods"
	"encoding/json"
	"fmt"
	"putil/cache"
	"strconv"
)

// GoodsLogic 实例
var GoodsLogic = &goodLogicInstance
var goodLogicInstance goodLogic

type goodLogic struct {
}

func (g *goodLogic) GetItemsBygoodIds(goodsIds []int64) ([]*ptGoods.GoodsItem, int) {
	//首先redis里查找
	var gIds []int64
	data := []*ptGoods.GoodsItem{}
	fmt.Println(goodsIds)
	for _, v := range goodsIds {
		itemData, err := g.GetItemsFromRedis(v)
		if itemData == "" || err != nil { //没有查到数据或有误，则查db
			gIds = append(gIds, v)
		} else {
			tmp := new(ptGoods.GoodsItem)
			json.Unmarshal([]byte(itemData), &tmp)
			data = append(data, tmp)
		}
	}
	if len(gIds) == 0 {
		return data, 0
	}
	goodsdata, err := model.Goods.Items(gIds)
	//goodlist没有记录

	if len(goodsdata) == 0 {
		return []*ptGoods.GoodsItem{}, 2014
	}
	//获取失败
	if err != nil {
		return []*ptGoods.GoodsItem{}, 2100
	}
	priceInfo, err := model.Prices.PriceItems(gIds)
	//goodprice没有记录
	if len(priceInfo) == 0 {
		return []*ptGoods.GoodsItem{}, 2014
	}
	//获取失败
	if err != nil {
		return []*ptGoods.GoodsItem{}, 2100
	}
	for _, v := range goodsdata {
		tmp := new(ptGoods.GoodsItem)
		tmp.GoodsId = v.GoodsID
		tmp.Name = v.Name
		tmp.Conditions = v.Conditions
		tmp.Desc = v.Desc
		tmp.Img = v.Img
		tmp.Pay = priceInfo[v.GoodsID].Pay
		tmp.Add = priceInfo[v.GoodsID].Added
		tmp.GoodsExt = v.GoodsExt
		tmp.Label = v.Label
		tmp.GoodsType = v.Type
		tmp.AppVersion = v.Appversion
		data = append(data, tmp)
		g.SetItemsToRedis(tmp)
	}
	return data, 0
}

/**
 * 缓存查找物品信息
 */
func (g *goodLogic) GetItemsFromRedis(gId int64) (string, error) {
	ItemsRedis := cache.LoadCache("bag")
	key := g.GetItemKey(gId)
	data, err := ItemsRedis.Get(key)
	if err != nil {
		return data, err
	}
	return data, nil
}

/**
 * 缓存保存物品信息
 */
func (g *goodLogic) SetItemsToRedis(data *ptGoods.GoodsItem) {
	ItemsRedis := cache.LoadCache("bag")
	goodId := data.GoodsId
	key := g.GetItemKey(goodId)
	if data, err := json.Marshal(data); err == nil {
		ItemsRedis.Save(key, string(data), 0)
	}
}

func (g *goodLogic) GetItemKey(gId int64) string {
	return "goodsitem_" + strconv.Itoa(int(gId))
}
