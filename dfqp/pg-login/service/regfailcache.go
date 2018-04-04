package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"putil/log"
)

var RegFailCacheService = &regFailCacheServide{}

//用户注册失败处记录表
type regFailCacheServide struct {
}

//获取用户资料缓存HASH的key
//@param int64 mid 用户id
//@return string 缓存hash的key
func (this *regFailCacheServide) getKey(cid int64) string {
	return fmt.Sprintf("go:regfail:%d", cid)
}
//设置缓存
func (this *regFailCacheServide) Set(cid int64, step int32) bool {
	cache, err := NewRedisCache("userLogin")
	if err == nil {
		ret, err := redis.String(cache.Do("SET", this.getKey(cid), step))
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
func (this *regFailCacheServide) Get(cid int64) (int, bool) {
	cache, err := NewRedisCache("userLogin")
	var (
		step int
	)
	if err == nil {
		step, err = redis.Int(cache.Do("GET", this.getKey(cid)))
		plog.Debug("get: ", step, err)
		if err != nil {
			return 0, false
		}
	} else {
		return 0, false
	}
	return step, true
}
//删除缓存
//@param int64 mid 用户id
//@return bool
func (this *regFailCacheServide) Del(cid int64) bool {
	cache, err := NewRedisCache("userLogin")
	if err != nil {
		return false
	}
	key := this.getKey(cid)
	res, err := redis.Bool(cache.Do("DEL", key))
	plog.Debug("del: ", err, res)
	if err != nil && err != redis.ErrNil {
		return false
	}
	return res
}

