package service

import (
	"PGMatch/app/entity"
	"strings"
	"PGMatch/app/proto"
	"time"
	"github.com/astaxie/beego/logs"
	"reflect"
	"encoding/json"
	"encoding/base64"
)

type BaseMatch struct {
	App        uint64                 // 省包 id
	RegionId   uint32                 // 地区 id
	GameId     uint64                 // 游戏 id
	Version    int                    // 版本号
	Vertical   bool                   // 是否竖屏
	OnlineInfo *entity.OnlineInfo     // 用户线上信息
	Response   *proto.ListsResponse   // 返回值
	Games      map[string]interface{} // 游戏列表
	GameWeight map[string]interface{} // 游戏权重列表
	MatchList  map[string]interface{} // 比赛列表
	curr       CurrentMatch           // 当前比赛
	TAG        string                 // 该结构体名称
}

// 当前比赛
type CurrentMatch struct {
	matchId int64      // 比赛 id
	bean    *MatchBean // 比赛信息
}

const (
	// 获取比赛信息 key
	matchInfoKey = "SRV_COMMON_MATCH_HASH"

	matchPlayersKey = "MATCH_USER_COUNT"

	timingMatchCompletedPlayers = "FIXED_MATCH_THROUGH_SHOW_NUMBER"
)

var NormalRoll *uint32

// 已经申请的定时赛的起始排序号
var applyTimingRoll *uint32

// 初始化 BaseMatch
func NewBaseMatch(mid, appId, gameId uint64, regionId uint32, sessionId string, resp *proto.ListsResponse) *BaseMatch {
	baseMatch := &BaseMatch{
		TAG:      "BaseMatch : ",
		App:      appId,
		RegionId: regionId,
		GameId:   gameId,
	}
	baseMatch.curr = CurrentMatch{}
	baseMatch.Response = resp
	baseMatch.Response.List = map[uint32]proto.List{}
	var normal uint32 = 0
	NormalRoll = &normal
	var applyRoll uint32 = 200
	applyTimingRoll = &applyRoll
	baseMatch.Vertical = baseMatch.isVertical()
	baseMatch.OnlineInfo = SessionService.GetOnlineInfo(mid)
	baseMatch.Version = HallVerToLang(baseMatch.OnlineInfo.Version)
	if baseMatch.OnlineInfo.Mid == 0 {
		baseMatch.OnlineInfo = SessionService.GetOnlineInfo(sessionId)
	}
	// 初始化 games 配置信息
	baseMatch.Games = GenericConfService.GetGamesConf()
	params := entity.GenericConfParam{Api: 0}
	// 初始化 gameWeight 配置信息
	baseMatch.GameWeight = GenericConfService.GetGameWeightConf(params)
	return baseMatch
}

// 判断当前 app 是否为竖屏
func (bm *BaseMatch) isVertical() bool {
	verticalApps := bm.getApps()
	for _, verticalAppId := range verticalApps {
		if verticalAppId == uint32(bm.App) {
			return true
		}
	}
	return false
}

// 获取竖屏 app 列表
func (bm *BaseMatch) getApps() map[string]uint32 {
	apps := map[string]uint32{}
	//TODO CFG("verticalapps", "apps")
	return apps
}

// 获取用户线上信息
// id : mid int | sessionId string
func (bm *BaseMatch) getOnlineInfo(id interface{}) *entity.OnlineInfo {
	var oi entity.OnlineInfo
	return &oi
}

// 处理图片链接(安卓返回HTTP图片链接)
func (bm *BaseMatch) getIconUrl(url string) (iconUrl string) {
	android := ClientIsAndroid(bm.App)
	iconUrl = url
	if android && bm.Version < 800 {
		iconUrl = strings.Replace("https", "http", url, -1)
	}
	return
}

