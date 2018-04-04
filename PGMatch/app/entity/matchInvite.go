package entity

// 比赛匹配服务接口
type MatchInviter interface {
	ApplyNum(key string) int
}
