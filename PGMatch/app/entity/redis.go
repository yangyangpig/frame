package entity

import "time"

// Redis 工具接口
type RedisClient interface {
	Get(key string) interface{}
	HGet(key, field string) interface{}
	HGetAll(key string) interface{}
	Hlen(key string) interface{}
	Put(key string, val interface{}, timeout time.Duration) error
	Delete(key string) error
	IsExist(key string) bool
}

// redis address
const (
	RedisMatchConf   = "matchConf"   // 比赛配置
	RedisFreeTime    = "freeTime"    // 免费次数
	RedisMatchInvite = "matchInvite" // 比赛匹配
	RedisSession     = "session"     // session
	RedisMatchUser   = "matchUser"   // 比赛用户
	RedisServerConf  = "serverConf"  // server 配置存储
)
