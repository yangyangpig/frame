package service

import (
	"github.com/astaxie/beego/logs"
	"PGMatch/app/entity"
)

// 后台配置服务
type matchConfService struct {
	rd entity.RedisClient
}

// 初始化比赛配置服务
func NewMatchConfService() entity.MatchConfer {
	rd, err := NewRedisCache(entity.RedisMatchConf)
	if err != nil {
		logs.Error("Init redis [%v] Fail", "matchConf")
	}
	return &matchConfService{rd: rd}
}

// 获取指定配置段内容
func (ms *matchConfService) GetAllConfBy(key string) (resp interface{}) {
	resp = ms.rd.HGetAll(key)
	return
}

// 获取配置项的原值
func (ms *matchConfService) GetConfBy(key, filed string) (resp interface{}) {
	resp = ms.rd.HGet(key, filed)
	return
}