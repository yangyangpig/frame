package controller

import (
	"dfqp/proto/login"
	"dfqp/pg-login/service"
)

type LogInfoController struct {}

func (this *LogInfoController)GetLogInfo(req *pgLogin.LogInfoRequest) *pgLogin.LogInfoResponse {
	// 初始化
	resp := new(pgLogin.LogInfoResponse)
	resp.Status = 1
	resp.Msg = "请求失败！"
	// 参数限定
	if req.Cid <= 0 {
		return resp
	}
	data, err := service.LoginInfoService.GetInfoByCid(req.Cid)
	if err != nil {
		return resp
	}
	resp.Status = 0
	resp.Msg = "请求成功"

	resp.Data.Cid           = data.Cid
	resp.Data.FirstApp      = data.FirstApp
	resp.Data.LastApp       = data.LastApp
	resp.Data.FirstVersion  = data.FirstVersion
	resp.Data.LastVersion   = data.LastVersion
	resp.Data.RegTime       = data.RegTime
	resp.Data.LoginTime     = data.LoginTime
	resp.Data.FirstIp       = data.FirstIp
	resp.Data.LastIp        = data.LastIp
	return resp
}
