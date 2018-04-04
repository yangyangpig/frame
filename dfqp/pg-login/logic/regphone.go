package logic

import (
	"dfqp/proto/login"
	"github.com/tidwall/gjson"
	"time"
	"fmt"
	"dfqp/pg-login/service"
	"dfqp/lib"
	"putil/log"
	"github.com/garyburd/redigo/redis"
	"math/rand"
)

var PhoneRegLogic = &phoneRegLogic{}

//手机账号注册
type phoneRegLogic struct {
	BaseLogic
}
//注册登录
func (this *phoneRegLogic) Reg(req *pgLogin.LoginRequest) *RegResponse {
	resp := new(RegResponse)
	resp.ErrCode = 1000 //默认登录失败
	resp.IsReg = false //默认是登录

	captcha := req.PhoneParam.Captcha
	phone := req.PhoneParam.Phone
	bid := req.PhoneParam.Bid
	accessToken := req.PhoneParam.AccessToken
	if !lib.IsTelephone(phone, "CHN") {
		resp.ErrCode = 1015
		return resp
	}
	guid := req.Guid
	//检查token和bid是否合法
	if len(bid) > 0 {
		errCode := this.checkToken(req)
		plog.Debug("phoneRegLogic.checkToken resp: ", errCode)
		if errCode == 0 { //校验成功
			cid, err := service.Platform2cidService.GetCid(bid, boyaaucType)
			if err != nil {
				return resp
			}
			resp.ErrCode = 0
			resp.Cid = cid
			resp.Bid = bid
			resp.AccessToken = accessToken
		} else {
			resp.ErrCode = 1012
		}
	} else {
		//验证码登录
		if len(captcha) > 0 {
			bid, errCode := this.phoneCaptchaLogin(req)
			plog.Debug("phoneRegLogic.phoneCaptchaLogin resp: ", bid, errCode)
			if errCode == 2 { //验证码错误
				resp.ErrCode = 1006
			} else if errCode == 0 { //成功
				//是否存在cid
				cid, err := service.Platform2cidService.GetCid(bid, boyaaucType)
				if err != nil {
					return resp
				}
				if cid > 0 {
					//登录
					resp.ErrCode = 0
					resp.Cid = cid
					result := this.setToken(bid, phone, guid)
					resp.Bid = result["bid"]
					resp.AccessToken = result["accessToken"]
					return resp
				} else {
					//是否存在游客账号
					var bindGuestFlag int32
					guestCid, err := service.Platform2cidService.GetCid(req.Guid, guestType)
					if err != nil {
						return resp
					}
					if guestCid > 0 {
						unionId, err := service.CidMapService.GetPlatformId(guestCid, wechatType)
						if err != nil {
							return resp
						}
						newBid, err := service.CidMapService.GetPlatformId(guestCid, boyaaucType)
						if err != nil {
							return resp
						}
						if len(unionId) > 0 || len(newBid) > 0 { //说明已经绑定过其他账号
							//生成新账号
							bindGuestFlag = 1
						} else { //绑定游客账号
							bindGuestFlag = 2
						}
					} else {
						bindGuestFlag = 1
					}
					plog.Debug("phoneRegLogic bind guest flag: ", bindGuestFlag)
					//生成新账号
					if bindGuestFlag == 1 {
						//生成新账号
						cid := this.addUserInfo(bid, boyaaucType, 0, this.getNickName(), "", phone)
						if cid > 0 {
							resp.ErrCode = 0
							resp.Cid = cid
							resp.IsReg = true
							result := this.setToken(bid, phone, guid)
							resp.Bid = result["bid"]
							resp.AccessToken = result["accessToken"]
						}
					} else {//绑定游客账号
						ret := this.modifyUserInfo(cid, bid, phone)
						if ret {
							resp.ErrCode = 0
							resp.Cid = guestCid
							result := this.setToken(bid, phone, guid)
							resp.Bid = result["bid"]
							resp.AccessToken = result["accessToken"]
						}
					}
				}
			}
		} else {//密码登录
			ret := this.phoneIsReg(req)
			//手机号不存在，提示验证码登录
			if ret == 0 {
				resp.ErrCode = 1007
				return resp
			}
			pwd := req.PhoneParam.Pwd
			if len(pwd) < 6 || len(pwd) > 12 {
				resp.ErrCode = 1016
				return resp
			}
			bid, errCode := this.phonePwdLogin(req)
			if errCode == 0 { //登录成功
				cid, err := service.Platform2cidService.GetCid(bid, boyaaucType)
				if err != nil {
					return resp
				}
				resp.Cid = cid
				resp.ErrCode = 0
				result := this.setToken(bid, phone, guid)
				resp.Bid = result["bid"]
				resp.AccessToken = result["accessToken"]
			}
		}
	}
	return resp
}

//添加手机用户平台标识和修改用户信息
func (this *phoneRegLogic) modifyUserInfo(cid int64, bid string, phone string) bool {
	ret := this.addPlatform(cid, bid, boyaaucType)
	if ret {
		//将游客信息修改为用户信息
		ret = this.modifyUser(cid, 0, this.getNickName(), "", "", phone)
		if ret {
			return true
		}
	}
	return false
}

//通行证检测token是否过期
func (this *phoneRegLogic) checkToken(req *pgLogin.LoginRequest) int32 {
	accessToken := req.PhoneParam.AccessToken
	bid := req.PhoneParam.Bid
	redisObj, err := service.NewRedisCache("userLogin")
	if err == nil {
		token, err := redis.String(redisObj.Do("GET", bid))
		if err == nil {
			if accessToken == token {
				curTime := time.Now().Unix()
				redisObj.Do("EXPIRE", bid, curTime+byAccountExpires)
				return 0
			} else {
				return 1 //token不一致
			}
		} else {
			return 2 //error
		}
	} else {
		return 2 //error
	}
}

