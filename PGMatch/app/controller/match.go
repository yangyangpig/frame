package controller

import (
	"PGMatch/app/proto"
	"PGMatch/app/service"
	"github.com/astaxie/beego/logs"
	"PGMatch/app/entity"
	"time"
	"math"
	"logic/common"
)

// 比赛 Controller
type MatchController struct {
	adapters map[string]entity.MatchLister
}

const (
	// 比赛列表每页显示数
	perPageNumber = 10
	// 最大显示 icon 数
	maxDisplayIconNum = 5
)

// 比赛列表
func (m *MatchController) Lists(rq *proto.ListsRequest) *proto.ListsResponse {
	mid := rq.GetMid()
	appId := rq.GetApp()
	regionId := rq.GetAreaId()
	gameId := rq.GetGameId()
	sessionId := rq.GetSsid()
	page := 1
	if rq.GetPage() > 1 {
		page = int(rq.GetPage())
	}
	filterMode := false
	if rq.GetPage() == 0 && len(rq.GetIds()) != 0 {
		filterMode = true
	}
	Response := new(proto.ListsResponse)
	Response.Srvtime = uint64(time.Now().Unix())
	Response.Tpage = uint64(page)
	Response.Pnum = perPageNumber
	baseMatch := service.NewBaseMatch(mid, appId, gameId, regionId, sessionId, Response)
	// 各个比赛列表分发获取
	m.dispatchGetMatchList(baseMatch)
	// 列表过滤
	if baseMatch.Response.List != nil && len(baseMatch.Response.List) > 0 {
		// 没有指定分页 && 上报自己的比赛 id 进行过滤
		if rq.GetPage() != 0 && rq.GetIds() != nil && len(rq.GetIds()) != 0 {
			mFilter := service.NewDynamicFilterService(rq.GetIds(), baseMatch.Response)
			mFilter.DynamicFilter()
		}
		// 加载其他信息
		startPage := 0
		var totalPages uint64 = 1
		Response.Filter = 1
		if !filterMode {
			listLength := len(Response.List)
			startPage = (page - 1) * perPageNumber
			totalPages = uint64(math.Ceil(float64(listLength))) / perPageNumber
			temp := map[uint32]proto.List{}
			for k, v := range Response.List {
				if int(k) > startPage && k < perPageNumber {
					temp[k] = v
				}
			}
			Response.List = temp
			Response.Filter = 0
		} else {
			Response.Pnum = uint32(len(Response.List))
		}
		Response.Tpage = uint64(totalPages)
		Response.Iconmax = maxDisplayIconNum
		// TODO CFG()
		Response.Sort = service.GenericConfService.GetMatchSortConf().(map[string]uint32)
	}
	return baseMatch.Response
}

const (
	betting       = iota // 投注赛
	fullStart            // 快速赛
	fullStart1650
	adversity
	group
)

// 不同比赛匹配不同条件
func (m *MatchController) dispatchGetMatchList(baseMatch *service.BaseMatch) {
	// 版本限制
	var versionLimit = map[int]int{
		betting:       1350,
		fullStart:     1650,
		fullStart1650: 1650,
		adversity:     1900,
		group:         900,
	}

	// 投注赛
	if baseMatch.Version >= versionLimit[betting] {
		betMatch := service.NewBettingMatchService(baseMatch)
		m.registeredMatch("投注赛", betMatch)
	}

	// 人满开赛
	if baseMatch.Version >= versionLimit[fullStart] {
		fullStartMatch := service.NewFullStartMatchService(baseMatch)
		m.registeredMatch("人满开赛", fullStartMatch)
	}

	// 逆袭赛
	if baseMatch.Version >= versionLimit[adversity] {
		adversityMatch := service.NewAdversityMatchService(baseMatch)
		m.registeredMatch("逆袭赛", adversityMatch)
	}

	// 定时赛
	timingMatch := service.NewTimingMatchService(baseMatch)
	m.registeredMatch("定时赛", timingMatch)

	// 快速赛
	fastMatch := service.NewFastMatchService(baseMatch)
	m.registeredMatch("快速赛", fastMatch)

	// 集团赛
	if baseMatch.Version < versionLimit[group] {
		groupMatch := service.NewGroupMatchService(baseMatch)
		m.registeredMatch("集团赛", groupMatch)
	}

	// 加载各比赛列表
	m.loadMatchList()
}

// 注册比赛
func (m *MatchController) registeredMatch(matchTypeName string, adapter entity.MatchLister) {
	if m.adapters == nil {
		m.adapters = map[string]entity.MatchLister{}
	}
	m.adapters[matchTypeName] = adapter
}

// 获取不同比赛列表
func (m *MatchController) loadMatchList() {
	matchMap := m.adapters
	for matchName, matchAdapter := range matchMap {
		logs.Debug("The Match matchType Is %v !", matchName)
		matchAdapter.Init()
		matchAdapter.Filter()
	}
}
