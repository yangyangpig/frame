package controller

import (
	"sync"
	"strconv"
	"net/http"
	"dfqp/proto/user"
	"dfqp/proto/login"
	"dfqp/pg-http/service"
)

type UserController struct {
	base
}

const (
	Success = iota
	Faild
)

// 获取用户信息
func (this *UserController)GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	cid, err := strconv.ParseInt(request.Form.Get("cid"), 10, 64)
	// 初始化
	code := Success

	if err != nil {
		code = Faild
		this.output(writer, code, "", "")
		return
	}

	respData := make(map[string]interface{})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		// 获取用户信息
		req := new(pgUser.GetUserInfoRequest)
		req.Mid = cid
		reqBytes, _ := req.Marshal()
		response := service.Client.SendAndRecvRespRpcMsg("pgUser.getUserInfo", reqBytes, 2000, 0)
		if response.ReturnCode != 0 {
			code = Faild
		} else {
			resData := new(pgUser.GetUserInfoResponse)
			resData.Unmarshal(response.Body)
			if resData.Status == 0 {
				respData["nick"]    = resData.Data.Nick
				respData["sex"]     = resData.Data.Sex
				respData["icon"]    = resData.Data.Icon
				respData["icon_big"]= resData.Data.IconBig
				respData["city"]    = resData.Data.City
				respData["sign"]    = resData.Data.Sign
				respData["status"]  = resData.Data.Status
				respData["icon_id"] = resData.Data.IconId
				respData["phone"]   = resData.Data.Phone
			}
		}
	}()
	go func() {
		defer wg.Done()
		// 获取登陆日志
		req := new(pgLogin.LogInfoRequest)
		req.Cid = cid
		reqBytes, _ := req.Marshal()
		response := service.Client.SendAndRecvRespRpcMsg("pgLogin.getLoginInfo", reqBytes, 2000, 0)
		if response.ReturnCode != 0 {
			code = Faild
		} else {
			resData := new(pgLogin.LogInfoResponse)
			resData.Unmarshal(response.Body)
			if resData.Status == 0 {
				respData["cid"]             = resData.Data.Cid
				respData["first_app"]       = resData.Data.FirstApp
				respData["last_app"]        = resData.Data.LastApp
				respData["first_version"]   = resData.Data.FirstVersion
				respData["last_version"]    = resData.Data.LastVersion
				respData["reg_time"]        = resData.Data.RegTime
				respData["login_time"]      = resData.Data.LoginTime
				respData["first_ip"]        = resData.Data.FirstIp
				respData["last_ip"]         = resData.Data.LastIp
			}
		}
	}()

	wg.Wait()

	this.output(writer, code, respData, "")
	return
}

// 修改用户昵称
func (this *UserController)ModifyNick(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	cid, err := strconv.ParseInt(request.PostFormValue("cid"), 10, 64)
	nick := request.PostFormValue("nick")
	// 初始化
	code := Success

	if err != nil {
		code = Faild
		this.output(writer, code, "", "")
		return
	}

	req := new(pgUser.ModifyUserInfoRequest)
	req.Mid  = cid
	req.Nick = nick
	reqBytes, _ := req.Marshal()
	response := service.Client.SendAndRecvRespRpcMsg("pgUser.modifyUserInfo", reqBytes, 2000, 0)
	if response.ReturnCode != 0 {
		code = Faild
	}
	this.output(writer, code, "", "")
	return
}
