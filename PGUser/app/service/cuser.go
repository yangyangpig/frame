package service

import (
	"github.com/astaxie/beego/orm"
	"PGUser/app/entity"
)

type cuserService struct {}

//表名
func (this *cuserService) table() string {
	return tableName("users.users")
}

//写user表
func (this *cuserService) Insert() {

}

//根据cid获取用户公共信息
func (this *cuserService) get(cid int64) (entity.Cuser, error) {
	var data entity.Cuser
	o.Using("common")
	err := o.QueryTable(this.table()).Filter("cid", cid).One(&data)
	//没有找到记录
	if err == orm.ErrNoRows {
		return entity.Cuser{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.Cuser{}, err
	}
	return data, nil
}