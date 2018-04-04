// 游戏信息服务
package service

import "PGMatch/app/entity"

type gameInfoService struct {
	mid int
}

// 构造游戏信息服务
func NewGameInfoService(mid int) entity.GameInfoInterface {
	return &gameInfoService{mid:mid}
}

func (gs *gameInfoService) GetRecord() (record map[string]interface{}) {
	return
}

// 获取用户游戏信息
func (gs *gameInfoService) getGameInfo(isFirstLoginToday bool) (gameInfo map[string]interface{}) {
	return
}
