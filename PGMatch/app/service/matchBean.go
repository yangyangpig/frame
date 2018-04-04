package service

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

// 比赛
type MatchBean struct {
	Desc    string                 `json:"desc"`  // 比赛 title
	Value   map[string]interface{} `json:"value"` // 比赛配置内容
	GetConf confFunc                              // 获取配置内容
}

type confFunc func(confName string) interface{}

// 比赛配置的 key 索引配置
const (
	dateStart               = "MatchStartDate"
	dateEnd                 = "MatchEndDate"
	timeStart               = "MatchStartTime"
	timeEnd                 = "MatchEndTime"
	matchTypeKey            = "MatchType"
	whiteListKey            = "EnableWhiteList"
	screenDirKey            = "ScreenDirection"
	regionLimitKey          = "MatchRegions"
	nameKey                 = "Name"
	matchLoopKey            = "MatchLoop"
	matchRoundKey           = "MatchLoopDay"
	perDayTimeSegmentsKey   = "PerDayTimeSegments"
	segmentsTimeStart       = "SegmentsFromTime"
	segmentsTimeEnd         = "SegmentsToTime"
	matchConfigIdKey        = "MatchConfigId"
	matchShortNameKey       = "MatchShortName"
	matchIconKey            = "MatchIcon"
	adIconKey               = "AdIcon"
	listSortKey             = "ListSort"
	matchUserCountKey       = "MatchUserCount"
	totalNumKey             = "TotalNum"
	rewardDescribeKey       = "RewardDescribe"
	verticalIconWeight      = "VerticalIconWeight"
	verticalMatchAdvertIcon = "VerticalMatchAdvertIcon"
	verticalMatchTags       = "VerticalMatchTags"
	iconWeightKey           = "IconWeight"
	matchAdvertIconKey      = "MatchAdvertIcon"
	matchTagsKey            = "MatchTags"
	thresholdKey            = "Threshold"
	feeKey                  = "Fee"
	newFeeKey               = "SignUpFeeInfo_2"
	tableUserCountKey       = "TableUserCount"
	matchLabelKey           = "mclabel"
	gameIdKey               = "GameId"
	allowWaitTimeKey        = "AllowWaitTime"
	isFinishMatchKey        = "isFhMatch"
	championRewardsKey      = "Prize_List"
	awardIconKey            = "AwardIcon"
	thresholdTypeKey        = "ThresholdType"
	openLimitKey            = "OpenLimit"
	openRecommendKey        = "OpenRecomm"
	entranceLimitKey        = "EntranceLimit"
	entranceRecommendKey    = "EntranceRecommend"
)

// 实例化一个 MatchBean
func NewMatchBean(matchJsonStr []byte) *MatchBean {
	var bean MatchBean
	json.Unmarshal(matchJsonStr, &bean)
	bean.GetConf = bean.getConfInfo()
	return &bean
}

// 获取配置信息
func (mb *MatchBean) getConfInfo() confFunc {
	return func(confName string) interface{} {
		confMap := mb.Value
		if confMap == nil || len(confMap) == 0 {
			logs.Error("Match configuration does not exist, the specific value is %v!", confMap)
			return nil
		}
		confItemI, ok := confMap[confName]
		if !ok {
			logs.Error("Change the configuration[%v] does not exist!", confName)
			return nil
		}
		confItem, ok := confItemI.(map[string]interface{})
		if !ok {
			logs.Error("Type conversion failed! [%v]", confItemI)
			return nil
		}
		confInfo, ok := confItem["value"]
		if !ok {
			logs.Error("Match configuration is not configured value!")
			return confInfo
		}
		confDesc, ok := confItem["desc"]
		if ok {
			logs.Debug("ConfName: %v\n ConfDesc: %v ConfInfo: %v", confName, confDesc, confInfo)
		}
		return confInfo
	}
}

// 获取报名门槛类型
func (mb *MatchBean) GetThresholdType() interface{} {
	thresholdType := mb.GetConf(thresholdTypeKey)
	return thresholdType
}

// 获取报名门槛
func (mb *MatchBean) GetThreshold() interface{} {
	threshold := mb.GetConf(thresholdKey)
	return threshold
}

// 获取新版报名费
func (mb *MatchBean) GetNewFee() interface{} {
	newFee := mb.GetConf(newFeeKey)
	return newFee
}

// 获取比赛标签
func (mb *MatchBean) GetMatchLabel() interface{} {
	matchLabel := mb.GetConf(matchLabelKey)
	return matchLabel
}

// 获取限制开关
func (mb *MatchBean) GetOpenLimit() interface{} {
	openLimit := mb.GetConf(openLimitKey)
	return openLimit
}

// 获取推荐开关
func (mb *MatchBean) GetOpenRecommend() interface{} {
	openRecommend := mb.GetConf(openRecommendKey)
	return openRecommend
}

// 获取入口限制
func (mb *MatchBean) GetEntranceLimit() interface{} {
	entranceLimit := mb.GetConf(entranceLimitKey)
	return entranceLimit
}

// 获取入口推荐
func (mb *MatchBean) GetEntranceRecommend() interface{} {
	entranceRecommend := mb.GetConf(entranceRecommendKey)
	return entranceRecommend
}

// 获取冠军奖励
func (mb *MatchBean) GetChampionRewards() interface{} {
	championRewards := mb.GetConf(championRewardsKey)
	return championRewards
}

// 获取运营描述
func (mb *MatchBean) GetRewardDescription() interface{} {
	rewardDescribe := mb.GetConf(rewardDescribeKey)
	return rewardDescribe
}

// 获取允许等待时间
func (mb *MatchBean) GetAllowWaitTime() interface{} {
	allowWaitTime := mb.GetConf(allowWaitTimeKey)
	return allowWaitTime
}

// 获取完成比赛
func (mb *MatchBean) GetIsFhMatch() interface{} {
	isFinishMatch := mb.GetConf(isFinishMatchKey)
	return isFinishMatch
}

// 获取所有报名人数
func (mb *MatchBean) GetTotalNum() interface{} {
	totalNum := mb.GetConf(totalNumKey)
	return totalNum
}

// 获取最低开赛人数
func (mb *MatchBean) GetGameId() interface{} {
	gameId := mb.GetConf(gameIdKey)
	return gameId
}

// 获取比赛 icon
func (mb *MatchBean) GetMatchIcon() interface{} {
	matchIcon := mb.GetConf(matchIconKey)
	return matchIcon
}

// 获取开始时间
func (mb *MatchBean) GetStartTime() interface{} {
	startTime := mb.GetConf(timeStart)
	return startTime
}

// 获取结束时间
func (mb *MatchBean) GetEndTime() interface{} {
	endTime := mb.GetConf(timeEnd)
	return endTime
}

// 获取开始日期
func (mb *MatchBean) GetStartDate() interface{} {
	startDate := mb.GetConf(dateStart)
	return startDate
}

// 获取结束日期
func (mb *MatchBean) GetEndDate() interface{} {
	endDate := mb.GetConf(dateEnd)
	return endDate
}

// 获取比赛名称
func (mb *MatchBean) GetMatchName() interface{} {
	matchNameI := mb.GetConf(nameKey)
	return matchNameI
}

// 获取比赛简称
func (mb *MatchBean) GetMatchShortName() interface{} {
	matchShortNameI := mb.GetConf(matchShortNameKey)
	return matchShortNameI
}