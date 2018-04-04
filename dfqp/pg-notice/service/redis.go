package service

import (
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"time"
	"errors"
)

type redisService struct {
	p *redis.Pool
}

/**
 * 初始化redis服务
 */
 func NewRedisCache(alias string) (*redisService, error) {
 	//从配置文件获取redis配置
 	redisHost := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + "redishost")
 	redisPort := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + "redisport")
 	dialFunc := func() (c redis.Conn, err error) {
 		//tcp连接redis
 		c, err = redis.Dial("tcp",redisHost + ":" +redisPort)
 		if err != nil {
 			return nil,err
		}
		return
	}
	//初始化连接池
	p := &redis.Pool{
		MaxIdle:	3,
		IdleTimeout:100 * time.Second,
		Dial:	dialFunc,
	}

	c := p.Get()
	//操作完成后自动关闭
	defer c.Close()
	return &redisService{
		p :p,
	}, c.Err()
 }

 /**
  * 操作redis同一方法
  */
  func (this *redisService) Do(commandName string, args ...interface{}) (reply interface{}, err error){
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	conn := this.p.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
  }