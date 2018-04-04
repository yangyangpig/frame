package service

import (
	"dfqp/pg-sign/entity"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type configService struct {}

func (this *configService) table() string {
	return tableName("main.config")
}

func (this *configService) Add(config *entity.Config) bool {
	o.Using("main")
	_, err := o.Raw("INSERT INTO " + this.table() + "(app_id,config_key,config_value,time) VALUES(?, ?, ?, ?)", config.AppId, config.ConfigKey, config.ConfigValue, config.Time).Exec()
	if err != nil {
		logs.Error("configService.Add Fail:", err)
		return false
	}
	return true
}

func (this *configService) Get(app_id int32, config_key string) (entity.Config, error) {
	var data entity.Config
	o.Using("main")
	err := o.Raw("SELECT app_id,config_key,config_value,time FROM " + this.table() + " WHERE app_id = ? AND config_key = ?", app_id, config_key).QueryRow(&data)
	if err == orm.ErrNoRows {
		return entity.Config{}, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return entity.Config{}, err
	}
	return data, nil
}
