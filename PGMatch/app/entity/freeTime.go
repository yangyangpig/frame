package entity

// 免费次数服务接口
type FreeTimer interface {
	// 获取用户免费总次数
	GetTotalTimes(macthConfigId int) int
    // 获取用户已经使用的免费次数
	GetUsedTimes(matchConfigId, mid int) int
	// 获取用户剩余的免费次数
	GetRemainingTimes(matchConfigId, mid int) int
}
