package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
)

//实例化memcache对象
func NewMemcache(name string)(adapter cache.Cache, err error) {
	host := beego.AppConfig.String(beego.BConfig.RunMode+"::"+name+".memhost")
	port := beego.AppConfig.String(beego.BConfig.RunMode+"::"+name+".memport")
	config := `{"conn":"`+host+":"+port+`"}`
	return cache.NewCache("memcache", config)
}
