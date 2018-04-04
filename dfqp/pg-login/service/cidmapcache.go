package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"putil/log"
)

var CidMapCacheService = &cidMapCacheService{}
//key go:cidmap:cid:platformType
//cid和平台标识映射缓存模型类
type cidMapCacheService struct {
}

//获取缓存HASH的key
//@param int64 mid 用户id
//@return string 缓存hash的key
func (this *cidMapCacheService) getKey(cid int64, platformType int32) string {
	return fmt.Sprintf("go:cidmap:%d:%d", cid, platformType)
}
//获取redis配置
func (this *cidMapCacheService) getCfgKey() string {
	return "cidMap"
}
//设置缓存
func (this *cidMapCacheService) Set(cid int64, platformType int32, data interface{}) bool {
	cache, err := NewRedisCache(this.getCfgKey())
	if err == nil {
		ret, err := redis.String(cache.Do("SET", this.getKey(cid, platformType), data))
		plog.Debug("set: ", ret, err)
		if err != nil || ret != "OK" {
			return false
		} else {
			return true
		}
	}
	return false
}
//获取缓存
func (this *cidMapCacheService) Get(cid int64, platformType int32) string {
	cache, err := NewRedisCache(this.getCfgKey())
	var data string
	if err == nil {
		data, err = redis.String(cache.Do("GET", this.getKey(cid, platformType)))
		plog.Debug("get: ", data)
		if err != nil {
			return ""
		}
	}
	return data
}
//删除缓存
//@param int64 mid 用户id
//@return bool
func (this *cidMapCacheService) Del(cid int64, platformType int32) bool {
	cache, err := NewRedisCache(this.getCfgKey())
	if err != nil {
		return false
	}
	key := this.getKey(cid, platformType)
	res, err := redis.Bool(cache.Do("DEL", key))
	plog.Debug("del: ", err, res)
	if err != nil && err != redis.ErrNil {
		return false
	}
	return res
}