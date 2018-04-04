package service

import (
	"github.com/astaxie/beego/orm"
	"dfqp/pg-user/entity"
	"strings"
	"putil/log"
)

var UserService = &userService{}

//用户信息
type userService struct {}

//表名
func (this *userService) table(cid int64) string {
	return tableNameBy("users.users", cid, 2)
}

//写user表
//@param int64 cid 用户mid
//@param string nick 昵称
//@param int32 sex 性别
//@param string city 城市
//@param string phone 手机
//@param string sign 签名
//@return bool
func (this *userService) Add(cid int64, sex int32, nick, city, phone, sign, icon, icon_big string) bool {
	if cid <= 0 {
		return false
	}
	o.Using("default")
	_, err := o.Raw("INSERT INTO " + this.table(cid) + " (cid,nick,sex,icon,icon_big,city,phone,sign) VALUES(?, ?, ?, ?, ?, ?, ?, ?)", cid, nick, sex, icon,icon_big, city, phone, sign).Exec()
	if err != nil {
		return false
	}
	return true
}
//根据cid获取用户信息
//@param int64 cid 用户cid
//@return bool
func (this *userService) Get(cid int64) (*entity.User, error) {
	cacheData := UserCacheService.Get(cid)
	plog.Debug("userService.get cache data", cacheData)
	if cacheData.Cid > 0 {
		return cacheData, nil
	} else {
		var data *entity.User
		o.Using("default")
		//分表必须用原生查询
		err := o.Raw("SELECT cid, sex, nick, city, phone, sign, icon, icon_big, icon_id FROM " + this.table(cid) + " WHERE cid = ?", cid).QueryRow(&data)
		//没有找到记录
		if err == orm.ErrNoRows {
			return nil, nil
		}
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}
		UserCacheService.Set(cid, data)
		return data, nil
	}
}
//更新user表
//@param int64 cid 用户mid
//@param string nick 昵称
//@param int32 sex 性别
//@param string city 城市
//@param string phone 手机
//@param string sign 签名
//@return bool
func (this *userService) Update(cid int64, sex int32, nick, city, phone, sign, icon, iconBig string, iconId string) bool {
	if cid <= 0 {
		return false
	}
	o.Using("default")
	var data []string
	var values []interface{}
	if len(nick) > 0 {
		data = append(data, "nick=?")
		values = append(values, nick)
	}
	if len(city) > 0 {
		data = append(data, "city=?")
		values = append(values, city)
	}
	if len(phone) > 0 {
		data = append(data, "phone=?")
		values = append(values, phone)
	}
	if len(sign) > 0 {
		data = append(data, "sign=?")
		values = append(values, sign)
	}
	if len(icon) > 0 {
		data = append(data, "icon=?")
		values = append(values, icon)
	}
	if len(iconBig) > 0 {
		data = append(data, "icon_big=?")
		values = append(values, iconBig)
	}
	if sex > 0 {
		data = append(data, "sex=?")
		values = append(values, sex)
	}
	if len(iconId) > 0 {
		data = append(data, "icon_id=?")
		values = append(values, iconId)
	}
	if len(data) == 0 {
		return false
	}
	sql := "UPDATE " + this.table(cid) + " SET " + strings.Join(data, ",") + " WHERE cid=?"
	ret, err := o.Raw(sql, values, cid).Exec()
	plog.Debug("userService.update err", err)
	if err != nil {
		return false
	}
	num, _ := ret.RowsAffected()
	if num > 0 {
		//删除缓存
		UserCacheService.Del(cid)
	}
	return true
}