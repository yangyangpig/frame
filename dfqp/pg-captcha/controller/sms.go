package controller

import (
	"dfqp/proto/sms"
	"github.com/tidwall/gjson"
	"dfqp/pg-captcha/service"
	"putil/log"
	"fmt"
)

type SmsController struct {}

//获取短信验证码
func (this *SmsController) GetCaptcha(req *pgSms.GetCaptchaRequest) *pgSms.GetCaptchaResponse {
	resp := new(pgSms.GetCaptchaResponse)
	if len(req.Phone) <= 0 {
		resp.Status = 1 //失败
		return resp
	}
	ret := this.getUserToken(req.Phone, req.Type)
	if ret {
		resp.Status = 0 //成功
	} else {
		resp.Status = 1 //失败
	}
	return resp
}

//获取语音验证码
func (this *SmsController) GetVoiceCaptcha(req *pgSms.GetVoiceCaptchaRequest) *pgSms.GetVoiceCaptchaResponse {
	resp := new(pgSms.GetVoiceCaptchaResponse)
	phone := req.Phone
	flag := req.Type
	if len(phone) <= 0 {
		resp.Status = 1 //失败
		return resp
	}
	token := this.getToken(phone)
	if token == 0 {
		result := this.getUserToken(phone, flag)
		if result {
			token = this.getToken(phone)
		}
	}
	ret := service.BoyaaSmsService.Get(req.Phone, fmt.Sprintf("%v", token), "voice")
	if ret {
		resp.Status = 0
	} else {
		resp.Status = 1
	}
	return resp
}

//获取验证码
func (this *SmsController) getToken(phone string) int64 {
	params := make(map[string]interface{})
	params["phone"] = phone
	result, err := service.ByClientService.Get("user/gettoken", params)
	if err == nil {
		plog.Debug("getCaptcha====", result)
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		exist := jsonData.Get("result.token").Exists()
		if !exist || code != 200 {
			return 0
		} else {
			token := jsonData.Get("result.token").Int()
			plog.Debug("getToken====", token)
			return token
		}
	} else {
		return 0
	}
}

//发送动态验证码
func (this *SmsController) getUserToken(phone string, flag int32) bool  {
	params := make(map[string]interface{})
	params["phone"] = phone
	params["flag"] = flag
	result, err := service.ByClientService.Get("user/token", params)
	if err == nil {
		plog.Debug("getCaptcha====", result)
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		result := jsonData.Get("result").Int()
		if code != 200 || result != 100 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}