package service

import (
	"github.com/astaxie/beego/logs"
	"time"
	"strings"
	"PGMatch/app/proto"
	"PGMatch/app/entity"
)

// 投注赛服务
type BettingMatchService struct {
	*BaseMatch
}

const (
	// 投注赛比赛类型
	betMatch        = 8
	betMatchListKey = "BET_MATCH_KEY_LIST"    // 投注赛key
)

// 分段结束时间
var mSegmentsEndTime *int64

// 初始化投注赛
func NewBettingMatchService(baseMatch *BaseMatch) entity.MatchLister {

	return &BettingMatchService{baseMatch}
}

// 初始化
func (m *BettingMatchService) Init() {
	// 初始化比赛列表信息
	matchListI := MatchConfService.GetAllConfBy(betMatchListKey)
	matchList := matchListI.(map[string]interface{})
	m.MatchList = matchList
}

// 投注赛列表过滤
func (m *BettingMatchService) Filter() {
	matchList := m.MatchList
	if matchList == nil {
		return
	}
	for matchId, matchGameIdI := range matchList {
		matchGameIdByte, ok := matchGameIdI.([]byte)
		if !ok {
			logs.Debug("[BET] matchGameId type is %v。", matchGameIdByte)
		}
		matchGameId := ByteToStr(matchGameIdByte)
		if matchGameId != "" && m.GameId != StrToUint64(matchGameId) {
			continue
		}

		matchHash := MatchConfService.GetConfBy(matchInfoKey, "match"+matchId)
		if matchHash == nil {
			continue
		}
		m.LoadCurrentMatch(StrToInt64(matchId), matchHash)

		// 地区限制
		if !m.FilterRegion() {
			continue
		}

		// 横竖屏限制
		if !m.FilterScreenDirection() {
			continue
		}

		// iOS 提审过滤
		if !m.FilterIOSAudit() {
			continue
		}

		// 白名单过滤
		if !m.FilterWhiteList() {
			continue
		}

		// 日期过滤
		if !m.FilterMatchData() {
			continue
		}

		// 分段时间过滤
		if !m.filterSegment() {
			continue
		}

		// 获取不同机型版本
		if !m.FilterVersion() {
			continue
		}

		// 加载 Response
		m.SetResponse()
	}
}

// 加载到 response
func (m *BettingMatchService) SetResponse() {
	//len(m.curr.bean)
	matchInfo := proto.List{}
	matchInfo.Type = betMatch
	matchInfo.Id = uint64(m.curr.matchId)
	matchInfo.Configid = m.GetCurrConf(matchConfigIdKey).(string)
	// 比赛名称
	matchName := m.GetCurrConf(nameKey).(string)
	matchInfo.Mname = matchName
	matchShortNameI := m.GetCurrConf(matchShortNameKey)
	matchInfo.Sname = matchName
	matchShortName, ok := matchShortNameI.(string)
	if ok && matchShortName != "" {
		matchInfo.Sname = matchShortName
	}
	// game id
	matchInfo.Gameid = uint64(m.GameId)
	matchInfo.Gamesort = uint64(m.GameId)
	if gameWeight, ok := m.GameWeight[Uint64ToStr(m.GameId)]; ok {
		matchInfo.Gamesort = StrToUint64(gameWeight.(string))
	}
	if gameName, ok := m.Games[Uint64ToStr(m.GameId)]; ok {
		matchInfo.Gamename = gameName.(string)
	}
	matchIcon := m.GetCurrConf(matchIconKey)
	if matchIcon != nil && matchIcon != ""{
		matchInfo.Micon = m.getIconUrl(matchIcon.(string))
	}
	adIcon := m.GetCurrConf(adIconKey)
	if adIcon != nil && adIcon != "" {
		matchInfo.AdIcon = adIcon.(string)
	}
	listSort := m.GetCurrConf(listSortKey)
	if listSort != nil && listSort != "" {
		matchInfo.ListSort = listSort.(string)
	}
	if *mSegmentsEndTime > 0 {
		matchInfo.Etime = uint64(*mSegmentsEndTime)
	}
	// 最低开赛人数
	totalNum := m.GetCurrConf(totalNumKey)
	matchUserCount := m.GetCurrConf(matchUserCountKey)
	if totalNum != nil {
		matchInfo.Requestnum = totalNum.(string)
		if matchUserCount != nil {
			matchInfo.Requestnum = matchUserCount.(string)
		}
	}
	// 运营描述
	rewardDescribe := m.GetCurrConf(rewardDescribeKey)
	if rewardDescribe != nil {
		matchInfo.RewardDescribe = rewardDescribe.(string)
	}

	// 比赛标签和 icon
	matchInfo.Iconweight = m.GetIconWeight()
	matchInfo.AdIcon = m.GetAdvertIcon()
	matchInfo.Matchtags = m.GetMatchTags()

	// 加载报名门槛 threshold(阀门)
	if Threshold, ok := m.GetCurrConf(thresholdKey).(map[string]interface{}); ok {
		ThresholdType := StrToFloat64(Threshold["type"].(string))
		if ThresholdType != 0 {
			matchInfo.ThresholdType = ThresholdType
		}
		// 报名门槛
		ThresholdValue := StrToFloat64(Threshold["value"].(string))
		if ThresholdValue != 0 {
			matchInfo.Threshold = ThresholdValue
		}
	}
	// 加载报名费
	m.loadFee(matchInfo.Fee)
	// 最大奖池
	tableUserCount := m.GetCurrConf(tableUserCountKey)
	thresholdVal := matchInfo.Threshold
	if tableUserCount != nil {
		tableUserCountNum := StrToUint32(tableUserCount.(string))
		matchInfo.Maxawardpool = tableUserCountNum * uint32(thresholdVal)
	}
	// 比赛标签
	matchLabel := m.GetCurrConf(matchLabelKey)
	if matchLabel != nil {
		matchInfo.Mclabel = map[string]string{}
		if label, ok := matchLabel.(map[string]interface{}); ok {
			for k, v := range label {
				matchInfo.Mclabel[k] = v.(string)
			}
		}
	}
	// 添加到比赛列表
	m.AddToResponseList(&matchInfo)
}

