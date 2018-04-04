package model

import (
	"dfqp/pg-goods/entity"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

var Prices = &defaultPriceModel
var defaultPriceModel priceModel

// 支付方式model
type priceModel struct {
	biz string
	o   orm.Ormer
}

//表名
func (m *priceModel) table() string {
	return "pg_goods_price"
}

/**
 * 获取物品兑换方式
 */
func (m *priceModel) PriceItems(goodsIds []int64) (map[int64]*entity.GoodsPrice, error) {
	var (
		sql        string
		err        error
		goodsprice []*entity.GoodsPrice
		goodId     int64
		data       = make(map[int64]*entity.GoodsPrice)
	)
	table := m.table()
	goodIn := "("
	for _, goodid := range goodsIds {
		g := strconv.Itoa(int(goodid))
		goodIn = goodIn + g + ","
	}
	goodIn = strings.TrimSuffix(goodIn, ",") + ")"
	sql = "SELECT goods_id, pay,added FROM " + table + " WHERE status=1 AND goods_id in " + goodIn

	_, err = m.o.Raw(sql).QueryRows(&goodsprice)
	//没有找到记录
	if err == orm.ErrNoRows {
		return map[int64]*entity.GoodsPrice{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return map[int64]*entity.GoodsPrice{}, err
	}

	for _, v := range goodsprice {
		goodId = v.GoodsId
		data[goodId] = v
	}
	return data, nil
}
func (m *priceModel) SetOrm(o orm.Ormer) {
	m.o = o
	m.biz = bizGoods
	m.o.Using(bizGoods)
}
