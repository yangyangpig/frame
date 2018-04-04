package entity

// 比赛配置服务接口
type MatchConfer interface {
	// 获取所有配置
	GetAllConfBy(key string) (resp interface{})
	// 获取单独配置项
	GetConfBy(key, filed string) (resp interface{})
}
