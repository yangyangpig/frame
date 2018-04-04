package service

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"dfqp/pg-login/entity"
	"time"
	"putil/log"
)

var RegFailService = &regFailService{}

//用户信息
type regFailService struct {}

//表名
func (this *regFailService) table(cid int64) string {
	return tableName("logs.regfail")
}

//写user表
//@param int64 cid 用户cid
//@param string step 第几步
//@return bool
func (this *regFailService) Add(cid int64, step int32) bool {
	if cid <= 0 {
		return false
	}
	o.Using("logs")
	timeStamp := time.Now().Unix()
	_, err := o.Raw("INSERT INTO " + this.table(cid) + " (cid, step, time) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE step=?", cid, step, timeStamp, step).Exec()
	if err != nil {
		return false
	}
	//删除缓存
	RegFailCacheService.Del(cid)
	return true
}
//根据cid获取用户信息
//@param int64 cid 用户cid
//@return bool
func (this *regFailService) Get(cid int64) (int, error) {
	step, flag := RegFailCacheService.Get(cid)
	if flag {
		plog.Debug("regFailCacheService.Get Resp: ", step)
		return step, nil
	} else {
		var data *entity.RegFail
		o.Using("logs")
		//分表必须用原生查询
		err := o.Raw("SELECT cid, step FROM " + this.table(cid) + " WHERE cid = ?", cid).QueryRow(&data)
		plog.Debug("regFailService.Get Resp: ", data, err)
		//没有找到记录
		if err == orm.ErrNoRows {
			return 0, nil
		}
		if err != nil && err != orm.ErrNoRows {
			return 0, err
		}
		go RegFailCacheService.Set(cid, data.Step)
		return step, nil
	}
}
//更新user表
//@param int64 cid 用户mid
//@param string step 第几步
//@return bool
func (this *regFailService) Update(cid int64, step int32) bool {
	if cid <= 0 {
		return false
	}
	o.Using("logs")
	sql := "UPDATE " + this.table(cid) + " SET step = "+fmt.Sprintf("%d", step)+" WHERE cid=?"
	ret, err := o.Raw(sql, step, cid).Exec()
	if err != nil {
		return false
	}
	num, _ := ret.RowsAffected()
	if num > 0 {
		//删除缓存
		RegFailCacheService.Del(cid)
	}
	return true
}