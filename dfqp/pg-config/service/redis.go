package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"errors"
	"strings"
)

type redisService struct {
	p *redis.Pool
}

var (
	configErr = errors.New("配置错误")
)

//alias
func NewRedisCache(alias string) (*redisService, error) {
	cfg := ServiceConf["redis"].String(Runmode+"::redis."+alias+".host")
	cfgSlice := strings.Split(cfg, ":")
	if len(cfgSlice) != 2 {
		return nil, configErr
	}
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", cfgSlice[0]+":"+cfgSlice[1])
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
