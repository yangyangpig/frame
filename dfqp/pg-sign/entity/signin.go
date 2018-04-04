package entity

type Signin struct {
	Mid int64 `json:"mid" orm:"column(mid)"`
	Month int32 `json:"month" orm:"column(month)"`
	SignValue string `json:"sign_value" orm:"column(sign_value)"`
}