// 明确这个比赛是否有串联的 configId，如果有，取server算出的人数
// @param repeatedConfigIds 	重复的 configId 数组
// @param matchConfigId 		该比赛 configId
// @param configIdMatchIdMap 	configId 和 matchId 的映射关系表
// @param enrollmentMap 		报名人数映射表
// @return enrollment 			报名人数
func (bm *BaseMatch) getApplyNum(repeatedConfigIds []string, matchConfigId string,
	configIdMatchIdMap map[string]string, enrollmentMap map[string]interface{}) (enrollment uint64) {
	matchUserRd, err := NewRedisCache(entity.RedisMatchUser)
	if err != nil {
		logs.Error(bm.TAG + "[matchUser] Redis Init Failed!")
	}
	matchId := configIdMatchIdMap[matchConfigId]
	// 如果没有打通，则返回本场比赛人数
	if len(repeatedConfigIds) == 0 {
		applyNum := matchUserRd.HGet(matchPlayersKey, "match"+matchId)
		if applyNum != nil {
			enrollment = StrToUint64(ByteToStr(applyNum.([]byte)))
		}
		return
	}

	serverConfRd, err := NewRedisCache(entity.RedisServerConf)
	if err != nil {
		logs.Error(bm.TAG + "[serverConf] Redis Init Failed!")
	}
	enrollmentMapI := serverConfRd.HGetAll(timingMatchCompletedPlayers)
	enrollments := enrollmentMapI.(map[string]string)
	for _, repeatedConfigId := range repeatedConfigIds {
		players := 0
		if v, ok := enrollments[repeatedConfigId]; ok {
			players = StrToInt(v)
		}
		if matchConfigId == repeatedConfigId {
			enrollmentMap[repeatedConfigId] = players
			enrollment = enrollmentMap[repeatedConfigId].(uint64)
		} else {
			enrollmentMap[repeatedConfigId] = players
		}
	}
	return
}

// 添加列表到 Response
func (bm *BaseMatch) AddToResponseList(matchList *proto.List) {
	if bm.Response.List == nil {
		logs.Debug(bm.TAG + "AddToResponseList() Init List")
		bm.Response.List = map[uint32]proto.List{}
	}
	// 定时赛已经报名
	bm.Response.List[*NormalRoll] = *matchList
	*NormalRoll++
}

//----------------- 单个比赛调用 -----------------
// 必须赋值 curr

// 加载当前比赛
func (bm *BaseMatch) LoadCurrentMatch(matchId int64, jsonData interface{}) {
	bm.curr.matchId = matchId
	bm.curr.bean = NewMatchBean(jsonData.([]byte))
	logs.Info(bm.TAG+"Load current match success!\n MatchId[%v]", matchId)
}

// 获取当前比赛配置信息
func (bm *BaseMatch) GetCurrConf(confName string) interface{} {
	return bm.curr.bean.GetConf(confName)
}

// 比赛地区限制
// 此方法有的比赛类型没有该过滤，使用开关控制
func (bm *BaseMatch) FilterRegion() bool {
	regionId := bm.RegionId
	regionsObj := bm.curr.bean.GetConf(regionLimitKey)
	regionIds, ok := regionsObj.([]string)
	if !ok {
		return true
	}
	for _, matchRegionId := range regionIds {
		if matchRegionId == Uint32ToStr(regionId) {
			return true
		}
	}
	return false
}

// 横竖屏过滤(0:横屏App 1:横竖通用 2:竖屏App)
func (bm *BaseMatch) FilterScreenDirection() bool {
	vertical := bm.Vertical
	ScreenDirection := bm.curr.bean.GetConf(screenDirKey)
	direction, ok := ScreenDirection.(string)
	if !ok {
		return true
	}
	intDir := StrToInt(direction)
	if (!vertical || intDir != 0) && (vertical || intDir != 2) {
		return true
	}
	return false
}

