package entity

type Cuser struct {
	Cid int64 `json:"cid" orm:"column(cid);"`
	Nick string `json:"nick" orm:"column(nick);"`
	Sex int32 `json:"sex" orm:"column(sex);"`
	Icon string `json:"icon" orm:"column(icon);"`
	Icon_big string `json:"icon_big" orm:"column(icon_big);"`
	Hometown string `json:"hometown" orm:"column(hometown);"`
	City string `json:"city" orm:"column(city);"`
	Email string `json:"email" orm:"column(email);"`
	Phone string `json:"phone" orm:"column(phone);"`
	Realname string `json:"realname" orm:"column(realname);"`
	Idcard string `json:"idcard" orm:"column(idcard);"`
	Is_set int32 `json:"is_set" orm:"column(is_set);"`
}
