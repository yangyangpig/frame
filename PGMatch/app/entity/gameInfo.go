package entity

// 游戏信息接口
type GameInfoInterface interface {
	// 获取一条游戏记录
	GetRecord() (record map[string]interface{})
}
