package entity

type Logininfo struct {
	Cid int64 `json:"cid" orm:"column(cid);pk"`
	FirstApp int32 `json:"first_app" orm:"column(first_app)"`
	LastApp int32 `json:"last_app" orm:"column(last_app)"`
	FirstVersion string `json:"first_version" orm:"column(first_version)"`
	LastVersion string `json:"last_version" orm:"column(last_version)"`
	RegTime int32 `json:"reg_time" orm:"column(reg_time)"`
	LoginTime int32 `json:"login_time" orm:"column(login_time)"`
	FirstIp string `json:"first_ip" orm:"column(first_ip)"`
	LastIp string `json:"last_ip" orm:"column(last_ip)"`
}
