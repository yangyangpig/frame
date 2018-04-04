package entity

type Notice struct {
	NoticeId	int32	`json:"notice_id" orm:"column(notice_id);pk"`
	NoticeType	int32	`json:"notice_type" orm:"column(notice_type)"`
	AppId 		int		`json:"app_id" orm:"column(app_id)"`
	Weight		int32	`json:"weight" orm:"column(weight)"`
	Title		string	`json:"title" orm:"column(title)"`
	Content		string 	`json:"content" orm:"column(content)"`
	StartTime	string	`json:"start_time" orm:"column(start_time)"`
	EndTime		string	`json:"end_time" orm:"column(end_time)"`
	Conditions	string	`json:"conditions" orm:"column(conditions)"`
	Mids		string	`json:"mids" orm:"column(mids)"`
	Status 		int8	`json:"status" orm:"column(status)"`
	UpdateTime	int		`json:"update_time" orm:"column(update_time)"`
}
