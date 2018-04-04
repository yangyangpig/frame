package entity

// 公共配置参数结构体
type GenericConfParam struct {
	Api       int
	IsLastDoc bool // 是否最后一次
	GameId    int
	AppId     int
}

// 公共配置获取接口
type GenericConfer interface {
	// 获取游戏配置
	GetGamesConf() (gamesConf map[string]interface{})
	// 获取游戏权重配置
	GetGameWeightConf(params GenericConfParam) (gameWeightConf map[string]interface{})
	// 获取应用配置
	GetAppsConf(params GenericConfParam) (appsConf map[string]interface{})
	// 获取开关配置 control
	GetControlConf(params GenericConfParam) (controlConf interface{})

	// HTTP
	// 获取 比赛排序配置
	GetMatchSortConf() (matchSortConf interface{})
}
