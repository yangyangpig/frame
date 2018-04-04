// 每日免费次数服务

package service

import (
	"strings"
	"math"
	"github.com/astaxie/beego/logs"
	"PGMatch/app/entity"
)

type freeTimeService struct {
	rd entity.RedisClient
}

const (
	// 每日免费次数 	Redis : matchConfigIdKey(比赛配置 id) => allTimes(所有免费次数)
	freeTimesKey = "DAILY_FREE_TIMES_KEY"
	// 所有比赛配置 	Redis : allTimes(所有免费次数) => matchConfigIds(所有比赛配置 id "111000, 222000...")
	freeTimesMatchConfigKey = "DAILY_FREE_TIMES_REV_KEY"
	// 用户已使用次数 	Redis :
	freeTimesUsedKey = "DAILY_FREE_USED_TIMES_KEY"
)

// 初始化每日免费次数服务
func NewFreeTimeService() entity.FreeTimer {
	rd, err := NewRedisCache(entity.RedisFreeTime)
	if err != nil {
		logs.Error("Init redis [%v] Fail", "freeTime")
	}
	return &freeTimeService{rd: rd}
}

// 获取 比赛配置ID 对应的每日免费次数
func (f *freeTimeService) get(matchConfigId int) string {
	freeTimes := f.rd.HGet(freeTimesKey, IntToStr(matchConfigId))
	if v, ok := freeTimes.(string); ok {
		return v
	}
	return "0"
}

// 获取 比赛配置ID 的免费次数
func (f *freeTimeService) GetTotalTimes(macthConfigId int) int {
	freeTimes := f.get(macthConfigId)
	totalTimes := 0
	if freeTimes != "0" {
		totalTimes = StrToInt(strings.Split(freeTimes, "_")[0])
	}
	return totalTimes
}

// 获取已经使用的免费次数
func (f *freeTimeService) GetUsedTimes(matchConfigId, mid int) int {
	freeTimes := StrToInt(f.get(matchConfigId))
	var usedTimes int
	if freeTimes > 0 {
		matchConfRedis, _ := NewRedisCache(entity.RedisMatchConf)
		matchConfig := matchConfRedis.HGet(freeTimesMatchConfigKey, IntToStr(freeTimes))
		matchConfigId := matchConfig.(string)
		matchConfigIdArr := strings.Split(matchConfigId, ",")
		userMatchConfigIdArr := f.rd.HGetAll(IntToStr(mid)).(map[string]interface{})
		for _, matchConfigIds := range matchConfigIdArr {
			if times, ok := userMatchConfigIdArr[matchConfigIds]; ok {
				usedTimes += StrToInt(times.(string))
			}
		}
	}
	return usedTimes
}

// 获取用户剩余免费次数
func (f *freeTimeService) GetRemainingTimes(matchConfigId, mid int) int {
	totalTimes := f.GetTotalTimes(matchConfigId)
	var ret int
	if totalTimes > 0 {
		usedTimes := f.GetUsedTimes(matchConfigId, mid)
		ret = totalTimes - usedTimes
	}
	return int(math.Max(float64(ret), 0))
}
