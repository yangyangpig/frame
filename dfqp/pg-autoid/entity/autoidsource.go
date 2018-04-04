package entity

type AutoidSource struct {
	Btag       string `json:"btag" orm:"column(btag);pk"`
	Maxid      int64  `json:"max_id" orm:"column(max_id)"`
	Step       int32  `json:"step" orm:"column(step)"`
	Des        string `json:"des" orm:"column(des)"`
	Updatetime int32  `json:"update_time" orm:"column(update_time)"`
}
