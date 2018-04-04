// 定时赛服务
package service

import (
	"PGMatch/app/entity"
	"time"
	"encoding/json"
	"fmt"
	"PGMatch/app/proto"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"math"
	"strings"
	"strconv"
)

type TimingMatchService struct {
	*BaseMatch
	configIdMatchIdMap map[string]string
	gameInfoService    entity.GameInfoInterface
	whetherAddRobot    bool
	matchItem          map[string]interface{}
	mixApplyNums       map[string]interface{}
}

// 定时赛比赛类型
const (
	timingMatch         = 3 // 定时赛比赛类型
	timingMatchListKey  = "FIXED_MATCH_LIST_HASH_"
	matchUserListPrefix = "TMFX_MATCH_USER"
)

// 构造定时赛列表服务
func NewTimingMatchService(base *BaseMatch) entity.MatchLister {
	base.TAG = "TimingMatchService : "
	mixApplyNums := map[string]interface{}{}
	configIdMatchIdMap := map[string]string{}
	matchItem := map[string]interface{}{}
	return &TimingMatchService{
		base,
		configIdMatchIdMap,
		nil,
		false,
		matchItem,
		mixApplyNums,
		}
}

// 初始化
func (tm *TimingMatchService) Init() {
	// 初始化比赛列表
	regionKey := Uint32ToStr(tm.RegionId)
	matchListData := MatchConfService.GetAllConfBy(timingMatchListKey + regionKey)
	// 列表为空
	tm.MatchList = map[string]interface{}{}
	if matchListData != nil {
		// 按照开始时间排序
		t := 100
		nowTime := time.Now().Unix()
		for matchId, matchItemJson := range matchListData.(map[string]interface{}) {
			var matchItem entity.TimingMatchData
			json.Unmarshal(matchItemJson.([]byte), &matchItem)
			startTime := matchItem.Begintime
			if nowTime >= int64(startTime) {
				continue
			}
			t++
			matchListKey := fmt.Sprintf("%d%d", int(startTime), t)
			tm.MatchList[matchListKey] = []byte(matchItemJson.(string))
			tm.configIdMatchIdMap[matchItem.Configid] = matchId
		}
	}

	// 初始化游戏记录服务
	tm.gameInfoService = NewGameInfoService(tm.OnlineInfo.Mid)
}

// 过滤
func (tm *TimingMatchService) Filter() {
	matchList := tm.MatchList
	if matchList == nil || len(matchList) == 0 {
		return
	}
	keySlice := KeySortedSlice(matchList)
	tm.Response.Apply = make([]uint64, len(keySlice))
	for _, key := range keySlice {
		jsonI, ok := tm.MatchList[key]
		if !ok {
			continue
		}
		jsonStr := jsonI.(string)
		json.Unmarshal([]byte(jsonStr), &tm.matchItem)
		// 定时赛信息过滤
		if tm.filterMatchInfo(tm.matchItem) {
			continue
		}
		tm.setCurrMatch()
		matchGameId := tm.GetCurrentMatchId()
		// gameId 过滤
		if tm.GameId > 0 && tm.GameId != matchGameId {
			continue
		}

		// iOS 送审过滤
		if !tm.FilterIOSAudit() {
			continue
		}

		// 白名单过滤
		if !tm.FilterWhiteList() {
			continue
		}
		// 是否添加机器人
		addrobotflag := tm.matchItem["addrobotflag"].(int)
		if addrobotflag == 1 {
			tm.whetherAddRobot = true
		}
		// 加载比赛列表信息
		tm.SetResponse()
	}
}

// 设置当前比赛
func (tm *TimingMatchService) setCurrMatch() {
	// 有效比赛
	// 赋值当前比赛
	matchId, ok := tm.matchItem["matchid"]
	if !ok {
		panic(tm.TAG + "setCurrMatch() MatchId is nil!")
	}
	matchBean := MatchConfService.GetConfBy(matchInfoKey, "match"+Int64ToStr(matchId.(int64)))
	if matchBean == nil {
		panic(tm.TAG + "setCurrMatch() matchBean is nil!")
	}
	// 装载当前比赛
	tm.LoadCurrentMatch(matchId.(int64), matchBean)
}

