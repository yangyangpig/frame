package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"errors"
	"PGMatch/app/entity"
	"github.com/astaxie/beego/logs"
)

type redisService struct {
	p *redis.Pool
}

/**
 * 初始化redis服务
 */
func NewRedisCache(alias string) (entity.RedisClient, error) {
	if RedisConf == nil {
		panic("RedisConf is not configured!")
	}
	redisAddr, ok := RedisConf[alias]
	if !ok {
		panic("RedisConf is not have " + alias + " !")
	}
	redisHost, ok := redisAddr["host"]
	if !ok {
		panic("RedisConf [" + alias + "] is not have host!")
	}
	redisPort, ok := redisAddr["port"]
	if !ok {
		panic("RedisConf [" + alias + "] is not have port!")
	}
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", redisHost+":"+redisPort)
		if err != nil {
			return nil, err
		}
		logs.Debug("Redis[%v] init success!", alias)
		return
	}
	// 初始化连接池
	p := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
		Wait:        true,
	}

	c := p.Get()
	defer c.Close()
	return &redisService{
		p: p,
	}, c.Err()
}

/**
 * 操作redis统一方法
 */
func (rs *redisService) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	conn := rs.p.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

// Get cache from redis.
func (rs *redisService) Get(key string) interface{} {
	if v, err := rs.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetConfBy cache from redis
func (rs *redisService) HGet(key, field string) interface{} {
	if v, err := rs.do("HGET", key, field); err == nil {
		return v
	}
	return nil
}

// GetAllConfBy cache from redis
func (rs *redisService) HGetAll(key string) interface{} {
	values, err := redis.Values(rs.do("HGETALL", key))
	if err != nil {
		return nil
	}
	data := map[string]interface{}{}
	for i := 0; i < len(values); i += 2 {
		filed, _ := redis.String(values[i], nil)
		value := values[i+1]
		data[filed] = value
	}
	return data
}

func (rs *redisService) Hlen(key string) interface{} {
	if v, err := rs.do("HLEN", key); err == nil {
		return v
	}
	return nil
}

//// GetMulti get cache from redis.
//func (rs *redisService) GetMulti(keys []string) []interface{} {
//	c := rs.p.Get()
//	defer c.Close()
//	var args []interface{}
//	for _, key := range keys {
//		args = append(args, rs.associate(key))
//	}
//	values, err := redis.Values(c.Do("MGET", args...))
//	if err != nil {
//		return nil
//	}
//	return values
//}

// Put put cache to redis.
func (rs *redisService) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rs.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rs *redisService) Delete(key string) error {
	_, err := rs.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rs *redisService) IsExist(key string) bool {
	v, err := redis.Bool(rs.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rs *redisService) Incr(key string) error {
	_, err := redis.Bool(rs.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rs *redisService) Decr(key string) error {
	_, err := redis.Bool(rs.do("INCRBY", key, -1))
	return err
}

//// ClearAll clean all cache in redis. delete this redis collection.
//func (rs *redisService) ClearAll() error {
//	c := rs.p.Get()
//	defer c.Close()
//	cachedKeys, err := redis.Strings(c.Do("KEYS", rs.key+":*"))
//	if err != nil {
//		return err
//	}
//	for _, str := range cachedKeys {
//		if _, err = c.Do("DEL", str); err != nil {
//			return err
//		}
//	}
//	return err
//}
