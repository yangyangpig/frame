package entity

type Platform2cid struct {
	Id int64 `json:"id" orm:"column(cid);pk"`
	PlatformId string `json:"platform_id" orm:"column(platform_id)"`
	PlatformType int32 `json:"platform_type" orm:"column(platform_type)"`
	Cid int64 `json:"cid" orm:"column(cid);"`
}