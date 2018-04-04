package service

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"putil/log"
)

//新版本手机短信通知
type boyaaSmsService struct {}

//获取语音验证码
func (bs *boyaaSmsService) Get(phone string, captcha string, smsType string) bool {
	content := make(map[string]interface{})
	content["content"] = captcha
	if smsType == "voice" {
		apiUrl := ServiceConf["sms"].String(Runmode+"::noticeVoice.api")
		authId, _ := ServiceConf["sms"].Int("noticeVoice.authId")
		token := ServiceConf["sms"].String("noticeVoice.token")
		tplId, _ := ServiceConf["sms"].Int("noticeVoice.tplId")
		ret := bs.curl(apiUrl, authId, token, phone, tplId, content)
		return ret
	} else {
		return false
	}
}
//GET请求
func (bs *boyaaSmsService) curl(api string, authId int, authKey string, phone string, tplId int, params map[string]interface{}) bool {
	content, err := json.Marshal(params)
	if err != nil {
		return false
	}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		plog.Debug("boyaasms new request====", err)
		return false
	}
	q := req.URL.Query()
	q.Add("auth_id", fmt.Sprintf("%d", authId))
	q.Add("auth_token", authKey)
	q.Add("phone", phone)
	q.Add("tpl_id", fmt.Sprintf("%d", tplId))
	q.Add("zone_id", "86")
	q.Add("params", string(content))
	q.Add("p", "1")
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		plog.Debug("boyaasms do====", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		plog.Debug("boyaasms readAll====", err)
		return false
	}
	g := gjson.Parse(string(body))
	plog.Debug("boyaasms curl====", g)
	exist := g.Get("code").Exists()
	if exist {
		code := g.Get("code").Int()
		if code == 0 {
			return true
		}
	}
	return false
}