package logic

import (
	"dfqp/proto/goods"
	"dfqp/pg-goods/model"
	"putil/cache"
	"encoding/json"
	"time"
	"dfqp/pg-goods/entity"
)

// BagLogic 背包实例
var BagLogic = &bagLogicInstance
var bagLogicInstance bagLogic

const userBagExpireTime = 7 * 24 * 3600

// 物品道具类型
const propsType uint32 = 2

// 物品礼包类型
const giftType uint32 = 6

type bagLogic struct {
}

func (logic *bagLogic) checkItem(bagItem *entity.PgBagItem) (errCode int) {
	if bagItem.Status == 0 {
		return 2011
	}
	// 检验是否过期
	if bagItem.ExpireTime > 0 && bagItem.ExpireTime > time.Now().Unix() {
		return 2012
	}
	if bagItem.Num <= 0 {
		return 2013
	}
	return 0
}

// Use 用户使用道具或礼包
func (logic *bagLogic) Use(data *ptGoods.GoodsUseRequest) (num int32, errCode int) {
	if data.Mid <= 0 || data.GoodsId <= 0 || data.UgId <= 0 || data.Num <= 0 {
		return 0, 2
	}
	// 检验用户物品是否有效
	bagItem, err := model.Bag.Get(data.Mid, data.UgId)
	if err != nil {
		return bagItem.Num, 2010
	}
	if errCode := logic.checkItem(&bagItem); errCode != 0 {
		return bagItem.Num, errCode
	}
	// 检验数量是否大于零，是否满足使用数量
	if bagItem.Num < data.Num {
		return bagItem.Num, 2013
	}

	// 获取物品信息
	item, err := model.Goods.Info(data.GoodsId)
	if err != nil {
		return 0, 2014
	}
	// 判断物品是否是道具和礼包
	if item.Type != propsType && item.Type != giftType {
		return bagItem.Num, 2015
	}

	if item.Type == propsType {
		// 如果特殊处理需要处理，这里添加道具使用逻辑
		if !model.Bag.UpdateNum(data.Mid, bagItem.UgID, data.Num) {
			return bagItem.Num, 2018
		}
		if left, err := model.Bag.GetItemNum(data.Mid, bagItem.UgID); err == nil {
			return left, 0
		}
		return 0, 2025

	}

	// 使用礼包
	if item.Type == giftType {
		var gifts []*entity.GiftExt
		println(item.GoodsExt)
		err := json.Unmarshal([]byte(item.GoodsExt), &gifts)
		if err != nil {
			return bagItem.Num, 2019
		}
		if !model.Bag.UpdateNum(data.Mid, bagItem.UgID, data.Num) {
			return bagItem.Num, 2018
		}
		for _, gift := range gifts {
			println("礼包物品item：", gift.GoodsID)
			goods, err := model.Goods.Info(gift.GoodsID)
			if err != nil {
				errCode = 2016
				break
			}
			newBagItem := new(entity.PgBagItem)
			newBagItem.Mid = data.Mid
			newBagItem.GoodsID = gift.GoodsID
			newBagItem.GoodsType = goods.Type
			newBagItem.Num = gift.Num
			unixTime := time.Now().Unix()
			if goods.ExpireTime > 0 {
				newBagItem.ExpireTime = unixTime + int64(goods.ExpireTime)
			} else {
				newBagItem.ExpireTime = 0
			}
			newBagItem.Ext = ""
			newBagItem.Status = 1
			newBagItem.AddTime = unixTime
			if _, err := model.Bag.Add(newBagItem); err != nil {
				errCode = 2017
				break
			}
		}
		logic.UpdateBagCache(data.Mid)
	}
	if errCode > 0 {
		return bagItem.Num, errCode
	}
	// TODO 记录流水 物品在背包中使用结果
	// 统计使用物品的情况，把物品的使用详情记录，原始物品在背包中的
	// 的数量,以及成功兑换的物品，兑换失败的物品，全部记录上报
	if left, err := model.Bag.GetItemNum(data.Mid, bagItem.UgID); err == nil {
		return left, 0
	}
	return bagItem.Num, 2025

}