// 加载报名费
func (m *BettingMatchService) loadFee(mFee []proto.Fee) {
	feeI := m.GetCurrConf(feeKey)
	fee, ok := feeI.([]interface{})
	if !ok {
		logs.Error("fee's type is not []interface{}")
		return
	}
	if len(fee) <= 0 {
		logs.Error("fee's length is <= 0")
		return
	}
	mFee = make([]proto.Fee, len(fee))
	for i, feeI := range fee {
		mFee[i] = proto.Fee{}
		feeValueSlice, ok := feeI.([]interface{})
		if !ok {
			logs.Error("feeI's type is not []interface{}")
			return
		}
		subFees := make([]proto.SubFee, len(feeValueSlice))
		for subIndex, feeValue := range feeValueSlice {
			feeMap, ok := feeValue.(map[string]interface{})
			if !ok {
				logs.Error("feeValue's type is not map[string]interface{}")
				return
			}
			var tp uint32
			var num uint64
			var desc string
			if typeValue, ok := feeMap["type"].(string); ok {
				tp = StrToUint32(typeValue)
			}
			if numValue, ok := feeMap["num"].(string); ok {
				num = StrToUint64(numValue)
			}
			if descValue, ok := feeMap["desc"].(string); ok {
				desc = descValue
			}
			subFees[subIndex] = proto.SubFee{
				Type: tp,
				Num:  num,
				Desc: desc,
			}
		}
		mFee[i].Subfee = subFees
	}
	return
}

// 分段时间过滤
func (m *BettingMatchService) filterSegment() bool {
	var endTime int64 = 0
	mSegmentsEndTime = &endTime
	matchTypeI := m.GetCurrConf(matchTypeKey)
	// 新后台时间段判断
	matchType := StrToInt(matchTypeI.(string))
	// 周循环赛或间隔循环赛
	if matchType == entity.MatchType_Week || matchType == entity.MatchType_Interval {
		perDayTimeSegments := m.GetCurrConf(perDayTimeSegmentsKey)
		timeSeg, ok := perDayTimeSegments.(int)
		if !ok {
			logs.Debug("perDayTimeSegmentsKey[%T] is not int!", perDayTimeSegments)
			return true
		}
		//确保数据有"段"的信息
		if timeSeg <= 0 {
			return true
		}
		now := time.Now()
		HMS := IntToStr(now.Hour()) + IntToStr(now.Minute()) + IntToStr(now.Second())
		for i := 1; i < timeSeg; i++ {
			segmentsFromTimeI := m.GetCurrConf(segmentsTimeStart + IntToStr(i))
			segmentsToTimeI := m.GetCurrConf(segmentsTimeEnd + IntToStr(i))
			startHMS := strings.Replace(":", "", segmentsFromTimeI.(string), 10)
			endHMS := strings.Replace(":", "", segmentsToTimeI.(string), 10)
			if HMS < startHMS && HMS >= endHMS {
				*mSegmentsEndTime = FormatTime(segmentsToTimeI.(string)).Unix()
				return false
			}
		}
	}
	return true
}
