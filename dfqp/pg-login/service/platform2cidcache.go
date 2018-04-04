package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"putil/log"
	"dfqp/lib"
)

var Platform2cidCacheService = &platform2cidCacheService{}

//key md5(go:cidmap:platformId:platformType)
//cid和平台标识映射缓存模型类
type platform2cidCacheService struct {
}

//获取缓存HASH的key
//@param int64 mid 用户id
//@return string 缓存hash的key
func (this *platform2cidCacheService) getKey(platformId string, platformType int32) string {
	md5Str := lib.GetMd5(fmt.Sprintf("go:platform2cid:%s:%d", platformId, platformType))
	plog.Debug("md5Str: ", md5Str, platformId, platformType)
	return md5Str
}
//获取redis配置
func (this *platform2cidCacheService) getCfgKey() string {
	return "cidMap"
}
//设置缓存
func (this *platform2cidCacheService) Set(platformId string, platformType int32, cid int64) bool {
	cache, err := NewRedisCache(this.getCfgKey())
	if err == nil {
		ret, err := redis.String(cache.Do("SET", this.getKey(platformId, platformType), cid))
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
func (this *platform2cidCacheService) Get(platformId string, platformType int32) int64 {
	cache, err := NewRedisCache(this.getCfgKey())
	var data int64
	if err == nil {
		data, err = redis.Int64(cache.Do("GET", this.getKey(platformId, platformType)))
		plog.Debug("get: ", data, err)
		if err != nil {
			return 0
		}
	}
	return data
}
//删除缓存
//@param int64 mid 用户id
//@return bool
func (this *platform2cidCacheService) Del(platformId string, platformType int32) bool {
	cache, err := NewRedisCache(this.getCfgKey())
	if err != nil {
		return false
	}
	key := this.getKey(platformId, platformType)
	res, err := redis.Bool(cache.Do("DEL", key))
	plog.Debug("del: ", err, res)
	if err != nil && err != redis.ErrNil {
		return false
	}
	return res
}