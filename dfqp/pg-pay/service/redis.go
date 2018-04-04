package service

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisService struct {
	p *redis.Pool
}

/**
 * 初始化redis服务
 */
func (this *redisService) NewRedisCache(alias string) (*redisService, error) {
	redisConf := ServiceConf["redis"]
	redisHost := redisConf.String(Runmode + "::redis." + alias + ".host")
	redisPort := redisConf.String(Runmode + "::redis." + alias + ".port")
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", redisHost+":"+redisPort)
		if err != nil {
			return nil, err
		}
		return
	}
	// 初始化连接池
	p := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
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
func (this *redisService) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	conn := this.p.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
