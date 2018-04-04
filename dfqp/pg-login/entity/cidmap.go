package entity

type Cidmap struct {
	Id int64 `json:"id" orm:"column(id);pk"`
	Cid int64 `json:"cid" orm:"column(cid);"`
	PlatformId string `json:"platform_id" orm:"column(platform_id)"`
	PlatformType int32 `json:"platform_type" orm:"column(platform_type)"`
}