// 加载列表信息
func (tm *TimingMatchService) SetResponse() {
	matchInfo := proto.List{}
	matchInfo.Type = timingMatch
	matchInfo.Id = tm.GetCurrentMatchId()
	// 比赛名称
	matchInfo.Mname = tm.GetMatchName()
	// 比赛简称
	matchInfo.Sname = tm.GetMatchShortName()
	// 比赛 icon
	matchInfo.Micon = tm.GetMatchIcon()
	// adIcon
	matchInfo.AdIcon = tm.getAdIcon()
	// 列表排序
	matchInfo.ListSort = tm.getListSort()
	// 开始时间
	matchInfo.Stime = tm.GetStartDateTime()
	// 结束时间
	matchInfo.Endtime = tm.getEndTime()
	// 游戏 id
	matchInfo.Gameid = tm.GetGameId()
	// 游戏排序
	matchInfo.Gamesort = tm.GetGameSort()
	// 最低开赛人数
	matchInfo.Requestnum = tm.GetMinimumStarts()
	// 报名人数
	enrollment := tm.getEnrollment()
	// 已报名人数(包括机器人)
	matchInfo.Applynum = tm.adjustEnrollment(enrollment, tm.whetherAddRobot)
	// 比赛分场数
	matchInfo.Matchpartitions = tm.getMatchPartitions(enrollment)
	// 所有报名人数
	matchInfo.Allapplynum = enrollment
	// 允许等待时间
	matchInfo.Allowwaittime = tm.GetAllowWaitTime()
	// 完成时间？？？
	matchInfo.Isfhmatch = tm.GetIsFinishedMatch()
	// 冠军奖励
	matchInfo.Champion = tm.GetChampionRewardsInfo()
	// 加载循环间隔
	loopInfo := tm.getLoopInterval()
	matchInfo.Looptype = loopInfo.LoopType
	matchInfo.Loopinterval = loopInfo.LoopInterval
	matchInfo.Loopendtime = loopInfo.LoopendTime
	matchInfo.Loopintervaltime = loopInfo.LoopIntervalTime
	matchInfo.Firstbegintime = loopInfo.FirstBeginTime
	// 比赛验证码
	matchInfo.Matchentrycode = tm.getEntryCode()
	// 比赛验证码获取方式
	matchInfo.Matchentryinfo = tm.getEntryInfo()
	// 游戏名称
	matchInfo.Gamename = tm.getGameName()
	// 比赛配置 id
	matchInfo.Configid = tm.getConfigId()
	// 奖励 icon
	matchInfo.RewardUrl = tm.getRewardUrl()
	// 运营描述文字
	matchInfo.RewardDescribe = tm.GetDescription()
	// 加载比赛标签
	matchInfo.Iconweight = tm.GetIconWeight()
	matchInfo.Advicon = tm.GetAdvertIcon()
	matchInfo.Matchtags = tm.GetMatchTags()
	// 报名门槛类型
	thresholdInfo := tm.GetThresholdInfo()
	matchInfo.ThresholdType = thresholdInfo[0]
	// 报名门槛
	matchInfo.Threshold = thresholdInfo[1]
	// 报名费
	matchInfo.Fee = tm.GetNewFee()
	// 更新 status
	if matchInfo.Status != 0 {
		matchInfo.Status = tm.getNewStatus(matchInfo.Stime, matchInfo.Allowwaittime)
	}
	// 比赛标签，玩法标签，场次标签配置
	matchInfo.Mclabel = tm.GetMatchLabel()
	// 获取推荐或限制信息
	gameCurrencyLimit, recommendMatchConfigId := tm.getRecommendLimitInfo()
	matchInfo.Moneylimit = gameCurrencyLimit
	matchInfo.Recommendmatchcfgid = recommendMatchConfigId
	// 根据玩家是否报名更新 apply 和 status
	index, status := tm.getIndexStatus()
	matchInfo.Status = status
	// 添加到返回列表
	tm.AddToResponseList(index, &matchInfo)
}

// 添加到返回列表(覆盖 BaseMatch 函数)
// @param index 报名情况不同索引不同
func (tm *TimingMatchService) AddToResponseList(index *uint32, matchInfo *proto.List) {
	if tm.Response.List == nil {
		logs.Debug(tm.TAG + "AddToResponseList() Init List")
		tm.Response.List = map[uint32]proto.List{}
	}
	// 定时赛已经报名
	tm.Response.List[*index] = *matchInfo
	*index++
}

// 根据玩家是否报名更新 apply 和 status
func (tm *TimingMatchService) getIndexStatus() (index *uint32, status uint32) {
	isApply := tm.checkPlayerSignUp(tm.GetCurrentMatchId())
	index = NormalRoll
	// 报名
	if isApply {
		tm.PushToApply()
		index = applyTimingRoll
		status = 1
	}
	return
}

