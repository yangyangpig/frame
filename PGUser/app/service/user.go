package service

import (
	"github.com/astaxie/beego/orm"
	"PGUser/app/entity"
	"PGUser/app/proto"
	"strings"
)

type userService struct {}

//表名
func (this *userService) table(mid int64) string {
	return tableNameByMid("users.users", mid)
}

//写user表
func (this *userService) Add(rq *pgUser.InsertUserRequest) bool {
	o.Using("default")
	_, err := o.Raw("INSERT INTO " + this.table(rq.Mid) + " (mid, app_id) VALUES(?, ?)", rq.Mid, rq.AppId).Exec()
	if err != nil {
		return false
	}
	return true
}

//根据mid获取用户信息
func (this *userService) Get(mid int64) (entity.User, error) {
	var data entity.User
	o.Using("default")
	//分表必须用原生查询
	err := o.Raw("SELECT * FROM " + this.table(mid) + " WHERE mid = ?", mid).QueryRow(&data)
	//没有找到记录
	if err == orm.ErrNoRows {
		return entity.User{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.User{}, err
	}
	return data, nil
}
//更新user表
func (this *userService) Update(rq *pgUser.UpdateUserRequest) bool {
	o.Using("default")
	var upArr []string
	if len(rq.Nick) > 0 {
		upArr = append(upArr, "nick='"+rq.Nick+"'")
	}
	if len(upArr) == 0 {
		return false
	}
	sql := "UPDATE " + this.table(rq.Mid) + " SET " + strings.Join(upArr, ",") + " where mid=?"
	_, err := o.Raw(sql, rq.Mid).Exec()
	if err != nil {
		return false;
	}
	return true
}