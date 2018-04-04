package controller

import (
	"net/http"
	"encoding/json"
	"dfqp/proto/pay"
	"dfqp/pg-http/service"
)

type PayHttp struct {
	base
}

var (
	rqData map[string]string
	rqParam = []string{
		"apple_productid", "do", "ext", "game_item_id", "mid", "pamount",
		"pamount_change", "pamount_rate", "pamount_unit", "pamount_usd", "paychips_v2", "paycoins",
		"payconfid", "payprod_v2", "pc_appid", "pc_rate", "pc_sid", "pc_time",
		"pdealno", "pendtime", "pid", "pmode", "pnum_v2", "pstarttime",
		"sign", "sign_v2", "sign_v3", "sitemid", "time",
	}
	signFlag bool
)

type SendMoneyRequest struct {
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
	RequestIp	   string `json:"request_ip"`
}

func (this *PayHttp) SendMoney(writer http.ResponseWriter, rq *http.Request) {
	// 此处只做请求转发，预处理请求参数
	rq.ParseForm()
	rqData = make(map[string]string)
	for _, v := range rqParam {
		rqData[v] = rq.Form.Get(v)
	}
	rqData["request_ip"] = rq.Header.Get("Remote_addr")
	if rqData["request_ip"] == "" {
		rqData["request_ip"] = rq.RemoteAddr
	}
	if rqData["pid"] == "" {
		this.output(writer, 2000, "", "请求参数有误！")
	}
	rqJson,_ := json.Marshal(rqData)
	sendMoneyData := new(pgPay.SendMoneyRequest)
	sendMoneyData.Content = string(rqJson)
	reqBytes, _ := sendMoneyData.Marshal()
	response := service.Client.SendAndRecvRespRpcMsg("pgPay.sendMoney", reqBytes, 2000, 0)
	sendMoneyResponse := new(pgPay.SendMoneyResponse)
	sendMoneyResponse.Unmarshal(response.Body)
	if sendMoneyResponse.Status == 0 {
		this.output(writer, 1, "", "发货成功！")
	} else {
		this.output(writer, int(sendMoneyResponse.Status), "", sendMoneyResponse.Msg)
	}
}