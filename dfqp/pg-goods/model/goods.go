package model

import (
	"dfqp/pg-goods/entity"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

const bizGoods = "goods"

// PGModel interface
type PGModel interface {
	table(id int64) string
	SetOrm(orm.Ormer)

	//add other method
}

// Goods 物品model实例
var Goods = &defaultGoodsModel
var defaultGoodsModel goodsModel

// 物品model
type goodsModel struct {
	biz string
	o   orm.Ormer
}

//表名
func (m *goodsModel) table(id int64) string {
	return "pg_goods"
}

// Create a goods
func (m *goodsModel) Create(item *entity.PgGoods) (int64, error) {
	return m.o.Insert(item)
}

func (m *goodsModel) SetOrm(o orm.Ormer) {
	m.o = o
	m.o.Using(bizGoods)
	m.biz = bizGoods

}

// Get goods list
func (m *goodsModel) Info(goodsID int64) (entity.PgGoods, error) {
	var data entity.PgGoods
	m.o.Using("goods")
	table := m.table(goodsID)
	sql := "SELECT goods_id, `name`, `desc`, goods_type, img, price, price_org, goods_ext, expire_time,`status` FROM "+ table +" WHERE goods_id = ?"

	err := m.o.Raw(sql, goodsID).QueryRow(&data)
	// 没有找到记录
	if err == orm.ErrNoRows {
		return data, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return data, err
	}
	return data, nil
}

func (m *goodsModel) Items(goodsIds []int64) ([]*entity.PgGoods, error) {
	var (
		sql string
		err error
		//data      []*ptGoods.GoodsItem
		goodsdata []*entity.PgGoods
	)
	m.o.Using("goods")
	//分表必须用原生查询
	table := m.table(1)
	goodIn := "("
	for _, goodid := range goodsIds {
		g := strconv.Itoa(int(goodid))
		goodIn = goodIn + g + ","
	}
	goodIn = strings.TrimSuffix(goodIn, ",") + ")"
	sql = "SELECT goods_id, name, `desc`, goods_type, img, goods_ext, label, conditions, app_version FROM " + table + " WHERE status=1 AND goods_id in " + goodIn
	_, err = m.o.Raw(sql).QueryRows(&goodsdata)
	//没有找到记录
	if err == orm.ErrNoRows {
		return goodsdata, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return goodsdata, err
	}

	return goodsdata, nil
}

// GoodsLabel return a goods label
func (m *goodsModel) GoodsLabelID(goodsID int64) (int64, error) {
	var labelID struct {
		Label int64
	}
	m.o.Using("goods")
	table := m.table(goodsID)
	sql := "select label from " + table + " where goods_id=?"
	err := m.o.Raw(sql, goodsID).QueryRow(&labelID)
	if err == nil {
		return labelID.Label, nil
	}
	println(err.Error())
	return labelID.Label, err
}

// Update goods status
func (m *goodsModel) UpdateStatus(goodsID int64, status int32) bool {
	item := &entity.PgGoods{}
	item.GoodsID = goodsID
	item.Status = status
	return m.Update(item, "Status")
}

// Update
func (m *goodsModel) Update(info *entity.PgGoods, cols ...string) bool {
	if info.GoodsID <= 0 {
		return false
	}
	_, err := m.o.Update(info, cols...)
	if err == nil {
		return true
	}
	return false

}

// Delete goods
func (m *goodsModel) Delete(id int64) bool {

	return m.UpdateStatus(id, 2)
}
