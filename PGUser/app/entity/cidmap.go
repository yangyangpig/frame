package entity

type CidMap struct {
	Cid int64 `json:"cid" orm:"column(cid);pk"`
	PlatformId string `json:"platform_id" orm:"column(platform_id)"`
	PlatformType int32 `json:"platform_type" orm:"column(platform_type)"`
}