//手机验证码登录
func (this *phoneRegLogic) phoneCaptchaLogin(req *pgLogin.LoginRequest) (string, int32) {
	params := make(map[string]interface{})
	params["phone"] = req.PhoneParam.Phone
	params["token"] = req.PhoneParam.Captcha
	result, err := service.ByClientService.Get("user/phonelogin", params)
	plog.Debug("phoneRegLogic.phoneCaptchaLogin resp: ", result, err)
	if err != nil {
		return "", 1
	} else {
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		if code != 200 {
			if code == 212 { //验证码错误
				return "", 2
			} else {
				return "", 1
			}
		} else {
			bid := jsonData.Get("result.bid").String()
			return bid, 0
		}
	}
}

//设置token
func (this *phoneRegLogic) setToken(bid, phone, guid string) map[string]string {
	data := make(map[string]string)
	curTime := time.Now().Unix()
	timeStr := fmt.Sprintf("%d", curTime)
	token := lib.GetMd5("GODFQPTOKEN"+bid+phone+guid+timeStr)
	redisObj, err := service.NewRedisCache("userLogin")
	plog.Debug("NewRedisCache fail: ", err)
	if err == nil {
		redisObj.Do("SET", bid, token)
		redisObj.Do("EXPIRE", bid, curTime+byAccountExpires)
	}
	data["bid"] = bid
	data["accessToken"] = token
	return data
}

//手机密码登录
func (this *phoneRegLogic) phonePwdLogin(req *pgLogin.LoginRequest) (string, int32) {
	phone := req.PhoneParam.Phone
	params := make(map[string]interface{})
	params["phone"] = phone
	params["type"] = "PHONE"
	params["pwd"] = req.PhoneParam.Pwd
	clientId := lib.GetClientId(req.AppId)
	params["devid"] = this.getDeviceType(clientId)
	result, err := service.ByClientService.Get("user/check", params)
	plog.Debug("phoneRegLogic.phonePwdLogin resp: ", result, err)
	if err != nil {
		return "", 1
	} else {
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		if code != 200 {
			return "", 2
		} else {
			bid := jsonData.Get("result.bid").String()
			return bid, 0
		}
	}
}

//获取设备类型
func (this *phoneRegLogic) getDeviceType(clientId int32) int32 {
	if clientId == 1 {
		return byPassIphone
	} else if clientId == 2 {
		return byPassAndroid
	}
	return 0
}

//手机是否注册了
// 0手机不存在 1手机账号存在 2错误
func (this *phoneRegLogic) phoneIsReg(req *pgLogin.LoginRequest) int32 {
	phone := req.PhoneParam.Phone
	params := make(map[string]interface{})
	params["passportcount"] = phone
	params["type"] = "PHONE"
	result, err := service.ByClientService.Get("user/getuser", params)
	plog.Debug("phoneRegLogic.phoneIsReg resp: ", result, err)
	if err == nil {
		jsonData := gjson.Parse(result)
		code := jsonData.Get("code").Int()
		if code == 218 { //该手机帐号不存在
			return 0
		}
		if code == 200 {
			bid := jsonData.Get("result.bid").String()
			if len(bid) > 0 { //手机账号存在
				return 1
			} else {
				return 0
			}
		} else {
			return 2
		}
	} else {
		return 2
	}
}

//游客绑定手机
func (this *phoneRegLogic) GuestBindPhone(req *pgLogin.GuestBindPhoneRequest) *RegResponse {
	resp := new(RegResponse)
	resp.ErrCode = 1009
	if req.Captcha == "" {
		resp.ErrCode = 1001
		return resp
	}
	if !lib.IsTelephone(req.Phone, "CHN") {
		resp.ErrCode = 1015
		return resp
	}
	unionId, err := service.CidMapService.GetPlatformId(req.Mid, wechatType)
	if err != nil {
		return resp
	}
	if unionId != "" {
		resp.ErrCode = 1004
		return resp
	}
	bid, err := service.CidMapService.GetPlatformId(req.Mid, boyaaucType)
	if err != nil {
		return resp
	}
	if bid != "" {
		resp.ErrCode = 1010
		return resp
	}
	params := make(map[string]interface{})
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	pwd := fmt.Sprintf("%06v", rnd.Int31n(100000000))
	params["phone"] = req.Phone
	params["pwd"] = pwd
	params["token"] = req.Captcha
	result, err := service.ByClientService.Get("user/register", params)
	plog.Debug("phoneRegLogic.GuestBindPhone resp: ", result)
	if err == nil {
		jsonData := gjson.Parse(result)
		if jsonData.Exists() {
			code := jsonData.Get("code").Int()
			if code == 214 { //该手机帐号已经注册
				resp.ErrCode = 1003
			} else if code == 212 { //验证码错误
				resp.ErrCode = 1006
			} else if code == 200 {
				bid := jsonData.Get("result.bid").String()
				ret := this.addPlatform(req.Mid, bid, boyaaucType)
				if ret {
					res := this.modifyUser(req.Mid, 0, "", "", "", req.Phone)
					if !res {
						plog.Fatal("modify user info fail")
					}
					guid, _ := service.CidMapService.GetPlatformId(req.Mid, guestType)
					data := this.setToken(bid, req.Phone, guid)
					resp.ErrCode = 0
					resp.LoginType = boyaaucType
					resp.Pwd = pwd
					resp.Bid = data["bid"]
					resp.AccessToken = data["accessToken"]
				}
			}
		}
	}
	return resp
}
