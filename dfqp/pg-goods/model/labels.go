package model

import (
	"github.com/astaxie/beego/orm"
	"dfqp/proto/goods"
	"encoding/json"
	"putil/cache"
)

const bizGoodsLabel = "goods-label"

const cacheKeyGoodsLabels = "pg-goods-labels"

const statusAvailable = 1

// GoodsLabel 物品model 实例
var GoodsLabel = &defaultLabelModel
var defaultLabelModel goodsLabelModel

// 物品model
type goodsLabelModel struct {
	biz string
	o   orm.Ormer
}

// 表名
func (m *goodsLabelModel) table(id int64) string {
	return "pg_goods_label"
}

// SetOrm, set a orm for this model
func (m *goodsLabelModel) SetOrm(o orm.Ormer) {
	m.o = o
	m.biz = bizGoodsLabel
	//m.o.Using(bizGoodsLabel)
}

// Add a label
func (m *goodsLabelModel) Add(name string) (bool, error) {

	return true, nil
}

// Update label
func (m *goodsLabelModel) Update(id int64, name string) (bool, error) {
	return true, nil
}

// Delete a label
func (m *goodsLabelModel) Delete(id int64) (bool, error) {

	return true, nil
}


// All get goods labels from cache or db
func (m *goodsLabelModel) All() ([]*ptGoods.BagTab, error) {

	cache := cache.LoadCache("default")
	v, err := cache.Get(cacheKeyGoodsLabels)
	if err != nil {
		println("get goods label from cache err:", err.Error())
	}
	var labels []*ptGoods.BagTab
	if v != "" {
		err := json.Unmarshal([]byte(v), &labels)
		if err == nil {
			println("yes, get cache from redis")
			return labels, nil
		}

		println(labels, "get data from cache failed:", err.Error())
	}
	query := "SELECT label_id, name FROM pg_goods_label WHERE `status`=?"
	num, err := m.o.Raw(query, statusAvailable).QueryRows(&labels)
	if err == nil {
		println("goods label nums:", num)
		bytes, err := json.Marshal(labels)
		if err != nil {
			println("json marshal err:", err.Error())
		}
		v := string(bytes)
		result, err := cache.Save(cacheKeyGoodsLabels, v, 0)
		if err != nil {
			println(result, err.Error())
		}

	}

	return labels, err
}