// iOS 送审暂时屏蔽比赛
func (bm *BaseMatch) FilterIOSAudit() bool {
	appParams := entity.GenericConfParam{AppId: int(bm.App)}
	provinceParams := entity.GenericConfParam{AppId: bm.OnlineInfo.App_id}
	appSwitchI := GenericConfService.GetControlConf(appParams)
	provinceSwitchI := GenericConfService.GetControlConf(provinceParams)
	provinceSwitch := provinceSwitchI.(map[string]string)
	appSwitch := appSwitchI.(map[string]string)
	proHideFlag := 0
	proHideHallFlag := 0
	if hideMatchFlag, ok := provinceSwitch["hideMatchFlag"]; ok {
		if hideMatchFlag != "" && hideMatchFlag != "0" {
			proHideFlag = StrToInt(hideMatchFlag)
		}
	}
	if hideMatchHallVersion, ok := provinceSwitch["hideMatchHallVersion"]; ok {
		if hideMatchHallVersion != "" && hideMatchHallVersion != "0" {
			proHideHallFlag = StrToInt(hideMatchHallVersion)
		}
	}
	appHideFlag := 0
	if hideMatchFlag, ok := appSwitch["hideMatchFlag"]; ok {
		if hideMatchFlag != "" && hideMatchFlag != "0" {
			appHideFlag = StrToInt(hideMatchFlag)
		}
	}
	appHideHallFlag := 0
	if hideMatchHallVersion, ok := appSwitch["hideMatchHallVersion"]; ok {
		if hideMatchHallVersion != "" && hideMatchHallVersion != "0" {
			appHideHallFlag = StrToInt(hideMatchHallVersion)
		}
	}
	if (proHideFlag != 1 || proHideHallFlag != bm.Version) && (appHideFlag != 1 || appHideHallFlag != bm.Version) {
		return true
	}
	Name := bm.curr.bean.GetConf(nameKey)
	matchName, ok := Name.(string)
	if !ok {
		return true
	}
	if !strings.Contains(matchName, "金币") || !strings.Contains(matchName, "游戏币") || !strings.Contains(matchName, "银币") {
		return true
	}
	return false
}

// 白名单过滤
func (bm *BaseMatch) FilterWhiteList() bool {
	onlineMid := bm.OnlineInfo.Mid
	EnableWhiteList := bm.curr.bean.GetConf(whiteListKey)
	if EnableWhiteList == nil {
		return true
	}
	whiteList, ok := EnableWhiteList.(map[string]interface{})
	if !ok {
		return true
	}
	open, ok := whiteList["isopen"]
	if !ok {
		return true
	}
	if open != 1 {
		return true
	}
	midList, ok := whiteList["mids"]
	if !ok {
		return true
	}
	midMap, ok := midList.(map[string]string)
	if !ok {
		return true
	}
	for _, mid := range midMap {
		if StrToInt(mid) == onlineMid {
			return false
		}
	}
	return true
}

// 比赛时间过滤
func (bm *BaseMatch) FilterMatchData() bool {
	MatchStartDate := bm.curr.bean.GetConf(dateStart)
	MatchEndDate := bm.curr.bean.GetConf(dateEnd)
	startData, ok := MatchStartDate.(string)
	if !ok {
		return true
	}
	endData, ok := MatchEndDate.(string)
	if !ok {
		return true
	}
	nowData := time.Now().Format("20060102")
	matchStart := strings.Replace(startData, "-", "", -1)
	matchEnd := strings.Replace(endData, "-", "", -1)
	if nowData >= matchStart && nowData <= matchEnd {
		return true
	}
	return false
}

// 比赛类型过滤
func (bm *BaseMatch) FilterMatchType() bool {
	matchType := bm.curr.bean.GetConf(matchTypeKey)
	switch matchType {
	case 2: // 周循环赛
		return bm.FilterRoundRobinWeek()
	case 3: // 24小时循环赛
		return bm.FilterRoundRobin24Hour()
	case 4: // 间隔循环赛
		return bm.FilterRoundRobinInterval()
	default:
		logs.Debug("Match matchTypeKey %v is not exist!", matchType)
	}
	return true
}

// 周循环赛
func (bm *BaseMatch) FilterRoundRobinWeek() bool {
	week := time.Now().Weekday()
	weekLoopI := bm.curr.bean.GetConf(matchLoopKey)
	weekLoop, ok := weekLoopI.(string)
	if !ok {
		logs.Error("WeekLoop's type is %v !", reflect.TypeOf(weekLoopI))
		return true
	}
	weeks := strings.Split(weekLoop, ",")
	for _, weekday := range weeks {
		if StrToInt(weekday) == int(week) {
			return false
		}
	}
	return true
}

// 过滤24小时循环赛
func (bm *BaseMatch) FilterRoundRobin24Hour() bool {
	startTimeI := bm.curr.bean.GetConf(timeStart)
	endTimeI := bm.curr.bean.GetConf(timeEnd)
	startDate := bm.curr.bean.GetConf(dateStart)
	endDate := bm.curr.bean.GetConf(dateEnd)
	if startTimeI == nil || endTimeI == nil || startDate == nil || endDate == nil {
		logs.Warn("time is nil !!! Date[%v - %v] Time[%v - %v]\n", startTimeI, endTimeI, startDate, endDate)
		return true
	}
	startTime := startDate.(string) + startTimeI.(string)
	endTime := endDate.(string) + endTimeI.(string)
	nowTime := Int64ToStr(time.Now().Unix())
	if nowTime >= startTime && nowTime <= endTime {
		return true
	}
	return false
}

