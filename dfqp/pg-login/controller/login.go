package controller

import (
	"dfqp/proto/login"
	"time"
	"fmt"
	"dfqp/pg-login/logic"
	"dfqp/lib"
	"dfqp/proto/online"
	"dfqp/proto/money"
	"dfqp/pg-login/service"
	"putil/log"
	"dfqp/proto/user"
	"sync"
)
//登录注册服务
type LoginController struct {

}

//登录注册接口
func (this *LoginController) Login(req *pgLogin.LoginRequest) *pgLogin.LoginResponse {
	plog.Debug("login params: ", req)
	resp := new(pgLogin.LoginResponse)
	if len(req.Guid) == 0 || req.AppId <= 0 {
		resp.Status = 1001
		return resp
	}
	var (
		cid int64
		ssid string
		isReg bool
		ret *logic.RegResponse
	)
	loginType := req.LoginType
	if loginType == 1 { //游客账号
		ret = logic.GuestRegLogic.Reg(req)
	} else if loginType == 14 { //微信登录
		ret = logic.WechatRegLogic.Reg(req)
	} else if loginType == 2 { //博雅通行证
		ret = logic.PhoneRegLogic.Reg(req)
		resp.Data.PhoneParam.AccessToken = ret.AccessToken
		resp.Data.PhoneParam.Bid = ret.Bid
	} else {
		resp.Status = 1001
		return resp
	}
	//错误
	if ret.ErrCode > 0 {
		resp.Status = ret.ErrCode
		return resp
	} else {
		cid = ret.Cid
		if cid > 0 {
			//检测是否存在表注册时表写入失败得情况
			step, err := service.RegFailService.Get(cid)
			if step > 1 || err != nil {
				resp.Status = 1000
				return resp
			}
			isReg = ret.IsReg
			//注册需加银币
			if isReg {
				go this.AddMoney(cid)
			}
			ssid = this.generateSsid(cid)
			//设置用户状态
			result := this.setOnline(req.AppId, req.HallVer, req.ApkVer, 0, 0, ssid, cid)
			if !result {
				resp.Status = 1000
				return resp
			}
			//并发获取start
			isSuccess := true
			var (
				wg sync.WaitGroup
			)
			wg.Add(2)
			go func() {
				//获取用户信息
				userInfo := this.getUserInfo(cid)
				plog.Debug("get User info:", userInfo)
				if userInfo.Status == 0 {
					resp.Data.UserInfoParam.Nick = userInfo.Data.Nick
					resp.Data.UserInfoParam.Sex = userInfo.Data.Sex
					resp.Data.UserInfoParam.Icon = userInfo.Data.Icon
					resp.Data.PhoneParam.Phone = userInfo.Data.Phone
				} else {
					isSuccess = false
				}
				defer wg.Done()
			}()
			go func() {
				//获取用户资产
				userProperty, flag := this.getUserProperty(cid)
				plog.Debug("get User property:", userProperty, flag)
				if flag {
					attr := userProperty.AttrList
					for _, v := range attr {
						//银币
						if v.Type == int32(pgMoney.EUserAttr_MONEY) {
							resp.Data.UserPropertyParam.Silver = v.Value
						}
						//金条
						if v.Type == int32(pgMoney.EUserAttr_SILVER) {
							resp.Data.UserPropertyParam.Bullion = v.Value
						}
					}
				} else {
					isSuccess = false
				}
				defer wg.Done()
			}()
			wg.Wait()
			//并发获取end
			//成功
			if isSuccess {
				resp.Status = 0
				resp.Data.Cid = cid
				resp.Data.Mid = cid
				resp.Data.Ssid = ssid
				resp.Data.LoginType = loginType
				//放入协程处理
				go func() {
					onlineInfo := this.getOnlineInfo(cid)
					this.AddLog(cid, req.AppId, req.AppId, req.ApkVer, req.ApkVer, onlineInfo.Ip, onlineInfo.Ip, isReg)
				}()
			} else { //失败
				resp.Status = 1000
			}
			return resp
		} else {
			resp.Status = 1000
			return resp
		}
	}
}

