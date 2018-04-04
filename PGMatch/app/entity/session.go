package entity

// 用户线上信息
type OnlineInfo struct {
	Active_time   int    `json:"active_time"`
	App_id        int    `json:"app_id"`
	Broken_cd     int    `json:"broken_cd"`
	Channel_id    int    `json:"channel_id"`
	Game_id       int    `json:"game_id"`
	Hall_version  int    `json:"hall_version"`
	Ip            string `json:"ip"`
	Ispay         int    `json:"ispay"`
	Login_time    int    `json:"login_time"`
	Mid           int    `json:"mid"`
	Partner_info  int    `json:"partner_info"`
	Platformtype  int    `json:"platformtype"`
	Play_end      int    `json:"play_end"`
	Play_seconds  int    `json:"play_seconds"`
	Play_start    int    `json:"play_start"`
	Reg_time      int    `json:"reg_time"`
	Server_id     int    `json:"server_id"`
	Sess_id       string `json:"sess_id"`
	Table_id      int    `json:"table_id"`
	Table_lv      int    `json:"table_lv"`
	Version       string `json:"version"`
	ProvinceAppId int    `json:"provinceAppId"`
}

// sessionService 暴露接口
type Sessioner interface {
	// 获取用户线上信息
	GetOnlineInfo(id interface{}) *OnlineInfo
}
