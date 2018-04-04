package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"dfqp/pg-user/entity"
	"putil/log"
)

var UserCacheService = &userCacheService{}

//用户信息key go:userinfo:cid
//用户资料缓存模型类
type userCacheService struct {
}

//获取用户资料缓存HASH的key
//@param int64 mid 用户id
//@return string 缓存hash的key
func (this *userCacheService) getKey(cid int64) string {
	return fmt.Sprintf("go:userinfo:%d", cid)
}
//获取redis配置
func (this *userCacheService) getCfgKey(cid int64) string {
	return fmt.Sprintf("userCache%d", cid % 10)
}
//设置缓存
func (this *userCacheService) Set(cid int64, data interface{}) bool {
	cache, err := NewRedisCache(this.getCfgKey(cid))
	if err == nil {
		ret, err := redis.String(cache.Do("HMSET", redis.Args{}.Add(this.getKey(cid)).AddFlat(data)...))
		plog.Debug("userCacheService set======", ret)
		plog.Debug("userCacheService err======", err)
		if err != nil || ret != "OK" {
			return false
		} else {
			return true
		}
	}
	return false
}
//获取缓存
func (this *userCacheService) Get(cid int64) *entity.User {
	resp := new(entity.User)
	cache, err := NewRedisCache(this.getCfgKey(cid))
	if err == nil {
		data, err := redis.Values(cache.Do("HGETALL", this.getKey(cid)))
		if err != nil {
			return resp
		}
		if err := redis.ScanStruct(data, resp); err != nil {
			return resp
		}
	}
	plog.Debug("userCacheService get======", resp)
	return resp
}
//删除缓存
//@param int64 mid 用户id
//@return bool
func (this *userCacheService) Del(cid int64) bool {
	cache, err := NewRedisCache(this.getCfgKey(cid))
	if err != nil {
		return false
	}
	key := this.getKey(cid)
	res, err := redis.Bool(cache.Do("DEL", key))
	plog.Debug("userCacheService del====", err, res)
	if err != nil && err != redis.ErrNil {
		return false
	}
	return res
}

