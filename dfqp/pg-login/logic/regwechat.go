package logic

import (
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"dfqp/proto/login"
	"putil/log"
	"time"
	"fmt"
	"math/rand"
	"dfqp/pg-login/service"
)

//实例化
var WechatRegLogic = &wechatRegLogic{}

//微信账号注册
type wechatRegLogic struct {
	BaseLogic
}

//微信登录
func (this *wechatRegLogic) Reg(req *pgLogin.LoginRequest) *RegResponse {
	resp := new(RegResponse)
	resp.ErrCode = 1000
	resp.IsReg = false //默认是登录

	unionId := req.WechatParam.UnionId
	openId := req.WechatParam.OpenId
	accessToken := req.WechatParam.AccessToken
	guid := req.Guid

	//检测参数
	if len(unionId) == 0 || len(openId) == 0 || len(accessToken) == 0 || len(guid) == 0 {
		resp.ErrCode = 1001
		return resp
	}

	//先检测该微信是否注册过
	cid, err := service.Platform2cidService.GetCid(unionId, wechatType)
	if err != nil {
		return resp
	}
	//注册过
	if cid > 0 {
		resp.ErrCode = 0
		resp.Cid = cid
	} else {
		//获取微信用户资料
		data, errCode := this.getUserInfo(accessToken, openId)
		//获取用户信息成功
		if errCode == 0 {
			var (
				sex int32
				nick string
				icon string
				ok bool
			)
			if sex, ok = data["sex"].(int32); !ok {
				sex = 0
			}
			if icon, ok = data["icon"].(string); !ok {
				icon = ""
			}
			if nick, ok = data["nick"].(string); !ok {
				nick = ""
			}
			cid, err = service.Platform2cidService.GetCid(guid, guestType)
			if err != nil {
				return resp
			}
			//游客账号存在
			if cid > 0 {
				//是否绑定了手机
				newBid, err := service.CidMapService.GetPlatformId(cid, boyaaucType)
				if err != nil {
					return resp
				}
				//是否绑定其他微信
				newUnionId, err := service.CidMapService.GetPlatformId(cid, wechatType)
				if err != nil {
					return resp
				}
				//说明游客账号没有被绑定过手机或微信
				if len(newBid) == 0 && len(newUnionId) == 0 {
					ret := this.modifyUserInfo(cid, unionId, sex, nick, icon)
					if ret {
						resp.ErrCode = 0
						resp.Cid = cid
					}
				} else {
					//生成新的cid
					cid = this.addUserInfo(unionId, wechatType, sex, nick, icon, "")
					if cid > 0 {
						resp.ErrCode = 0
						resp.Cid = cid
					}
				}
			} else {
				//生成新的cid
				cid = this.addUserInfo(unionId, wechatType, sex, nick, icon, "")
				if cid > 0 {
					resp.ErrCode = 0
					resp.Cid = cid
					resp.IsReg = true
				}
			}
		} else if errCode == 2 { //token过期
			resp.ErrCode = 1005
		}
	}
	return resp
}

//添加微信用户平台标识和修改用户信息
func (this *wechatRegLogic) modifyUserInfo(cid int64, platformId string, sex int32, nick, icon string) bool {
	ret := this.addPlatform(cid, platformId, wechatType)
	if ret {
		//将游客信息修改为用户信息
		ret = this.modifyUser(cid, sex, nick, icon, icon, "")
		if ret {
			return true
		}
	}
	return false
}

//获取微信用户基本信息
func (this *wechatRegLogic) getUserInfo(accessToken, openId string) (map[string]interface{}, int32) {
	data := make(map[string]interface{})
	url := "https://api.weixin.qq.com/sns/userinfo?access_token="+accessToken+"&openid="+openId
	ret, err := http.Get(url)
	if err != nil {
		return nil, 1
	}
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, 1
	}
	plog.Debug("getUserInfo resp: ", string(body))
	jsonData := gjson.Parse(string(body))
	openId = jsonData.Get("openid").String()
	errCode := jsonData.Get("errcode").Int()
	if len(openId) > 0 {
		data["icon"] = jsonData.Get("headimgurl").String()
		data["sex"] = jsonData.Get("sex").Int() //值为1时是男性，值为2时是女性，值为0时是未知
		data["nick"] = jsonData.Get("nickname").String()
		data["unionId"] = jsonData.Get("unionid").String()
		return data, 0
	} else if errCode == 40001 || errCode == 42001 {//access_token过期
		return nil, 2
	} else {
		return nil, 1
	}
}

//游客绑定微信
func (this *wechatRegLogic) GuestBindWechat(req *pgLogin.GuestBindWechatRequest) *RegResponse {
	cid := req.Mid
	resp := new(RegResponse)
	resp.ErrCode = 1018
	//是否已经绑定
	platformId, err := service.CidMapService.GetPlatformId(cid, wechatType)
	if err != nil {
		return resp
	}
	if len(platformId) > 0 {
		if platformId == req.UnionId { //直接登录
			resp.ErrCode = 0
			resp.LoginType = wechatType
			return resp
		} else {
			resp.ErrCode = 1004
			return resp
		}
	}
	//检测该微信是否已经被其他cid绑定了
	ccid, err := service.Platform2cidService.GetCid(req.UnionId, wechatType)
	if err != nil {
		return resp
	}
	if ccid > 0 {
		resp.ErrCode = 1019
		return resp
	}
	data, errorCode := this.getUserInfo(req.AccessToken, req.OpenId)
	if errorCode == 2 {
		resp.ErrCode = 1005
	} else if errorCode == 0 { //验证成功
		var (
			sex int32
			nick string
			icon string
			ok bool
		)
		if sex, ok = data["sex"].(int32); ok {
			sex = 0
		}
		if icon, ok = data["icon"].(string); ok {
			icon = ""
		}
		if nick, ok = data["nick"].(string); ok {
			nick = ""
		}
		//写入微信平台标识
		ret := this.modifyUserInfo(cid, req.UnionId, sex, nick, icon)
		if ret {
			resp.ErrCode = 0
			resp.LoginType = wechatType
		}
	}
	return resp
}

//微信绑定手机
func (this *wechatRegLogic) WechatBindPhone(req *pgLogin.WechatBindPhoneRequest) *RegResponse {
	resp := new(RegResponse)
	resp.ErrCode = 1008
	captcha := req.Captcha
	phone := req.Phone

	//是否已经绑定了手机
	bid, err := service.CidMapService.GetPlatformId(req.Mid, boyaaucType)
	if err != nil {
		return resp
	}
	//已经绑定过了
	if bid != "" {
		resp.ErrCode = 1011
		return resp
	}
	params := make(map[string]interface{})
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	pwd := fmt.Sprintf("%v", rnd.Int31n(100000000))
	params["phone"] = phone
	params["pwd"] = pwd
	params["token"] = captcha
	result, err := service.ByClientService.Get("user/register", params)
	plog.Debug("wechatBindPhone resp: ", result)
	if err == nil {
		jsonData := gjson.Parse(result)
		if !jsonData.Exists() {
			resp.ErrCode = 1
		} else {
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
					resp.ErrCode = 0
					resp.Pwd = pwd
				}
			}
		}
	}
	return resp
}