// 获取新状态
func (tm *TimingMatchService) getNewStatus(startTime uint64, minutes uint32) (status uint32) {
	seconds := 60 * uint64(minutes)
	nowTime := time.Now().Unix()
	// 到达提前进入时间，并且提前5秒
	if uint64(nowTime) >= startTime-seconds-5 {
		status = 4
	}
	return
}

// 获取游戏名称
func (tm *TimingMatchService) getGameName() (gameName string) {
	gameNameI, ok := tm.matchItem["gamename"]
	if !ok {
		logs.Error(tm.TAG + "getGameName() gamename is nil!")
		return ""
	}
	gameName = gameNameI.(string)
	return
}

// 获取配置 id
func (tm *TimingMatchService) getConfigId() (configId string) {
	configIdI, ok := tm.matchItem["configid"]
	if !ok {
		logs.Error(tm.TAG + "getConfigId() configid is nil!")
		return ""
	}
	configId = configIdI.(string)
	return
}

// 获取验证码
func (tm *TimingMatchService) getEntryCode() (code string) {
	codeI, ok := tm.matchItem["matchentrycode"]
	if !ok {
		logs.Error(tm.TAG + "getEntryCode() matchentrycode is nil!")
		return ""
	}
	code = codeI.(string)
	return
}

// 获取验证码获取方式
func (tm *TimingMatchService) getEntryInfo() (info string) {
	infoI, ok := tm.matchItem["matchentryinfo"]
	if !ok {
		logs.Error(tm.TAG + "getEntryInfo() matchentryinfo is nil!")
		return ""
	}
	info = infoI.(string)
	return
}

// 获取 adIcon
func (tm *TimingMatchService) getAdIcon() (adIcon string) {
	adIconI, ok := tm.matchItem["AdIcon"]
	if ok {
		adIcon = adIconI.(string)
	}
	return
}

// 获取列表排序
func (tm *TimingMatchService) getListSort() (listSort string) {
	listSortI, ok := tm.matchItem["ListSort"]
	if !ok {
		logs.Error(tm.TAG + "getListSort() ListSort is nil")
		return ""
	}
	listSort = listSortI.(string)
	return
}

// 获取结束时间
func (tm *TimingMatchService) getEndTime() (endTime uint64) {
	endTimeI, ok := tm.matchItem["endtime"]
	if ok {
		endTime = endTimeI.(uint64)
	}
	return
}

// 获取报名人数
func (tm *TimingMatchService) getEnrollment() (enrollment uint64) {
	// 比赛配置 id
	matchConfigIdI, ok := tm.matchItem["configid"]
	if !ok {
		logs.Error(tm.TAG + "MatchItem not have key[configid]!")
		return
	}
	matchConfigId, ok := tm.mixApplyNums[matchConfigIdI.(string)]
	// 报名人数
	if ok {
		enrollment = StrToUint64(matchConfigId.(string))
	} else {
		configIds := NewTimingMatchCompletedService().GetMixMatchConfigId(matchConfigId.(string))
		enrollment = tm.getApplyNum(configIds.([]string), matchConfigId.(string), tm.configIdMatchIdMap, tm.mixApplyNums)
	}
	return
}

// 获取比赛分区数量
func (tm *TimingMatchService) getMatchPartitions(applyNum uint64) (matchPartitions uint64) {
	PartitionFlag, ok := tm.matchItem["partitionFlag"]
	if ok && PartitionFlag.(string) == "1" {
		maxUserCount := tm.matchItem["maxUserCount"].(string)
		maxPartitions := tm.matchItem["maxPartitions"].(string)
		matchPartitions = uint64(tm.getSubFieldNum(applyNum, StrToInt(maxUserCount), StrToInt(maxPartitions)))
	}
	return
}

// 获取推荐或限制信息
func (tm *TimingMatchService) getRecommendLimitInfo() (gameCurrencyLimit map[string]string, matchConfigId uint64) {
	// 游戏币金额限制
	gameCurrencyLimit = map[string]string{}
	// 获取一条游戏记录
	mid := tm.OnlineInfo.Mid
	gameInfo := NewGameInfoService(mid).GetRecord()
	// 携带的金条数
	goldBarsCarried, ok := gameInfo["crystal"]
	if !ok {
		logs.Error(tm.TAG + "getRecommendLimitInfo goldBarsCarried is nil!")
	}
	// 保险箱金条数
	goldBarsSafeBox, ok := gameInfo["crystalsafebox"]
	if !ok {
		logs.Error(tm.TAG + "getRecommendLimitInfo goldBarsSafeBox is nil!")
	}
	// 总金条数
	goldBars := goldBarsCarried.(int) + goldBarsSafeBox.(int)
	// 旧版本
	if tm.GetLimitSwitch() {
		tm.loadGameCurrencyLimit(goldBars, gameCurrencyLimit, &matchConfigId)
	}
	// 新版本开关(游戏币金额限制没有加载旧版本)
	if tm.GetRecommendSwitch() && len(gameCurrencyLimit) == 0 {
		tm.loadGameCurrencyLimit(goldBars, gameCurrencyLimit, &matchConfigId)
	}
	return
}