// 过滤间隔循环赛
func (bm *BaseMatch) FilterRoundRobinInterval() bool {
	startDateUnix := FormatTime(bm.curr.bean.GetConf(dateStart).(string)).Unix()
	round := bm.curr.bean.GetConf(matchRoundKey).(int) + 1
	diff := (time.Now().Unix() - startDateUnix) / (24 * 60 * 60) % int64(round)
	if diff > 0 {
		return false
	}
	return true
}

// 版本号过滤
func (bm *BaseMatch) FilterVersion() bool {
	appId := bm.App
	regionId := bm.RegionId
	ver := bm.Version
	var (
		androidMaxVer int
		androidMinVer int
		iOSMaxVer     int
		iOSMinVer     int
	)
	RegionsVersion := bm.curr.bean.GetConf(regionLimitKey)
	versions, ok := RegionsVersion.([]interface{})
	if !ok {
		logs.Error("versions is %T, not map[string]interface{}", versions)
		return true
	}
	componseRegion := func(version map[string]interface{}) bool {
		if !ok {
			logs.Debug("Android or iOS region version[%T] is not map[string]interface{} !", version)
			return false
		}
		MaxVersionAnd, ok := version["MaxVersionAnd"]
		if !ok {
			logs.Error("version['MaxVersionAnd'] is nil!")
			return false
		}
		MinVersionAnd, ok := version["MinVersionAnd"]
		if !ok {
			logs.Error("version['MinVersionAnd'] is nil!")
			return false
		}
		MaxVersionIos, ok := version["MaxVersionIos"]
		if !ok {
			logs.Error("version['MaxVersionIos'] is nil!")
			return false
		}
		MinVersionIos, ok := version["MinVersionIos"]
		if !ok {
			logs.Error("version['MinVersionIos'] is nil!")
			return false
		}
		androidMaxVer = MaxVersionAnd.(int)
		androidMinVer = MinVersionAnd.(int)
		iOSMaxVer = MaxVersionIos.(int)
		iOSMinVer = MinVersionIos.(int)
		return true
	}
	// 地区 id 对应版本
	regionVerI := versions[regionId]
	if regionVerI == nil {
		// 默认加载 index = 0
		regionVerI = versions[0]
	}
	regionVer, ok := regionVerI.(map[string]interface{})
	if !ok {
		return true
	} else {
		// index = regionId 存在则加载 region 对应配置
		isPassed := componseRegion(regionVer)
		if !isPassed {
			return true
		}
	}

	if ClientIsAndroid(appId) { // Android
		if ver < androidMinVer || (androidMaxVer > 0 && ver > androidMaxVer) {
			return false
		}
	} else if ClientIsIos(appId) { // iOS
		if ver < iOSMinVer || (iOSMaxVer > 0 && ver > iOSMaxVer) {
			return false
		}
	}

	return true
}

// ------------------------ 获取当前比赛配置

// 当前比赛已经报名，增加到 Apply
func (bm *BaseMatch) PushToApply() {
	bm.Response.Apply = append(bm.Response.Apply, bm.GetCurrentMatchId())
}

// 获取当前比赛的 id
func (bm *BaseMatch) GetCurrentMatchId() (matchId uint64) {
	matchId = uint64(bm.curr.matchId)
	return
}

// 获取比赛名称
func (bm *BaseMatch) GetMatchName() (matchName string) {
	matchNameI := bm.curr.bean.GetMatchName()
	if matchNameI != nil {
		matchName = matchNameI.(string)
	}
	return
}

// 获取比赛简称
func (bm *BaseMatch) GetMatchShortName() (matchShortName string) {
	matchShortNameI := bm.curr.bean.GetMatchShortName()
	if matchShortNameI != nil {
		matchShortName = matchShortNameI.(string)
		return
	}
	matchNameI := bm.curr.bean.GetMatchName()
	if matchNameI != nil {
		matchShortName = matchNameI.(string)
	}
	return
}

// 获取开始日期+时间
func (bm *BaseMatch) GetStartDateTime() (startDateTime uint64) {
	startDateI := bm.curr.bean.GetStartDate()
	startTimeI := bm.curr.bean.GetStartTime()
	if startTimeI != nil && startDateI != nil {
		startDateTimeStr := startDateI.(string) + startTimeI.(string)
		startDateTime = uint64(FormatTime(startDateTimeStr).Unix())
	}
	return
}

