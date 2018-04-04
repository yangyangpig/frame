package controllers

import (
	"Proto/login"
	"Proto/user"
	"Proto/online"
	"PGLogin/app/service"
	"PGLibrary"
	"time"
	"fmt"
	"github.com/astaxie/beego/logs"
	"Proto/money"
	"math/rand"
)
//登录注册服务
type LoginController struct {}

//错误提示说明
var (
	SUCCESS = "成功"
	FAIL = "失败"
)

//登录注册接口
func (this *LoginController) Login(req *pgLogin.LoginRequest) *pgLogin.LoginResponse {
	var (
		cid int64
		mid int64
		isNew int8
	)
	resp := new(pgLogin.LoginResponse)

	ssid := this.generateSsid(mid)

	cid = this.getCid(req)
	//失败
	if cid == 0 {
		resp.Status = 1
		resp.Msg = FAIL
		return resp
	}
	mid, isNew = this.getMid(req)
	//失败
	if mid == 0 {
		resp.Status = 1
		resp.Msg = FAIL
		return resp
	}
	//注册
	if isNew == 1 {
		//rpc调用用户服务新增
		result := this.addUser(ssid, mid, cid, req, resp)
		if !result {
			resp.Status = 1
			resp.Msg = FAIL
			return resp
		}
	} else { //登录
		//mid存在，获取用户信息
		result := this.getUserInfo(ssid, mid, req, resp)
		if !result {
			resp.Status = 1
			resp.Msg = FAIL
			return resp
		}
	}
	//设置用户状态
	result := this.setOnline(req.AppId, req.HallVer, req.ApkVer, 0, 0, ssid, mid)
	if !result {
		resp.Status = 1
		resp.Msg = FAIL
		return resp
	}
	return resp
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
		logs.Error(err)
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("coincache.CoinCacheService.Create", reqBytes, 1000, mid)
	if response.ReturnCode != 0 {
		logs.Error(fmt.Sprintf("LoginController.AddMoney Err: returncode %d", response.ReturnCode))
		return false
	}
	return true
}

//生成ssid
func (this *LoginController) generateSsid(mid int64) string {
	currentTime := time.Now()
	timeStamp := currentTime.Unix()
	tmpStr := fmt.Sprintf("%d@%d", mid, timeStamp)
	ssid := PGLibrary.GetMd5(tmpStr)
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
		logs.Error(err)
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("online.OnlineService.Report", reqBytes, 1000, 0)
	if response.ReturnCode != 0 {
		logs.Error(fmt.Sprintf("LoginController.setOnline Err: returncode %d", response.ReturnCode))
		return false
	}
	return true
}

//获取用户信息
func (this *LoginController) getUserInfo(ssid string, mid int64, rq *pgLogin.LoginRequest, resp *pgLogin.LoginResponse) bool {
	request := new(pgUser.GetUserInfoRequest)
	request.Mid = mid
	reqBytes, err := request.Marshal()
	if err != nil {
		logs.Error("LoginController::getUserInfo Err: Marshal fail")
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("PGUser.Get", reqBytes, 1000, 0)
	if response.ReturnCode != 0 {
		logs.Error(fmt.Sprintf("LoginController.getUserInfo Err: returncode %d", response.ReturnCode))
		return false
	}
	pgResp := new(pgUser.GetUserInfoResponse)
	pgResp.Unmarshal(response.Body)
	//成功
	if pgResp.Status == 0 {
		resp.Status = 0
		resp.Msg = SUCCESS
		resp.Data.Mid = mid
		resp.Data.Nick = pgResp.Data.Nick
		resp.Data.Sex = pgResp.Data.Sex
		resp.Data.Icon = pgResp.Data.Icon
		resp.Data.Ssid = ssid
		resp.Data.Money = 8888
		resp.Data.LoginType = rq.LoginType
		return true
	} else { //失败
		logs.Error("LoginController.getUserInfo Err: 获取用户信息失败")
		return false
	}
}

//rpc调用用户模块新增用户信息
func (this *LoginController) addUser(ssid string, mid, cid int64, rq *pgLogin.LoginRequest, resp *pgLogin.LoginResponse) bool {
	request := new(pgUser.InsertUserRequest)
	request.Mid = mid
	request.AppId = rq.AppId
	request.Cid = cid
	request.PlatformId = rq.Guid
	request.PlatformType = rq.LoginType
	request.Nick = "Nick" + fmt.Sprintf("%d", rand.Intn(10000))
	reqBytes, err := request.Marshal()
	if err != nil {
		logs.Error("LoginController.addUser Err: Marshal fail")
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("PGUser.Add", reqBytes, 1000, 0)
	if response.ReturnCode != 0 {
		logs.Error(fmt.Sprintf("LoginController.addUser Err: returncode %d", response.ReturnCode))
		return false
	}
	pgResp := new(pgUser.InsertUserResponse)
	pgResp.Unmarshal(response.Body)
	//成功
	if pgResp.Status == 0 {
		resp.Status = 0
		resp.Msg = SUCCESS
		resp.Data.Mid = mid
		resp.Data.Nick = request.Nick
		resp.Data.Sex = 1
		resp.Data.Icon = ""
		resp.Data.Ssid = ssid
		resp.Data.Money = 8888
		resp.Data.LoginType = rq.LoginType
		return true
	} else { //失败
		logs.Error("LoginController.addUser Err: 写入用户信息失败")
		return false
	}
}

//检测cid是否存在，不存在就写入cid
func (this *LoginController) getCid(rq *pgLogin.LoginRequest) int64 {
	cid, err := service.CidMapService.GetCid(rq.Guid, rq.LoginType)
	if err != nil {
		logs.Error("LoginController.cidIsExist Err: 获取CID失败")
		return 0
	} else if cid == 0 {//cid不存在
		cid, err = service.CidMapService.Insert(rq.Guid, rq.LoginType)
		if err != nil {
			logs.Error("LoginController.cidIsExist Err: 写入CID失败")
			return 0
		}
	}
	return cid
}

//检测mid是否存在，不存在就写入mid
func (this *LoginController) getMid(rq *pgLogin.LoginRequest) (int64, int8) {
	region := PGLibrary.GetRegionId(rq.AppId)
	//查mid是否存在
	mid, err := service.MidMapService.GetMid(rq.Guid, rq.LoginType, region)
	if err != nil {
		logs.Error("LoginController.midIsExist Err: 获取MID失败")
		return 0, 0
	} else if mid == 0 {
		mid, err = service.MidMapService.Insert(rq.Guid, rq.LoginType, region)
		if err != nil {
			logs.Error("LoginController.midIsExist Err: 写入MID失败")
			return 0, 0
		}
		return 0, 1
	}
	return mid, 0
}
