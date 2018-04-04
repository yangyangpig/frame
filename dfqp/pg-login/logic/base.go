package logic

import (
	"dfqp/pg-login/service"
	"dfqp/proto/user"
	"time"
	"math/rand"
	"putil/log"
	"dfqp/proto/autoidpro"
)

const (
	//游客类型
	guestType = 1
	//博雅通行证
	boyaaucType = 2
	//微信类型
	wechatType = 14
	//博雅通行证过期时间
	byAccountExpires = 2592000
	//注册客户端类型
	byPassIphone = 2
	byPassAndroid = 1
	//个性签名
	profileSign = "这家伙很懒，什么也没留下！"

	addCidMapFail = 1 //写入pg_cidmap表失败
	addPlatformId2CidFail = 2 //写入pg_platform2cid表失败
	addUsersFail = 3 //写入pg_users表失败
)

//登录注册绑定返回值
type RegResponse struct {
	ErrCode int32
	ErrMsg string
	Cid int64
	IsReg bool
	LoginType int32
	Pwd string
	AccessToken string
	Bid string
}

//基类
type BaseLogic struct {}

//生成cid
func (this *BaseLogic) generateCid(platformType int32) int64 {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//cid := int64(r.Intn(1000000))
	request := new(autoidpro.AutoidRequest)
	request.Btag = "user"
	reqBytes, _ := request.Marshal()
	response := service.Client.SendAndRecvRespRpcMsg("autoidmanager.GetId", reqBytes, 3000, 0)
	plog.Debug("response address is: ", &response, "seq is: ", response.Seq, "value is: ", response.Body)
	if response.ReturnCode != 0 {
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		return 0
	} else {
		arithResp := new(autoidpro.AutoidResponse)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return id  = ", arithResp.Bid)
		if arithResp.Bid >= 1  {
			platformId, err := service.CidMapService.GetPlatformId(arithResp.Bid, platformType)
			if err != nil {
				return 0
			}
			if len(platformId) > 0 {
				plog.Warn("bid duplicate: ", arithResp.Bid)
				return 0
			}
			return arithResp.Bid
		} else {
			plog.Warn("bid value is: ", arithResp.Bid)
			return 0
		}
	}
}

//生成昵称
func (this *BaseLogic) getNickName() string {
	firstArr := []string{"partA", "partB"}
	lastArr := []string{"partC", "partD"}
	n := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(firstArr))
	firstNickArr := nickNameMap[firstArr[n]]
	firstLen := len(firstNickArr)
	nn := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(firstLen)

	m := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(lastArr))
	lastNickArr := nickNameMap[lastArr[m]]
	lastLen := len(lastNickArr)
	mm := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(lastLen)

	nick := firstNickArr[nn] + lastNickArr[mm]

	return nick
}

//添加用户平台标识和新增用户信息
//@param platformId 平台ID
//@param platformType 平台类型
//@param sex 性别
//@param nick 昵称
//@param icon 头像
//@param phone 手机
func (this *BaseLogic) addUserInfo(platformId string, platformType int32, sex int32, nick, icon string, phone string) int64 {
	cid := this.generateCid(platformType)
	if cid == 0 {
		return 0
	}
	ret := this.addPlatform(cid, platformId, platformType)
	if ret {
		ret = this.addUser(cid, sex, nick, icon, icon, phone, profileSign)
		if ret {
			return cid
		}
	}
	return 0
}

//添加平台信息
//@param cid 用户公共ID
//@param platformId 平台ID
//@param platformType 平台类型
func (this *BaseLogic) addPlatform(cid int64, platformId string, platformType int32) bool {
	ret := service.Platform2cidService.Insert(cid, platformId, platformType)
	plog.Debug("BaseLogic.addPlatform platform resp: ", ret)
	if ret {
		ret = service.CidMapService.Insert(cid, platformId, platformType)
		plog.Debug("BaseLogic.addPlatform cidMap resp: ", ret)
		if !ret {
			//记录失败位置
			go service.RegFailService.Add(cid, addPlatformId2CidFail)
		}
	} else {
		go service.RegFailService.Add(cid, addCidMapFail)
	}
	return ret
}

//添加用户
//@param cid 用户公共ID
//@param sex 性别
//@param nick 昵称
//@param icon 头像
//@param bigIcon 头像大图
//@param phone 手机
//@param sign 签名
func (this *BaseLogic) addUser(cid int64, sex int32, nick, icon, bigIcon, phone, sign string) bool {
	request := new(pgUser.InsertUserInfoRequest)
	request.Cid = cid
	request.Nick = nick
	request.Sex = sex
	request.Icon = icon
	request.IconBig = bigIcon
	request.Phone = phone
	request.Sign = sign
	reqBytes, err := request.Marshal()
	if err != nil {
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("pgUser.addUserInfo", reqBytes, 2000, 0)
	if response.ReturnCode != 0 {
		go service.RegFailService.Add(cid, addUsersFail)
		return false
	}
	pgResp := new(pgUser.InsertUserInfoResponse)
	pgResp.Unmarshal(response.Body)
	plog.Debug("BaseLogic.addUser resp: ", response)
	//成功
	if pgResp.Status == 0 {
		return true
	} else { //失败
		go service.RegFailService.Add(cid, addUsersFail)
		return false
	}
}
//修改用户基本信息
//@param cid 用户公共ID
//@param sex 性别
//@param nick 昵称
//@param icon 头像
//@param bigIcon 头像大图
//@param phone 手机
func (this *BaseLogic) modifyUser(cid int64, sex int32, nick, icon, bigIcon, phone string) bool {
	request := new(pgUser.ModifyUserInfoRequest)
	request.Mid = cid
	request.Icon = icon
	request.IconBig = icon
	request.Sex = sex
	request.Nick = nick
	request.Phone = phone
	reqBytes, err := request.Marshal()
	if err != nil {
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("pgUser.modifyUserInfo", reqBytes, 2000, 0)
	if response.ReturnCode != 0 {
		return false
	}
	pgResp := new(pgUser.ModifyUserInfoResponse)
	pgResp.Unmarshal(response.Body)
	plog.Debug("BaseLogic.modifyUser resp: ", pgResp)
	//成功
	if pgResp.Status == 0 {
		return true
	} else { //失败
		return false
	}
}

