package service

import (
	"github.com/astaxie/beego/logs"
	"PGMatch/app/entity"
)

type matchInviteService struct {
	rd entity.RedisClient
}

const matchUserPrefix = "INVIT_MATCH_USER_"; //单场邀请赛的报名用户

// 初始化比赛匹配服务
func NewMatchInviteService() entity.MatchInviter {
	rd, err := NewRedisCache(entity.RedisMatchInvite)
	if err != nil {
		logs.Error("Init redis [%v] Fail", "session")
	}
	return &matchInviteService{rd:rd}
}

func (m *matchInviteService) ApplyNum(key string) int {
	num := 0
	matchKey := matchUserPrefix + key
	num = m.rd.Hlen(matchKey).(int)
	return num
}