// Synthesis 物品合成
func (logic *bagLogic) Synthesis(data *ptGoods.SynthesisRequest) (leftNum int32, errCode int) {
	// 验证参数
	if data.Mid <= 0 || data.GoodsId <= 0 || data.UgId <= 0 ||
		data.Num <= 0 || data.TargetId == 0 {
		return 0, 2
	}
	// 检验用户物品是否有效
	bagItem, err := model.Bag.Get(data.Mid, data.UgId)
	if err != nil {
		return bagItem.Num, 2010
	}
	if bagItem.Status == 0 {
		return bagItem.Num, 2011
	}
	// 检验是否过期
	if bagItem.ExpireTime > 0 && bagItem.ExpireTime > time.Now().Unix() {
		return bagItem.Num, 2012
	}
	// 检验数量是否大于零，是否满足使用数量
	if bagItem.Num < 0 || bagItem.Num < data.Num {
		return bagItem.Num, 2013
	}
	// 获取物品信息
	item, err := model.Goods.Info(data.GoodsId)
	if err != nil {
		return bagItem.Num, 2014
	}
	// 合成目标是否存在
	var targets []*entity.ComponentExt
	println("碎片扩展字段：", item.GoodsExt)
	if err := json.Unmarshal([]byte(item.GoodsExt), &targets); err != nil {
		println(err.Error())
		return bagItem.Num, 2019
	}
	var has bool
	has = false
	var target *entity.ComponentExt
	for _, v := range targets {
		if data.TargetId == v.Target {
			has = true
			target = v
		}
	}
	if !has {
		return bagItem.Num, 2020
	}
	if bagItem.Num < target.NeedNum {
		return bagItem.Num, 2013
	}
	println("-----------碎片正在合成:", data.Mid, bagItem.UgID, data.Num)
	// 扣除碎片数量
	if !model.Bag.UpdateNum(data.Mid, bagItem.UgID, data.Num) {
		return bagItem.Num, 1
	}
	// 新增合成的物品
	goods, err := model.Goods.Info(target.Target)
	newBagItem := new(entity.PgBagItem)
	newBagItem.Mid = data.Mid
	newBagItem.GoodsID = goods.GoodsID
	newBagItem.GoodsType = goods.Type
	newBagItem.Num = 1
	unixTime := time.Now().Unix()
	if goods.ExpireTime > 0 {
		newBagItem.ExpireTime = unixTime + int64(goods.ExpireTime)
	} else {
		newBagItem.ExpireTime = 0
	}
	newBagItem.Ext = ""
	newBagItem.Status = 1
	newBagItem.AddTime = unixTime
	println("生成新的物品:", newBagItem.String())
	if _, err := model.Bag.Add(newBagItem); err != nil {
		return bagItem.Num, 2017
	}
	logic.UpdateBagCache(data.Mid)
	if left, err := model.Bag.GetItemNum(data.Mid, data.UgId); err == nil {
		return left, 0
	}
	return bagItem.Num, 2025
}

// ExchangeRealGoods 兑换实物
// TODO 记录流水
func (logic *bagLogic) ExchangeRealGoods(data *ptGoods.ExchangeRealGoodsRequest) (int32, int) {
	// 验证参数
	if data.Mid <= 0 || data.GoodsId <= 0 || data.UgId <= 0 || data.RealName == "" ||
		data.Addr == "" || data.Phone == "" {
		return 0, 2
	}
	// TODO 校验 姓名，地址，以及手机号码
	// 检验用户物品是否有效
	bagItem, err := model.Bag.Get(data.Mid, data.UgId)
	if err != nil {
		return bagItem.Num, 2010
	}
	if bagItem.Status == 0 {
		return bagItem.Num, 2011
	}
	// 检验是否过期
	if bagItem.ExpireTime > 0 && bagItem.ExpireTime > time.Now().Unix() {
		return bagItem.Num, 2012
	}
	// 检验数量是否大于零，是否满足使用数量
	if bagItem.Num <= 0 {
		return bagItem.Num, 2013
	}
	// 获取物品信息
	goods, err := model.Goods.Info(data.GoodsId)
	if err != nil {
		return bagItem.Num, 2014
	}
	if !model.Bag.UpdateNum(data.Mid, bagItem.UgID, 1) {
		return bagItem.Num, 2008
	}
	userInfo := new(entity.UserAddrInfo)
	userInfo.UserName = data.RealName
	userInfo.Addr = data.Addr
	userInfo.Phone = data.Phone
	tmp, _:= json.Marshal(userInfo)
	if err != nil {
		println(err.Error())
	}
	item := new(entity.AddUserExchangeHistoryItem)
	item.Mid = data.Mid
	item.GoodsID = data.GoodsId
	item.GoodsName = goods.Name
	item.GoodsType = goods.Type
	item.CreateTime = time.Now().Unix()
	item.Ext = string(tmp)
	// 提交申请
	item.Status = model.SubmitStatus
	if !model.ExHistory.Add(item) {
		return bagItem.Num, 2021
	}
	logic.UpdateBagCache(data.Mid)
	if left, err := model.Bag.GetItemNum(data.Mid, data.UgId); err == nil {
		return left, 0
	}
	return bagItem.Num, 2025

}