// 获取比赛 icon
func (bm *BaseMatch) GetMatchIcon() (matchIcon string) {
	matchIconI := bm.curr.bean.GetMatchIcon()
	if matchIconI != nil {
		matchIcon = bm.getIconUrl(matchIconI.(string))
	}
	return
}

// 获取游戏 id 排序
func (bm *BaseMatch) GetGameSort() (gameSort uint64) {
	gameIdI := bm.curr.bean.GetGameId()
	gameIdStr, ok := gameIdI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetGameId() Failed!")
		return 0
	}
	gameSort = StrToUint64(gameIdStr)
	gameWeight, ok := bm.GameWeight[gameIdStr]
	if !ok {
		logs.Error(bm.TAG+"GetGameSort() gameWeight not has %v", gameIdStr)
		return
	}
	gameSort = StrToUint64(gameWeight.(string))
	return
}

func (bm *BaseMatch) GetGameId() (gameId uint64) {
	gameIdI := bm.curr.bean.GetGameId()
	gameIdStr, ok := gameIdI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetGameId() Failed!")
		return 0
	}
	gameId = StrToUint64(gameIdStr)
	return
}

// 获取最低开赛人数
func (bm *BaseMatch) GetMinimumStarts() (minimumStarts string) {
	totalNumI := bm.curr.bean.GetTotalNum()
	minimumStarts, ok := totalNumI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetMinimumStarts() Failed!")
		return "0"
	}
	return
}

// 获取所有报名人数
func (bm *BaseMatch) GetAllEnrollment() (allEnrollment int) {
	totalNumI := bm.curr.bean.GetTotalNum()
	allEnrollmentStr, ok := totalNumI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetAllEnrollment() Failed!")
		return 0
	}
	allEnrollment = StrToInt(allEnrollmentStr)
	return
}

// 获取允许等待时间
func (bm *BaseMatch) GetAllowWaitTime() (allowWaitingTime uint32) {
	allowWaitTimeI := bm.curr.bean.GetAllowWaitTime()
	allowWaitingTimeStr, ok := allowWaitTimeI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetAllowWaitTime() Failed!")
		return 0
	}
	allowWaitingTime = StrToUint32(allowWaitingTimeStr)
	return
}

// 获取是否完成比赛
func (bm *BaseMatch) GetIsFinishedMatch() (isFinishMatch float64) {
	isFinishMatchI := bm.curr.bean.GetIsFhMatch()
	isFinishMatch, ok := isFinishMatchI.(float64)
	if !ok {
		logs.Error(bm.TAG + "GetIsFinishedMatch() Failed!")
		return 0
	}
	return
}

// 获取运营描述信息
func (bm *BaseMatch) GetDescription() (desc string) {
	rewardDescribeI := bm.curr.bean.GetRewardDescription()
	desc, ok := rewardDescribeI.(string)
	if !ok {
		logs.Error(bm.TAG + "GetDescription() Failed!")
		return ""
	}
	return
}

// 获取比赛标签
func (bm *BaseMatch) GetMatchTags() (matchTags string) {
	tags := matchTagsKey
	if bm.Vertical {
		tags = verticalMatchTags
	}
	mTags := bm.GetCurrConf(tags)
	if mTags != nil {
		if tags, ok := mTags.([]interface{}); ok {
			tagSlice := make([]string, len(tags))
			for i, v := range tags {
				tagSlice[i] = v.(string)
			}
			matchTags = strings.Join(tagSlice, ",")
		} else {
			matchTags = mTags.(string)
		}
	}
	return
}

// 获取广告图 icon
func (bm *BaseMatch) GetAdvertIcon() (advertI string) {
	advertIcon := matchAdvertIconKey
	if bm.Vertical {
		advertIcon = verticalMatchAdvertIcon
	}
	mAdvertIcon := bm.GetCurrConf(advertIcon)
	if mAdvertIcon != nil {
		advertI = bm.getIconUrl(mAdvertIcon.(string))
	}
	return
}

