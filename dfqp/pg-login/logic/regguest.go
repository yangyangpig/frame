package logic

import (
	"dfqp/proto/login"
	"dfqp/pg-login/service"
)

var GuestRegLogic = &guestRegLogic{}

//游客账号注册
type guestRegLogic struct {
	BaseLogic
}

//注册和登录一个游客账号
func (this *guestRegLogic) Reg(req *pgLogin.LoginRequest) *RegResponse {
	resp := new(RegResponse)
	resp.ErrCode = 1000
	resp.IsReg = false

	guid := req.Guid
	if guid == "" {
		resp.ErrCode = 1001
		return resp
	}
	cid, err := service.Platform2cidService.GetCid(req.Guid, guestType)
	if err != nil {
		return resp
	}
	if cid == 0 {
		cid = this.addUserInfo(guid, guestType, 0, this.getNickName(), "", "")
		if cid > 0 {
			resp.ErrCode = 0
			resp.Cid = cid
			resp.IsReg = true
		}
	} else {
		resp.ErrCode = 0
		resp.Cid = cid
	}
	return resp
}