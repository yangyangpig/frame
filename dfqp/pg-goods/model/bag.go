package model

import (
	"github.com/astaxie/beego/orm"
	"dfqp/pg-goods/entity"
	"dfqp/proto/goods"
	"strconv"
	"time"
)

var bizBag = "bag"

// 物品model
type bagModel struct {
	biz string
	o   orm.Ormer
}

// 表名
func (m *bagModel) table(id int64) string {
	return "pg_bag"
}

// GetItem 获取物品
func (m *bagModel) GetItem(mID, goodsID int64) (*entity.PgBagItem, error) {
	var item *entity.PgBagItem
	table := m.table(mID)
	sql := "select ug_id, num, expire_time, ext, status from " + table + " where mid=? AND goods_id=? AND status in (0, 1)"
	err := m.o.Raw(sql, mID, goodsID).QueryRow(&item)
	return item, err
}


// CheckItem 从背包中取一个无过期时间的物品
func (m *bagModel) GetExistItem(mID, goodsID int64 ) (*entity.PgBagItem, error){
	var item *entity.PgBagItem
	table := m.table(mID)
	sql := "select ug_id, num, ext, status from " + table + " where mid=? AND goods_id=? AND expire_time=0 AND status in (0, 1) ORDER BY ug_id DESC limit 1"
	err := m.o.Raw(sql, mID, goodsID).QueryRow(&item)
	return item, err
}


// Add  create goods to bag
// TODO 加锁
func (m *bagModel) Add(item *entity.PgBagItem) (int64, error) {
	if item.Num <= 0 {
		return 0, nil
	}
	table := m.table(item.Mid)
	if item.ExpireTime == 0 {
		hasItem, err := m.GetExistItem(item.Mid, item.GoodsID)
		// 背包中存在物品
		// 没有过期时间的物品可以叠加
		if err == nil {
			setString := ""
			setString += " SET num=num+" + strconv.Itoa(int(item.Num))
			setString += " ,add_time=" + strconv.FormatInt(item.AddTime, 10)

			whereString := " WHERE mid=? AND ug_id=?"
			if hasItem.Status == 0 {
				whereString += " `status`=1 "
			}
			sql := "UPDATE " + table + setString + whereString
			result, err := m.o.Raw(sql, item.Mid, hasItem.UgID).Exec()
			if err == nil {
				return result.RowsAffected()
			}
			return 0, err

		}
	}


	midString := strconv.FormatInt(item.Mid, 10)
	vString := "(" + midString + "," + strconv.FormatInt(item.GoodsID, 10) + "," +
		strconv.Itoa(int(item.Num)) + "," +strconv.Itoa(int(item.GoodsType))+ "," +
		strconv.FormatInt(item.ExpireTime, 10) + ",\"" + item.Ext + "\", 1, " +
		strconv.FormatInt(time.Now().Unix(), 10) + ")"
	sql := "INSERT INTO " + table + "(`mid`, `goods_id`, `num`,`goods_type`,`expire_time`, `ext`, `status`, `add_time`) VALUES " + vString
	println("add goods to bag :", sql)
	result, err := m.o.Raw(sql).Exec()
	if err != nil {
		println(err.Error())
	}
	return result.RowsAffected()

}

// Get goods list from user bag
func (m *bagModel) Get(mid, ugID int64) (entity.PgBagItem, error) {
	table := m.table(mid)
	var data entity.PgBagItem
	status := ptGoods.AVAILABLE
	sql := "SELECT ug_id, goods_id, goods_type, num, expire_time, ext, status FROM " + table +
		" WHERE mid=? AND ug_id=? AND `status`=?"
	err := m.o.Raw(sql, mid, ugID, status).QueryRow(&data)
	// 没有找到记录
	if err == orm.ErrNoRows {
		return entity.PgBagItem{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.PgBagItem{}, err
	}
	return data, nil
}

// Items return user's bag
func (m *bagModel) Items(mid int64) ([]*ptGoods.BagItem, error) {
	var items []*ptGoods.BagItem
	table := m.table(mid)
	sql := "SELECT ug_id, goods_id, num, expire_time, ext FROM " + table + " WHERE mid=? and `status`=?"
	_, err := m.o.Raw(sql, 1, mid).QueryRows(&items)
	if err == nil {
		// cache

		return items, nil
	}
	return items, err

}

// UseProps 使用道具
func (m *bagModel) UseProps(item *entity.PgBagItem) {
	if item.Status == 0 {

	}
}

// Update update a item
// TODO 优化在redis中加锁
func (m *bagModel) UpdateNum(mID, ugID int64, useNum int32) bool {

	table := m.table(mID)
	sql := "UPDATE " + table + " SET num=num-" + strconv.Itoa(int(useNum)) + " WHERE ug_id=? AND mid=?"
	result, err := m.o.Raw(sql, ugID, mID).Exec()
	if err != nil {
		return false
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return affected > 0

}

// GetItemNum 获取物品数量
func (m *bagModel) GetItemNum(mID, ugID int64) (int32, error) {
	table := m.table(mID)
	query := "SELECT num FROM " + table + " WHERE ug_id=? AND mid=?"
	var num struct {
		Num int32
	}
	err := m.o.Raw(query, ugID, mID).QueryRow(&num)
	if err == nil {
		return num.Num, nil
	}
	return 0, err
}

// Delete goods
func (m *bagModel) Delete(id int64) bool {

	return true
}

// SetOrm set orm
func (m *bagModel) SetOrm(o orm.Ormer) {
	m.o = o
	m.biz = bizBag
	m.o.Using(bizBag)
}

// Bag model bag instance
var Bag = &defaultBagModel
var defaultBagModel bagModel