// 加载游戏金额限制信息和比赛配置 id
func (tm *TimingMatchService) loadGameCurrencyLimit(goldBars int, gameCurrencyLimit map[string]string, matchConfigId *uint64) {
	// 入口限制信息加载
	entranceLimitInfo := tm.GetEntranceLimitInfo()
	if entranceLimitInfo != nil {
		if entranceLimitType, ok := entranceLimitInfo["type"]; ok {
			gameCurrencyLimit["type"] = entranceLimitType.(string)
		}

		if entranceLimitValue, ok := entranceLimitInfo["value"]; ok {
			gameCurrencyLimit["value"] = entranceLimitValue.(string)
		}
	}
	// 入口推荐信息加载
	entranceRecommendInfo := tm.GetEntranceRecommendInfo()
	if entranceRecommendInfo != nil && len(entranceRecommendInfo) > 0 {
		for _, infoI := range entranceRecommendInfo {
			info, ok := infoI.(map[string]interface{})
			if !ok {
				panic("entranceRecommendInfo load failed!")
			}
			start, ok := info["start"]
			if !ok {
				continue
			}
			end, ok := info["end"]
			if !ok {
				continue
			}
			if goldBars < start.(int) {
				continue
			}
			if end.(string) != "无限大" && goldBars > end.(int) {
				continue
			}
			*matchConfigId = info["matchCfgId"].(uint64)
			break
		}
	}
}

// 覆盖获取比赛标签方法
func (tm *TimingMatchService) GetMatchTags() (matchTag string) {
	matchTagsKey := "matchtags"
	if tm.Vertical {
		matchTagsKey = "vmatchtags"
	}
	matchTagsI, ok := tm.matchItem[matchTagsKey]
	if ok {

		if matchTags, ok := matchTagsI.([]string); ok {
			matchTag = strings.Join(matchTags, ",")
		} else {
			matchTag = matchTagsI.(string)
		}
	}
	return
}

// 获取奖励 icon url
func (tm *TimingMatchService) getRewardUrl() (rewardUrl string) {
	AwardIcon := tm.GetCurrConf(awardIconKey).(string)
	if AwardIcon != "" {
		awardIcons := map[string]string{}
		json.Unmarshal([]byte(AwardIcon), &awardIcons)
		if url, ok := awardIcons["1"]; ok {
			rewardUrl = tm.getIconUrl(url)
		}
	}
	return
}

// 循环信息
type loopIntervalInfo struct {
	LoopType         uint32
	LoopendTime      string
	LoopInterval     float32
	LoopIntervalTime uint64
	FirstBeginTime   uint64
}

