package controller

import (
	"dfqp/proto/login"
	"dfqp/pg-login/service"
	"github.com/tidwall/gjson"
	"putil/log"
	"dfqp/lib"
)

//修改密码
type PasswordController struct {}

//重置密码
func (this *PasswordController) Reset(request *pgLogin.ResetPwdRequest) *pgLogin.ResetPwdResponse {
	resp := new(pgLogin.ResetPwdResponse)
	resp.Status = 0
	pwd := request.Pwd
	phone := request.Phone
	if !lib.IsTelephone(phone, "CHN") {
		resp.Status = 1015
		return resp
	}
	captcha := request.Captcha
	if len(pwd) < 6 || len(pwd) > 12 {
		resp.Status = 1017
		return resp
	}
	params := make(map[string]interface{})
	params["type"] = "PHONE"
	params["phone"] = phone
	params["pwd"] = pwd
	params["token"] = captcha
	result, err := service.ByClientService.Get("user/resetpwd", params)
	if err != nil {
		resp.Status = 1014
		return resp
	} else {
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		plog.Debug("code-======", code)
		if code != 200 {
			if code == 212 {
				resp.Status = 1006
			}
			resp.Status = 1014
			return resp
		} else {
			return resp
		}
	}
}