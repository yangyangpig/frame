package service

import (
	"encoding/json"
	"PGMatch/app/entity"
	"github.com/astaxie/beego/logs"
)

type sessionService struct {
	rd entity.RedisClient
}

const (
	midPrefix       = "ssid_mid_"
	sessionIdPrefix = "ssid_byssid_"
)

// 初始化 session 服务
func NewSessionService() entity.Sessioner {
	rd, err := NewRedisCache(entity.RedisSession)
	if err != nil {
		logs.Error("Redis initialization failed! name[%s]", entity.RedisSession)
	}
	return &sessionService{rd: rd}
}

// 获取用户线上信息
func (s *sessionService) GetOnlineInfo(id interface{}) *entity.OnlineInfo{
	var onlineKey string
	switch t := id.(type) {
	case uint64:
		onlineKey = midPrefix + Uint64ToStr(id.(uint64))
	case string:
		onlineKey = sessionIdPrefix + id.(string)
	default:
		_ = t
	}
	onlineInfo := s.rd.Get(onlineKey)
	if onlineInfo == nil {
		logs.Error("Did not find the user online information! mid[%v]", id)
		panic("OnlineInfo is nil")
	}
	var oi entity.OnlineInfo
	json.Unmarshal(onlineInfo.([]byte), &oi)
	return &oi
}
