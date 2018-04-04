package entity

// 支付宝调起支付必要参数
type AliPayParam struct {
	Porder     string `json:"porder"`
	Pname      string `json:"pname"`
	Pamount    string `json:"pamount"`
	Udesc      string `json:"udesc"`
	NotifyUrl  string `json:"notify_url"`
	Partner    string `json:"PARTNER"`
	Seller     string `json:"SELLER"`
	RsaPrivate string `json:"RSA_PRIVATE"`
}

// 微信调起支付必要参数
type WeChatParam struct {
	PartnerId    string `json:"partnerId"`
	PrepayId     string `json:"prepayId"`
	NonceStr     string `json:"nonceStr"`
	TimeStamp    string `json:"timeStamp"`
	PackageValue string `json:"packageValue"`
	Sign         string `json:"sign"`
	ExtData      string `json:"extData"`
	Pmode        string `json:"pmode"`
}

// 银联调起支付必要参数
type UnionPayParam struct {
	Tn    string `json:"tn"`
	Pmode string `json:"pmode"`
}

// 苹果调起支付必要参数
type ApplePayParam struct {
	Order      string `json:"ORDER"`
	Appstoreid string `json:"appstoreid"`
}

// 发货请求回调
type SendMoneyParam struct {
	AppleProductid string `json:"apple_productid"`
	Do             string `json:"do"`
	Ext            string `json:"ext"`
	GameItemID     string `json:"game_item_id"`
	Mid            string `json:"mid"`
	Pamount        string `json:"pamount"`
	PamountChange  string `json:"pamount_change"`
	PamountRate    string `json:"pamount_rate"`
	PamountUnit    string `json:"pamount_unit"`
	PamountUsd     string `json:"pamount_usd"`
	PaychipsV2     string `json:"paychips_v2"`
	Paycoins       string `json:"paycoins"`
	Payconfid      string `json:"payconfid"`
	PayprodV2      string `json:"payprod_v2"`
	PcAppid        string `json:"pc_appid"`
	PcRate         string `json:"pc_rate"`
	PcSid          string `json:"pc_sid"`
	PcTime         string `json:"pc_time"`
	Pdealno        string `json:"pdealno"`
	Pendtime       string `json:"pendtime"`
	Pid            string `json:"pid"`
	Pmode          string `json:"pmode"`
	PnumV2         string `json:"pnum_v2"`
	Pstarttime     string `json:"pstarttime"`
	Sign           string `json:"sign"`
	SignV2         string `json:"sign_v2"`
	SignV3         string `json:"sign_v3"`
	Sitemid        string `json:"sitemid"`
	Time           string `json:"time"`
	RequestIp      string `json:"request_ip"`
}
