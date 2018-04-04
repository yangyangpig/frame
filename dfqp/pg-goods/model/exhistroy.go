package model

import (
	"github.com/astaxie/beego/orm"
	"dfqp/pg-goods/entity"
	"fmt"
)
// ExHistory 兑换记录model实例
var ExHistory = &defaultExHistoryModel
var defaultExHistoryModel exHistory

// SubmitStatus 提交申请状态
const SubmitStatus = 3

// 物品model
type exHistory struct {
	biz string
	o   orm.Ormer
	env string
}

// SetEvn 设置环境
func (m *exHistory) SetEvn(env string) {
	m.env = env
}

// 表名分表分库100个表
func (m *exHistory) table(id int64) string {
	table := "pg_user_exchange_history"
	println("----env-----", m.env)
	if m.env == entity.EnvDev {
		return table
	}
	return fmt.Sprintf(table+"_%d", id%100)
}

// Add 添加兑换历史纪录
func (m *exHistory) Add (item *entity.AddUserExchangeHistoryItem) bool {
	table := m.table(item.Mid)
	sql := "INSERT INTO " + table + "(mid, goods_id, goods_name, goods_type, ext, create_time, status) VALUES" +
		"(?,?,?,?,?,?,?)"
	_, err := m.o.Raw(sql, item.Mid, item.GoodsID, item.GoodsName, item.GoodsType, item.Ext, item.CreateTime, item.Status).Exec()
	if err != nil {
		println(err.Error())
		return false
	}
	return true
}

// Get 用户兑换历史纪录å
func (m *exHistory) Get(mid, preIndex int64, isNew, size uint32) ([]*entity.AddUserExchangeHistoryItem, error) {
	var data []*entity.AddUserExchangeHistoryItem
	table := m.table(mid)
	if size == 0 {
		size = 30
	}
	whereString := " WHERE mid=? AND "
	if isNew > 0 { // 最新
		whereString += "ueh_id > ? "
	} else {
		whereString += "ueh_id < ? "
	}
	whereString += " ORDER BY ueh_id DESC LIMIT ?"
	sql := "SELECT ueh_id, mid, goods_id, goods_name, goods_type, ext, status FROM " + table + whereString
	_, err := m.o.Raw(sql, mid, preIndex, size).QueryRows(&data)
	// 没有找到记录
	return data, err
}

func (m *exHistory) SetOrm(o orm.Ormer) {
	m.o = o
	m.biz = bizBag
	m.o.Using(bizBag)
}
