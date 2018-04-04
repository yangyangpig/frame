package controller

import (
	"dfqp/proto/login"
	"dfqp/pg-login/logic"
	"dfqp/lib"
)

//登录注册服务
type BindController struct {}

//游客绑定微信
func (this *BindController) GuestBindWechat(request *pgLogin.GuestBindWechatRequest) *pgLogin.GuestBindWechatResponse {
	resp := new(pgLogin.GuestBindWechatResponse)
	if request.Mid <= 0 || len(request.OpenId) == 0 || len(request.AccessToken) == 0 {
		resp.Status = 1001
		return resp
	}
	bindResp := logic.WechatRegLogic.GuestBindWechat(request)
	resp.Status = bindResp.ErrCode
	resp.Data.LoginType = bindResp.LoginType
	return resp
}

//游客绑定手机
func (this *BindController) GuestBindPhone(request *pgLogin.GuestBindPhoneRequest) *pgLogin.GuestBindPhoneResponse {
	resp := new(pgLogin.GuestBindPhoneResponse)
	if request.Mid <= 0 || len(request.Captcha) == 0 {
		resp.Status = 1001
		return resp
	}
	if len(request.Phone) == 0 {
		resp.Status = 1015
		return resp
	}
	bindResp := logic.PhoneRegLogic.GuestBindPhone(request)
	resp.Status = bindResp.ErrCode
	resp.Data.LoginType = bindResp.LoginType
	resp.Data.Pwd = bindResp.Pwd
	resp.Data.Bid = bindResp.Bid
	resp.Data.AccessToken = bindResp.AccessToken
	return resp
}

//微信绑定手机
func (this *BindController) WechatBindPhone(request *pgLogin.WechatBindPhoneRequest) *pgLogin.WechatBindPhoneResponse {
	resp := new(pgLogin.WechatBindPhoneResponse)
	if request.Mid <= 0 || len(request.Captcha) == 0 {
		resp.Status = 1001
		return resp
	}
	if !lib.IsTelephone(request.Phone, "CHN") {
		resp.Status = 1015
		return resp
	}
	bindResp := logic.WechatRegLogic.WechatBindPhone(request)
	resp.Status = bindResp.ErrCode
	resp.Data.Pwd = bindResp.Pwd
	return resp
}
