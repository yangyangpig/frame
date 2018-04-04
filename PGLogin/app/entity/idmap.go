package entity

type Idmap struct {
	Mid int64 `json:"mid" orm:"column(mid);pk"`
	PlatformId string `json:"platform_id" orm:"column(platform_id)"`
	PlatformType int32 `json:"platform_type" orm:"column(platform_type)"`
	Region int32 `json:"region" orm:"column(region)"`
}