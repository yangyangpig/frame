package entity

type MatchLister interface {
	// 初始化
	Init()
	// 比赛过滤(如:超时)
	Filter()
	// 组合比赛信息到 Response
	SetResponse()
}
