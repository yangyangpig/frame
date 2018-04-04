package controllers

import (
	"PGUser/app/proto"
	"PGUser/app/service"
	"PGUser/app/entity"
)

type UserController struct {}

//写入user用户表
func (this *UserController) Add(rq *pgUser.InsertUserRequest) *pgUser.InsertUserResponse {
	resp := new(pgUser.InsertUserResponse)
	if rq.Mid <= 0 {
		return resp
	}
	service.UserService.Add(rq)
	return resp
}

//获取用户信息
func (this *UserController) Get(rq *pgUser.GetUserInfoRequest) *pgUser.GetUserInfoResponse {
	resp := new(pgUser.GetUserInfoResponse)
	if rq.Mid <= 0 {
		resp.Status = 1
		resp.Msg = "参数不合法"
		return resp
	}
	var (
		userInfo entity.User
		err error
	)
	userInfo, err = service.UserService.Get(rq.Mid)
	if err != nil {
		resp.Status = 2
		resp.Msg = "获取用户信息失败"
		return resp
	}
	resp.Data.Mid = userInfo.Mid
	resp.Data.AppId = userInfo.App_id
	resp.Data.Nick = userInfo.Nick
	return resp
}