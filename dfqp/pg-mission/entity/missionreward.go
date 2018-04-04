package entity

type Missionreward struct {
	Id          int64 `json:"id" orm:"column(id);pk"`
	Mid         int64 `json:"mid" orm:"column(mid)"`
	App_id      int32 `json:"app_id" orm:"column(app_id)"`
	Mission     int64 `json:"mission" orm:"column(mission)"`
	Reward_type string `json:"reward_type" orm:"column(reward_type)"`
	Status      int32 `json:"status" orm:"column(status)"`
	Time        int64 `json:"time" orm:"column(time)"`
}
