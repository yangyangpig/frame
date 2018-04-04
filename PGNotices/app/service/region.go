package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
)

type regionService struct {
}

//地区相关缓存key
const (
	REGION_KEY = "GO_region"
	PROVINCE_KEY = "GO_province"
)

/**
 * 获取城市相关配置
 */
func (this *regionService) GetRegion() string {
	redisObj, err := NewRedisCache("deploy")
	if err != nil {
		beego.Error("[regionService] new deploy fail:", err.Error())
		return ""
	}
	v, err := redis.String(redisObj.Do("GET", REGION_KEY))
	if err != nil {
		beego.Error("[regionService] GetRegion fail:", err.Error())
		return ""
	}
	return v
}

/**
 * 获取城市和省份相关数据
 */
func (this *regionService) GetProvince() string {
	redisObj, err := NewRedisCache("deploy")
	if err != nil {
		beego.Error("[regionService] new deploy fail:", err.Error())
		return ""
	}
	v, err := redis.String(redisObj.Do("GET", PROVINCE_KEY))
	if err != nil {
		beego.Error("[regionService] GetProvince fail:", err.Error())
		return ""
	}
	return v
}
