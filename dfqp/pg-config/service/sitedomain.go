package service

import (
	"github.com/garyburd/redigo/redis"
)

var SiteDomainService = &siteDomainService{}

type siteDomainService struct {}

//获取access域名
func (this *siteDomainService) Get() string {
	redisObj, err := NewRedisCache("sitetake0")
	if err != nil {
		return ""
	}
	ret, err := redis.String(redisObj.Do("GET", this.getKey()))
	if err != nil {
		return ""
	}
	return ret
}
//获取key
func (this *siteDomainService) getKey() string {
	return "PG:GLOBAL:HALL"
}


