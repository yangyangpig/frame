package entity

type Config struct {
	AppId int32  `json:"app_id" orm:"column(app_id)"`
	ConfigKey string  `json:"config_key" orm:"column(config_key)"`
	ConfigValue string  `json:"config_value" orm:"column(config_value)"`
	Time int32 `json:"time" orm:"column(time)"`
}
