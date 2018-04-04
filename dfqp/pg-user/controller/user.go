package controller

import (
	"dfqp/proto/user"
	"dfqp/pg-user/service"
	"putil/log"
)

type UserController struct {}

//获取用户信息
func (this *UserController) Get(rq *pgUser.GetUserInfoRequest) *pgUser.GetUserInfoResponse {
	resp := new(pgUser.GetUserInfoResponse)
	mid := rq.GetMid()
	if mid <= 0 {
		resp.Status = 1
		return resp
	}
	userInfo, err := service.UserService.Get(mid)
	if err != nil {
		resp.Status = 2
		return resp
	}
	if userInfo == nil {
		resp.Status = 3
	} else {
		//成功
		resp.Status = 0
		resp.Data.Nick = userInfo.Nick
		resp.Data.Sex = userInfo.Sex
		resp.Data.IconBig = userInfo.Icon_big
		resp.Data.Icon = userInfo.Icon
		resp.Data.City = userInfo.City
		resp.Data.Sign = userInfo.Sign
		resp.Data.IconId = userInfo.IconId
		resp.Data.Phone = userInfo.Phone
	}
	plog.Debug("UserController.Get resp:", resp)
	resp.Data.Phone = userInfo.Phone
	return resp
}

//修改用户信息
func (this *UserController) Modify(req *pgUser.ModifyUserInfoRequest) *pgUser.ModifyUserInfoResponse {
	resp := new(pgUser.ModifyUserInfoResponse)
	mid := req.Mid
	if mid <= 0 {
		resp.Status = 1
		return resp
	}
	sex := req.GetSex()
	nick := req.GetNick()
	city := req.GetCity()
	sign := req.GetSign()
	icon := req.GetIcon()
	iconBig := req.GetIconBig()
	phone := req.GetPhone()
	iconId := req.GetIconId()
	ret := service.UserService.Update(mid, sex, nick, city, phone, sign, icon, iconBig, iconId)
	if !ret {
		//失败
		resp.Status = 2
	}
	return resp
}

//写入user
func (this *UserController) Add(rq *pgUser.InsertUserInfoRequest) *pgUser.InsertUserInfoResponse {
	resp := new(pgUser.InsertUserInfoResponse)
	cid := rq.GetCid()
	if cid <= 0 {
		resp.Status = 1 //1失败 0成功
		return resp
	}
	sex := rq.GetSex()
	nick := rq.GetNick()
	city := rq.GetCity()
	phone := rq.GetPhone()
	sign :=  rq.GetSign()
	icon := rq.GetIcon()
	iconBig := rq.GetIconBig()
	ret := service.UserService.Add(cid, sex, nick, city, phone, sign, icon, iconBig)
	//成功
	if ret {
		resp.Status = 0
	} else {
		resp.Status = 2
	}
	return resp
}