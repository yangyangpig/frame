package entity

type RegFail struct {
	Cid int64 `redis:"cid" json:"cid" orm:"column(cid);pk"`
	Step int32 `redis:"step" json:"step" orm:"column(step);"`
}