//添加日志
func (this *LoginController) AddLog(cid int64, firstApp, lastApp int32, firstVersion, lastVersion, firstIp, lastIp string, isReg bool) {
	plog.Debug("isReg: ", isReg)
	if isReg {
		ret := service.LoginInfoService.Insert(cid, firstApp, lastApp, firstVersion, lastVersion, firstIp, lastIp)
		if !ret {
			//todo
		}
	} else {
		service.LoginInfoService.Update(cid, lastApp, lastVersion, lastIp)
	}
}

//加银币
func (this *LoginController) AddMoney(mid int64) bool {
	req := new(pgMoney.CoinCacheCreateRequest)
	userAtt := pgMoney.UserAttr{}
	userAtt.Type = int32(pgMoney.EUserAttr_MONEY)
	userAtt.Value = 10000
	req.Uid = mid
	req.Userattr = append(req.Userattr, userAtt)
	reqBytes, err := req.Marshal()
	if err != nil {
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("userserver.UserServerService.PhpCreateRecord", reqBytes, 1000, mid)
	if response.ReturnCode != 0 {
		return false
	}
	return true
}

//生成ssid
func (this *LoginController) generateSsid(cid int64) string {
	currentTime := time.Now()
	timeStamp := currentTime.Unix()
	tmpStr := fmt.Sprintf("%d@%d", cid, timeStamp)
	ssid := lib.GetMd5(tmpStr)
	return ssid
}

//用户在线状态
func (this *LoginController) setOnline(app int32,hallVer int32,apkVer string, cityApp int32, forbid int32, ssid string, mid int64) bool {
	request := new(pgOnline.ReportOnlineRequest)
	request.Uid = mid
	request.App = app
	request.HallVer = hallVer
	request.ApkVer = apkVer
	request.CityApp = cityApp
	request.Forbid = forbid
	request.Ssid = ssid
	reqBytes, err := request.Marshal()

	if err != nil {
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("online.OnlineService.Report", reqBytes, 1000, 0)
	if response.ReturnCode != 0 {
		return false
	}
	return true
}

//获取用户在线信息
func (this *LoginController) getOnlineInfo(mid int64) *pgOnline.GetOnlineResponse {
	response := new(pgOnline.GetOnlineResponse)
	request := new(pgOnline.GetOnlineRequest)
	request.Uid = mid
	reqBytes, err := request.Marshal()
	if err != nil {
		return nil
	}
	resp := service.Client.SendAndRecvRespRpcMsg("online.OnlineService.Get", reqBytes, 1000, 0)
	if resp.ReturnCode != 0 {
		return nil
	}
	response.Unmarshal(resp.Body)
	return response
}

//获取用户基本信息
func (this *LoginController) getUserInfo(mid int64) *pgUser.GetUserInfoResponse {
	response := new(pgUser.GetUserInfoResponse)
	request := new(pgUser.GetUserInfoRequest)
	request.Mid = mid
	reqBytes, err := request.Marshal()
	if err != nil {
		return response
	}
	resp := service.Client.SendAndRecvRespRpcMsg("pgUser.getUserInfo", reqBytes, 1000, 0)
	if resp.ReturnCode != 0 {
		return response
	}
	err = response.Unmarshal(resp.Body)
	if err != nil {
		return response
	}
	return response
}

//获取用户资产
func (this *LoginController) getUserProperty(mid int64) (*pgMoney.CoinGetUserInfoResponse, bool) {
	response := new(pgMoney.CoinGetUserInfoResponse)
	req := new(pgMoney.CoinGetUserInfoRequest)
	req.Uid = int32(mid)
	reqBytes, err := req.Marshal()
	if err != nil {
		return response, false
	}
	resp := service.Client.SendAndRecvRespRpcMsg("coincache.CoinCacheService.CoinGetUserInfo", reqBytes, 1000, mid)
	if resp.ReturnCode != 0 {
		return response, false
	}
	err = response.Unmarshal(resp.Body)
	if err != nil {
		return response, false
	}
	return response, true
}


