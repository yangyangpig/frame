package service

import (
	"github.com/astaxie/beego/logs"
	"logic/common"
	"logic/Match"
	"PGMatch/app/entity"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"encoding/json"
)

// 通用配置服务
type genericConfService struct {
	rd entity.RedisClient
}

// Http 返回结构体
type result struct {
	code   uint32                 `json:"code"`
	error  string                 `json:error`
	result map[string]interface{} `json:result`
}

const (
	confGames      = "games"
	confGameWeight = "gameweight"
	confApps       = "apps"
	confControl    = "control"

	publicMatchSort = "matchsort"
	// 接口密钥
	SECRET = "d!f@@qp()20170630"
	// 获取配置的 url
	HttpConfUrl = "http://192.168.200.21/dfqp/index.php"
)

// 初始化比赛配置服务
func NewGenericConfService() entity.GenericConfer {
	return &genericConfService{}
}

// 拉取游戏配置
func (s *genericConfService) GetGamesConf() (gamesConf map[string]interface{}) {
	params := entity.GenericConfParam{Api: -1}
	gamesConfI := s.getConfigInfo(confGames, params)
	logs.Debug("GetGamesConf() Param:%v", params)
	gamesConf, ok := gamesConfI.(map[string]interface{})
	if ok {
		logs.Error("GenericConfService.GetGamesConf() Fail! gamesConf'type is not map!\n")
		logs.Error("gamesConf : %v", gamesConfI)
		return
	}
	return map[string]interface{}{}
}

// 拉取游戏权重配置
func (s *genericConfService) GetGameWeightConf(params entity.GenericConfParam) (gameWeightConf map[string]interface{}) {
	gameWeightConfI := s.getConfigInfo(confGameWeight, params)
	logs.Debug("GetGameWeightConf() Param:%v", params)
	gameWeightConf, ok := gameWeightConfI.(map[string]interface{})
	if ok {
		logs.Error("GenericConfService.GetGameWeightConf() Fail! gameWeightConf'type is not map!\n")
		logs.Error("gameWeightConf : %v", gameWeightConfI)
		return
	}
	return map[string]interface{}{}
}

// 获取 app 配置
func (s *genericConfService) GetAppsConf(params entity.GenericConfParam) (appsConf map[string]interface{}) {
	appsConfI := s.getConfigInfo(confApps, params)
	logs.Debug("GetAppsConf() Param:%v", params)
	appsConf, ok := appsConfI.(map[string]interface{})
	if ok {
		logs.Error("GenericConfService.GetAppsConf() Fail! appsConf'type is not map!\n")
		logs.Error("appsConf : %v", appsConfI)
		return
	}
	return map[string]interface{}{}
}

// 获取开关配置
func (s *genericConfService) GetControlConf(params entity.GenericConfParam) (controlConf interface{}) {
	controlConf = s.getConfigInfo(confControl, params)
	logs.Debug("GetControlConf() Param:%v", params)
	return
}

// 获取比赛排序配置 <=> CFG("matchsort", "")
func (s *genericConfService) GetMatchSortConf() (matchSortConf interface{}) {
	url := HttpConfUrl
	url += "?"
	params := s.getHttpParam(publicMatchSort, "")
	for k, v := range params {
		url += k + "=" + v + "&"
	}
	url += "servertype=lua"
	resp, _ := httplib.Get(url).Response()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var result result
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		panic(err)
	}
	return result.result
}

// 获取 HTTP 请求的 params
func (s *genericConfService) getHttpParam(section, key string) map[string]string {
	return map[string]string{
		"action":  "externals.cfg4Lua",
		"section": section,
		"key":     key,
		"secret":  SECRET,
	}
}

// 拉取配置信息
func (s *genericConfService) getConfigInfo(confKey string, confParams entity.GenericConfParam) (resp interface{}) {
	if confKey == "" {
		return nil
	}
	var isLastDoc = confParams.IsLastDoc
	switch confKey {
	case confGames: // 获取游戏信息
		gameId := confParams.GameId
		gamesI := s.getConf(confGames, 0, isLastDoc)
		games, ok := gamesI.(map[string]interface{})
		if !ok {
			logs.Error("gamesI type is not map[string]interface{}")
			return
		}
		if len(games) > 0 {
			if gameId != 0 {
				if game, ok := games[IntToStr(gameId)]; ok {
					resp = game.(map[string]interface{})
				}
			} else {
				resp = games
			}
		} else {
			// 暂时使用老的 CFG (logic)
			if gameId != 0 {
				param := confGames + "." + IntToStr(gameId)
				// TODO CFG("system", key)
				resp = common.CFG("system", param)
			} else {
				// TODO CFG("system", "games")
				resp = common.CFG("system", "games")
			}
		}
	default:
		if confParams.Api != -1 {
			resp = s.getConf(confKey, confParams.Api, isLastDoc)
		}
	}
	return
}

// 获取配置文件
func (s *genericConfService) getConf(confKey string, api int, isLastDoc bool) (resp interface{}) {
	// 暂时使用老的配置文件获取函数
	resp = match.NewEquip().GetGenericConfig(confKey, api, isLastDoc)
	return
}
