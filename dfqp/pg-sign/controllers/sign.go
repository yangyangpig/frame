package controllers

import (
	"dfqp/proto/sign"
	"dfqp/pg-sign/service"
)

// 签到注册服务
type SignController struct {}

// 获取签到信息
func (this *SignController) GetSigninInfos(req *pgSign.GetSigninInfosRequest) *pgSign.GetSigninInfosResponse {
	resp := new(pgSign.GetSigninInfosResponse)

	resp.TodayInfos.Num = 1

	return resp
}

// 签到
func (this *SignController) Signin(req *pgSign.SigninRequest) *pgSign.SigninResponse {
	resp := new(pgSign.SigninResponse)

	// 插入数据库测试
	service.SigninService.Add(2233, 12, "sdfsdfsfd")
	resp.Status = 1
	return resp
}