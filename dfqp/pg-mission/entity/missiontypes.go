package entity

type Missiontypes struct {
	Id           int64  `json:"id" orm:"column(id);pk"`
	Region_id    int32  `json:"region_id" orm:"column(region_id)"`
	Name         string `json:"name" orm:"column(name)"`
	Desc         string `json:"desc" orm:"column(desc)"`
	Icon         string `json:"icon" orm:"column(icon)"`
	Reward       int64  `json:"reward" orm:"column(reward)"`
	Reward_type  int64  `json:"reward_type" orm:"column(reward_type)"`
	Sort_order   int32  `json:"sort_order" orm:"column(sort_order)"`
	Conditions   string `json:"conditions" orm:"column(conditions)"`
	Jump_code    string `json:"jump_code" orm:"column(jump_code)"`
	Status       int32  `json:"status" orm:"column(status)"`
	Circle_type  int32  `json:"circle_type" orm:"column(circle_type)"`
	Maxsupported string `json:"maxsupported" orm:"column(maxsupported)"`
	Cities       string `json:"cities" orm:"column(cities)"`
}