// 加载循环周期时间和循环赛
func (tm *TimingMatchService) getLoopInterval() (info loopIntervalInfo) {
	// 比赛循环周期 1:单场 2:多场
	loopType, ok := tm.matchItem["looptype"]
	if ok {
		info.LoopType = uint32(StrToInt(loopType.(string)))
	}
	// 间隔时间分钟数
	loopIntervalMinute := StrToInt(tm.matchItem["loopinterval"].(string))
	// 间隔时间秒数
	loopIntervalSeconds := 0
	if v, ok := tm.matchItem["loopintervalsecond"]; ok {
		loopIntervalSeconds = StrToInt(v.(string))
	}
	// 组合间隔总时间秒
	cycleTimeIntervalSecond := loopIntervalMinute*60 + loopIntervalSeconds
	// 组合间隔总时间分钟
	cycleTimeIntervalMinute := fmt.Sprintf("%.1f", float64(cycleTimeIntervalSecond/60))
	// 转换为 float32
	loopInterval, err := strconv.ParseFloat(cycleTimeIntervalMinute, 32)
	if err != nil {
		logs.Error(tm.TAG+"loopInterval 转换失败: err: %v", err)
	}

	// 加载
	info.LoopInterval = float32(loopInterval)
	info.LoopIntervalTime = uint64(cycleTimeIntervalSecond)
	// 默认24小时循环赛
	if cycleTimeIntervalSecond != 0 && loopType.(string) == "" {
		info.LoopType = entity.MatchType_24Hours
	}
	// 第一次开始时间
	firstBeginTime := tm.matchItem["firstBeginTime"].(int64)
	info.FirstBeginTime = uint64(firstBeginTime - int64(cycleTimeIntervalSecond))
	if firstBeginTime <= time.Now().Unix() {
		info.FirstBeginTime = uint64(firstBeginTime)
	}
	// 周循环赛
	if StrToInt(loopType.(string)) == entity.MatchType_Week {
		if tm.matchItem["begintime"].(int64) > time.Now().Unix()+int64(cycleTimeIntervalSecond) {
			info.FirstBeginTime = uint64(tm.matchItem["begintime"].(int64) - int64(cycleTimeIntervalSecond))
		}
	}
	// 循环结束时间
	// 去掉秒
	loopEndTime, ok := tm.matchItem["loopendtime"]
	if ok {
		parts := strings.Split(loopEndTime.(string), ":")
		if len(parts) > 1 && parts[2] != "" {
			parts = parts[:2]
		}
		tm.matchItem["loopendtime"] = strings.Join(parts, ":")
	}
	info.LoopendTime = tm.matchItem["loopendtime"].(string)
	return
}

// 检测玩家是否报名
func (tm *TimingMatchService) checkPlayerSignUp(matchId uint64) bool {
	matchUserRd, _ := NewRedisCache(entity.RedisMatchUser)
	key := matchUserListPrefix + Uint64ToStr(matchId)
	field := IntToStr(tm.OnlineInfo.Mid)
	status := matchUserRd.HGet(key, field)
	return status != nil
}

// 定时赛过滤
func (tm *TimingMatchService) filterMatchInfo(matchItem map[string]interface{}) (filtered bool) {
	direction, ok := matchItem["ScreenDirection"]
	if !ok {
		logs.Error(tm.TAG + "MatchItem not have key[ScreenDirection]!")
		return
	}
	// 横竖屏过滤
	dir := direction.(int)
	if (tm.Vertical && dir == 0) || (!tm.Vertical && dir == 2) {
		return true
	}
	// 最低显示版本过滤
	minVersion, ok := matchItem["minversion"]
	if ok && tm.Version < minVersion.(int) {
		return true
	}
	// Android iOS 过滤
	if ClientIsAndroid(tm.App) {
		minVersionAndroid, ok := matchItem["minversionand"]
		if ok && (tm.Version < minVersionAndroid.(int)) {
			return true
		}
		maxVersionAndroid, ok := matchItem["maxversionand"]
		if ok && maxVersionAndroid.(int) > 0 && (tm.Version > maxVersionAndroid.(int)) {
			return true
		}
	} else if ClientIsIos(tm.App) {
		minVersionIos, ok := matchItem["minversionios"]
		if ok && (tm.Version < minVersionIos.(int)) {
			return true
		}
		maxVersionIos, ok := matchItem["maxversionios"]
		if ok && maxVersionIos.(int) > 0 && (tm.Version > maxVersionIos.(int)) {
			return true
		}
	}

	// 时间过滤
	displayTime, ok := matchItem["displaytime"]
	if !ok {
		return
	}
	beginTime, ok := matchItem["begintime"]
	if !ok {
		return
	}
	now := time.Now().Unix()
	if now < displayTime.(int64) || now > beginTime.(int64) {
		return true
	}
	return
}

// 更新报名人数
func (tm *TimingMatchService) adjustEnrollment(realNum uint64, isAddRobot bool) (enrollment uint64) {
	lestStarters := tm.GetAllEnrollment()
	enrollment = realNum
	area := beego.AppConfig.String("area")
	if isAddRobot && area == "lsqp" {
		random := RandInt(10, 15)
		enrollment = uint64(math.Floor(float64(lestStarters * random)))
	}
	return
}

// 获取分场数量
func (tm *TimingMatchService) getSubFieldNum(applyNum uint64, maxUserCount, maxSubFieldNum int) (subFieldNum int) {
	partitionNum := 1
	if float64(applyNum) > float64(maxUserCount)*1.5 {
		partitionNum = int(math.Floor(float64(applyNum / uint64(maxUserCount))))
	}
	subFieldNum = partitionNum
	if partitionNum > maxSubFieldNum {
		subFieldNum = maxSubFieldNum
	}
	return
}
