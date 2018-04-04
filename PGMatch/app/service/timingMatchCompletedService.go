package service

import (
	"PGMatch/app/entity"
	"github.com/astaxie/beego/logs"
	"strings"
)

// 定时赛完成服务
type timingMatchCompletedService struct {
	TAG string
	rd  entity.RedisClient
}

const (
	//配置id对应的赛制ID_主键 hash
	timingMatchCompletedKey = "FIXED_MATCH_THROUGH_KEY"
	//赛制ID_主键对应多个配置id hash
	timingMatchCompletedRevKey = "FIXED_MATCH_THROUGH_REV_KEY"
)

//　初始化
func NewTimingMatchCompletedService() entity.TimingMatchCompletedInterface {
	tag := "timingMatchCompletedService : "
	rd, err := NewRedisCache(entity.RedisServerConf)
	if err != nil {
		logs.Error(tag + "Redis initialization failed!")
	}
	return &timingMatchCompletedService{TAG: tag, rd: rd}
}

// 获取 configId 赛制组的 configIds 数组
func (ts *timingMatchCompletedService) GetMixMatchConfigId(matchConfigId string) interface{} {
	group := ts.rd.HGet(timingMatchCompletedKey, matchConfigId)
	configIds := ts.rd.HGet(timingMatchCompletedRevKey, group.(string))
	if configIds != nil {
		return strings.Split(ByteToStr(configIds.([]byte)), ",")
	}
	return configIds
}
