package entity

type User struct {
	Mid int64 `json:"mid" orm:"column(mid);pk"`
	App_id int32 `json:"app_id" orm:"column(app_id);"`
	Reg_time int64 `json:"reg_time" orm:"column(reg_time);"`
	Login_time int64 `json:"login_time" orm:"column(login_time);"`
	Status int32 `json:"status" orm:"column(status);"`
	Vip_time int64 `json:"vip_time" orm:"column(vip_time);"`
	Vip_level int32 `json:"vip_level" orm:"column(vip_level);"`
	Platform_id string `json:"platform_id" orm:"column(platform_id);"`
	Platform_type int32 `json:"platform_type" orm:"column(platform_type);"`
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
	Org_app int32 `json:"org_app" orm:"column(org_app);"`
	Org_channel int32 `json:"org_channel" orm:"column(org_channel);"`
	Channel_id int32 `json:"channel_id" orm:"column(channel_id);"`
	Timeout int64 `json:"timeout" orm:"column(timeout);"`
	Bagvol int64 `json:"bagvol" orm:"column(bagvol);"`
	Ispay int32 `json:"ispay" orm:"column(ispay);"`
	First_match int32 `json:"first_match" orm:"column(first_match);"`
	Fast_match int32 `json:"fast_match" orm:"column(fast_match);"`
	Version string `json:"version" orm:"column(version);"`
	Is_set int32 `json:"is_set" orm:"column(is_set);"`
	Partner_info int32 `json:"partner_info" orm:"column(partner_info);"`
}