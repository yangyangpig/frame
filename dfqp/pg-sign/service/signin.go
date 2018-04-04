package service

import (
	"dfqp/pg-sign/entity"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type signinService struct {}

func (this *signinService) table() string {
	return tableName("main.signin")
}

func (this *signinService) Add(mid int64, month int32, sign_value string) bool {
	o.Using("main")
	_, err := o.Raw("INSERT INTO " + this.table() + "(mid,month,sign_value) VALUES(?, ?, ?)", mid, month, sign_value).Exec()
	if err != nil {
		logs.Error("signinService.Add Fail:", err)
		return false
	}
	return true
}

func (this *signinService) Get(mid int64, month int32) (entity.Signin, error) {
	var data entity.Signin
	o.Using("main")
	err := o.Raw("SELECT mid,month,sign_value FROM " + this.table() + " WHERE mid = ? AND month = ?", mid, month).QueryRow(&data)
	if err == orm.ErrNoRows {
		return entity.Signin{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.Signin{}, err
	}
	return data, nil
}
