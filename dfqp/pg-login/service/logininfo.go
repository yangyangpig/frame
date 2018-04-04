package service

import (
	"time"
	"dfqp/pg-login/entity"
	"github.com/astaxie/beego/orm"
)

var LoginInfoService = &loginInfoService{}
//登录信息日志
type loginInfoService struct {}

func (this *loginInfoService) table(cid int64) string {
	return tableNameBy("logs.logininfo", cid, 2)
}

//插入
func (this *loginInfoService) Insert(cid int64, firstApp, lastApp int32, firstVersion, lastVersion, firstIp, lastIp string) bool {
	if cid <= 0 {
		return false
	}
	o.Using("logs")
	timeStamp := time.Now().Unix()
	_, err := o.Raw("INSERT INTO " + this.table(cid) + " (cid,first_app,last_app,first_version,last_version,reg_time,login_time,first_ip,last_ip) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", cid, firstApp, lastApp, firstVersion, lastVersion, timeStamp, timeStamp, firstIp, lastIp).Exec()
	if err != nil {
		return false
	}
	return true
}

//更新
func (this *loginInfoService) Update(cid int64, lastApp int32, lastVersion string, lastIp string) bool {
	if cid <= 0 || lastApp <= 0 || len(lastVersion) == 0 || len(lastIp) == 0 {
		return false
	}
	o.Using("logs")
	timeStamp := time.Now().Unix()
	_, err := o.Raw("UPDATE " + this.table(cid) + " SET last_app=?, last_version=?, last_ip=?, login_time=? WHERE cid=?", lastApp, lastVersion, lastIp, timeStamp, cid).Exec()
	if err != nil {
		return false
	}
	return true
}

func (this *loginInfoService) GetInfoByCid(cid int64) (entity.Logininfo, error) {
	var data entity.Logininfo

	if cid <= 0 {
		return entity.Logininfo{}, nil
	}

	o.Using("logs")
	err := o.Raw("SELECT cid, first_app, last_app, first_version, last_version, reg_time, login_time, first_ip, last_ip FROM " + this.table(cid) + " WHERE cid = ?", cid).QueryRow(&data)
	// 没有找到记录
	if err == orm.ErrNoRows {
		return entity.Logininfo{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.Logininfo{}, err
	}
	return data, nil
}