// ExchangeRealGoods 兑换话费
// TODO 记录流水
func (logic *bagLogic) ExchangeTelFee(data *ptGoods.ExChangeTelFeeRequest) (int32, int) {
	// 验证参数
	if data.Mid <= 0 || data.GoodsId <= 0 || data.UgId <= 0 || data.Phone == "" {
		return 0, 2
	}
	// TODO 手机号码
	// 检验用户物品是否有效
	bagItem, err := model.Bag.Get(data.Mid, data.UgId)
	if err != nil {
		return bagItem.Num, 2010
	}
	if bagItem.Status == 0 {
		return bagItem.Num, 2011
	}
	// 检验是否过期
	if bagItem.ExpireTime > 0 && bagItem.ExpireTime > time.Now().Unix() {
		return bagItem.Num, 2012
	}
	// 检验数量是否大于零，是否满足使用数量
	if bagItem.Num <= 0 {
		return bagItem.Num, 2013
	}
	// 获取物品信息
	goods, err := model.Goods.Info(data.GoodsId)
	if err != nil {
		return bagItem.Num, 2014
	}
	hItem := bagItem.Num
	if !model.Bag.UpdateNum(data.Mid, bagItem.UgID, 1) {
		return hItem, 2008
	}
	item := new(entity.AddUserExchangeHistoryItem)
	item.Mid = data.Mid
	item.GoodsID = data.GoodsId
	item.GoodsName = goods.Name
	item.GoodsType = goods.Type
	item.CreateTime = time.Now().Unix()
	item.Ext = "{\"phone\":\"" + data.Phone + "\"}"
	// 提交申请
	item.Status = model.SubmitStatus
	if !model.ExHistory.Add(item) {
		return hItem, 2021
	}
	logic.UpdateBagCache(data.Mid)
	if left, err := model.Bag.GetItemNum(data.Mid, data.UgId); err == nil {
		return left, 0
	}
	return bagItem.Num, 0
}

// ExHistory 兑换历史记录
func (logic *bagLogic) ExHistory(mid, preIndex int64, isNew, size uint32) ([]*ptGoods.ExChangeItem, error) {
	items, err := model.ExHistory.Get(mid, preIndex, isNew, size)
	var result []*ptGoods.ExChangeItem
	if err != nil {
		return result, err
	}
	result = make([]*ptGoods.ExChangeItem, len(items))
	for k, v := range items {
		item := new(ptGoods.ExChangeItem)
		item.GoodsName = v.GoodsName
		item.Ext = v.Ext
		item.Status = logic.getExItemStatusString(v.GoodsType, v.Status)
		item.CreateTime = v.CreateTime
		result[k] = item
	}
	return result, nil
}

// getExItemStatusString 通过物品类型和兑换记录状态返回状态描述
// TODO 实现
func (logic *bagLogic) getExItemStatusString(goodsType uint32, status uint32) string {
	if goodsType == propsType {

	}
	return "测试"
}

// Add 向背包添加一个
func (logic *bagLogic) Add(mid, goodsID int64, num uint) bool {
	// 取物品分信息析物品类型
	// 插入数据库
	// model.Bag.Add()
	// update cache
	logic.UpdateBagCache(mid)
	// 记录流水
	return true
}

// Delete a bag item
func (logic *bagLogic) Delete(ugID, mid, goodsID int64) bool {
	// 从数据库中删除数据
	// update cache
	logic.UpdateBagCache(mid)
	// 记录操作流水
	return true
}


// UpdateBagCache 更新用户背包缓存
func (logic *bagLogic) UpdateBagCache(mid int64) {
	items, err := model.Bag.Items(mid)
	if err != nil || len(items) <= 0 {
		return
	}
	// 获取标签id TODO 优化 一次性取多个标签id
	for k, item := range items {
		labelId, err := model.Goods.GoodsLabelID(item.GoodsId)
		if err == nil {
			items[k].TabId = labelId
		}
	}
	if bytes, err := json.Marshal(items); err == nil {
		value := string(bytes)
		bagCache := cache.LoadCache("bag")
		key := entity.UserBagCacheKey(mid)
		bagCache.Save(key, value, userBagExpireTime)
	}
}

// ClearCache 清除背包缓存
func (logic *bagLogic) ClearCache(mid int64) {

}

// GetAll 获取背包所有的物品
func (logic *bagLogic) GetAll(mid int64) ([]*ptGoods.BagItem, error) {
	var items []*ptGoods.BagItem
	var err error
	// 从缓存中读取
	bagCache := cache.LoadCache("bag")
	key := entity.UserBagCacheKey(mid)
	if v, err := bagCache.Get(key); err == nil {
		if err := json.Unmarshal([]byte(v), &items); err == nil {
			println("get bag from cache")
			return items, nil
		}
	}
	// 从db中读用户所有的物品
	items, err = model.Bag.Items(mid)
	if err != nil || len(items) <= 0 {
		return items, err
	}

	// 获取标签id
	// TODO 优化 一次性取多个标签id
	for k, item := range items {
		labelId, err := model.Goods.GoodsLabelID(item.GoodsId)
		if err == nil {
			items[k].TabId = labelId
		}
	}
	// cache
	if bytes, err := json.Marshal(items); err == nil {
		value := string(bytes)
		bagCache.Save(key, value, userBagExpireTime)
	}

	return items, nil
}

// Tabs 获取背包标签
func (logic *bagLogic) Tabs() ([]*ptGoods.BagTab, error) {
	return model.GoodsLabel.All()
}

// GetNewGoods 获取最新的物品列表
func (logic *bagLogic) GetNewGoods(mid int64) ([]*ptGoods.BagItem, error) {
	var list []*ptGoods.BagItem
	pgCache := cache.LoadCache("bag")
	key := entity.UserNewGoodsCacheKey(mid)
	value, err := pgCache.Get(key)
	if err != nil {
		return list, nil
	}
	if err := json.Unmarshal([]byte(value), &list); err != nil {
		return list, err
	}
	return list, nil
}
