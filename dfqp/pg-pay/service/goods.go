package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"putil/log"
)

type goodsService struct{}

func (this *goodsService) GetGoodsConf(tabId, clientId int32) map[string]string {
	redisObj, err := RedisService.NewRedisCache("pay")
	if err != nil {
		plog.Debug("error................")
	}
	goodsKey := this.getGoodsConfKey(tabId, clientId)
	//redis.Do("info")
	result, err := redis.StringMap(redisObj.Do("HGETALL", goodsKey))
	plog.Debug("result:", result)
	return result
}

func (this *goodsService) GetTabConf() string {
	redisObj, err := RedisService.NewRedisCache("pay")
	if err != nil {
		plog.Debug("error................")
	}
	tabKey := this.getTabConfKey()
	result, err := redis.String(redisObj.Do("GET", tabKey))
	return result
}

func (this *goodsService) GetGoodsInfo(tabId, clientId, gid int32) string {
	if tabId <= 0 || clientId <= 0 || gid <= 0 {
		return ""
	}
	redisObj, err := RedisService.NewRedisCache("pay")
	if err != nil {
		plog.Debug("error................")
	}
	goodsKey := this.getGoodsConfKey(tabId, clientId)
	result, err := redis.String(redisObj.Do("HGET", goodsKey, gid))
	return result
}

func (this *goodsService) GetPayCenterConf(appId int32) string {
	if appId <= 0 {
		return ""
	}
	redisObj, err := RedisService.NewRedisCache("pay")
	if err != nil {
		plog.Debug("error................")
	}
	key := fmt.Sprintf("paycenter|conf|%d", appId)
	plog.Debug("key ===>", key)
	result, err := redis.String(redisObj.Do("GET", key))
	return result
}

func (this *goodsService) getGoodsConfKey(tabId, clientId int32) string {
	if tabId <= 0 {
		return ""
	}
	return fmt.Sprintf("pay|goodsconf|%d|%d", tabId, clientId)
}

func (this *goodsService) getTabConfKey() string {
	return "pay|tabconf"
}