// 获取比赛 icon 权重
func (bm *BaseMatch) GetIconWeight() (iconW string) {
	iconWeight := iconWeightKey
	if bm.Vertical {
		iconWeight = verticalIconWeight
	}
	mIconWeight := bm.GetCurrConf(iconWeight)
	if mIconWeight != nil {
		iconWeight = mIconWeight.(string)
	}
	return
}

// 获取报名门槛信息
func (bm *BaseMatch) GetThresholdInfo() (thresholdArray [2]float64) {
	// 第一种情况
	thresholdI := bm.curr.bean.GetThreshold()
	threshold, ok := thresholdI.(float64)
	if ok {
		thresholdTypeI := bm.curr.bean.GetThresholdType()
		thresholdType, ok := thresholdTypeI.(float64)
		if ok {
			thresholdArray[0] = thresholdType
		}
		thresholdArray[1] = threshold
		return
	}
	// 第二种情况
	thresholdMap, ok := thresholdI.(map[string]string)
	if ok {
		// 门槛类型
		if thresholdTypeStr, ok := thresholdMap["type"]; ok {
			thresholdArray[0] = StrToFloat64(thresholdTypeStr)
		}
		// 报名门槛
		if thresholdStr, ok := thresholdMap["value"]; ok {
			thresholdArray[1] = StrToFloat64(thresholdStr)
		}
	}
	return
}

// 获取新版报名费
func (bm *BaseMatch) GetNewFee() (feeSlice []proto.Fee) {
	newFeeI := bm.curr.bean.GetNewFee()
	if newFeeI == nil {
		return nil
	}
	newFee, ok := newFeeI.(string)
	if ok {
		json.Unmarshal([]byte(newFee), &feeSlice)
	}
	return
}

// 获取比赛标签
func (bm *BaseMatch) GetMatchLabel() (matchLabelMap map[string]string) {
	matchLabelI := bm.curr.bean.GetMatchLabel()
	if matchLabelI == nil {
		return nil
	}
	matchLabel, ok := matchLabelI.(string)
	if ok {
		json.Unmarshal([]byte(matchLabel), &matchLabelMap)
	}
	return
}

// 获取限制开关
func (bm *BaseMatch) GetLimitSwitch() (limitSwitch bool) {
	openLimit := bm.curr.bean.GetOpenLimit()
	// 开关打开
	if openLimit != nil && openLimit.(float64) == 1 {
		limitSwitch = true
	}
	return
}

// 获取限制开关
func (bm *BaseMatch) GetRecommendSwitch() (recommendSwitch bool) {
	openRecommend := bm.curr.bean.GetOpenRecommend()
	// 开关打开
	if openRecommend != nil && openRecommend.(float64) == 1 {
		recommendSwitch = true
	}
	return
}

// 获取入口限制信息
func (bm *BaseMatch) GetEntranceLimitInfo() (entranceLimitInfo map[string]interface{}) {
	entranceLimitI := bm.curr.bean.GetEntranceLimit()
	entranceLimitInfo, ok := entranceLimitI.(map[string]interface{})
	if ok {
		return entranceLimitInfo
	}
	return
}

// 获取入口推荐信息
func (bm *BaseMatch) GetEntranceRecommendInfo() (entranceRecommendInfo []interface{}) {
	entranceRecommendInfoI := bm.curr.bean.GetEntranceRecommend()
	if entranceRecommendInfoStr, ok := entranceRecommendInfoI.(string); ok {
		json.Unmarshal([]byte(entranceRecommendInfoStr), entranceRecommendInfo)
	} else {
		if entranceRecommendInfoMap, ok := entranceRecommendInfoI.([]interface{}); ok {
			entranceRecommendInfo = entranceRecommendInfoMap
		}
	}
	return
}

// 获取冠军奖励
func (bm *BaseMatch) GetChampionRewardsInfo() (championRewardsDesc string) {
	championRewardsMap := map[string]map[string]map[string]interface{}{}
	championRewardsI := bm.curr.bean.GetChampionRewards()
	json.Unmarshal([]byte(championRewardsI.(string)), &championRewardsMap)
	if championRewards, ok := championRewardsMap["1"]; ok {
		if championReward, ok := championRewards["1"]; ok {
			desc := championReward["desc"].(string)
			bytes, err := base64.StdEncoding.DecodeString(desc)
			if err != nil {
				logs.Error(bm.TAG + "GetChampionRewards() decode failed!")
			}
			championRewardsDesc = strings.Trim(string(bytes[:]), " ")
		}
	}
	return
}
