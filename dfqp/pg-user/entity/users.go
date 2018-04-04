package entity

type User struct {
	Cid int64 `redis:"cid" json:"cid" orm:"column(cid);pk"`
	Nick string `redis:"nick" json:"nick" orm:"column(nick);"`
	Sex int32 `redis:"sex" json:"sex" orm:"column(sex);"`
	Icon string `redis:"icon" json:"icon" orm:"column(icon);"`
	Icon_big string `redis:"icon_big" json:"icon_big" orm:"column(icon_big);"`
	City string `redis:"city" json:"city" orm:"column(city);"`
	Phone string `redis:"phone" json:"phone" orm:"column(phone);"`
	Sign string `redis:"sign" json:"sign" orm:"column(sign);"`
	Status int32 `redis:"status" json:"status" orm:"column(status);"`
	IconId string `redis:"icon_id" json:"icon_id" orm:"column(icon_id);"`
}