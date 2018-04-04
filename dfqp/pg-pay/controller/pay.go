package controller

import (
	"dfqp/proto/pay"
	//"dfqp/lang/zh"
	"crypto/sha1"
	"dfqp/lib"
	"dfqp/pg-pay/entity"
	"dfqp/pg-pay/service"
	"dfqp/proto/online"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"putil/log"
	"strconv"
	"strings"
	"time"
)

type PayController struct{}

func (this *PayController) Config(rq *pgPay.ConfigRequest) *pgPay.ConfigResponse {
	resp := new(pgPay.ConfigResponse)
	mid := rq.GetMid()
	tabId := rq.GetTabId()
	if mid <= 0 || tabId < 0 {
		resp.Status = 1
		resp.Msg = "请求失败！"
		return resp
	}

	// 取得tab配置
	var (
		tabConf   []pgPay.TabRespData
		goodsConf []pgPay.GoodsRespData
	)
	tabConfJson := service.GoodsService.GetTabConf()
	tabEndResult := gjson.Parse(tabConfJson)
	for _, v := range tabEndResult.Array() {
		itemStr := v.String()
		itemName := gjson.Get(itemStr, "name").String()
		itemId := gjson.Get(itemStr, "id").Int()
		itemOrder := gjson.Get(itemStr, "order").Int()
		tabConfItem := pgPay.TabRespData{
			Id:    int32(itemId),
			Name:  itemName,
			Order: int32(itemOrder),
		}
		tabConf = append(tabConf, tabConfItem)
	}
	resp.Data.TabConf = tabConf

	if tabId == 0 {
		tabId = 1
	}
	var clientId int32 = 3
	goodsConfMap := service.GoodsService.GetGoodsConf(tabId, clientId)
	//if len(goodsConfMap) == 0 {
	//	resp.Status = 1
	//	resp.Msg = zh.Response_Err
	//	return resp
	//}
	for _, v := range goodsConfMap {
		goodsItemStr := gjson.Parse(v)
		goodsId := goodsItemStr.Get("gid").Int()
		goodsName := goodsItemStr.Get("name").String()
		goodsIcon := goodsItemStr.Get("icon").String()
		goodsOrder := goodsItemStr.Get("order").Int()
		//goodsFlagIcon := goodsItemStr.Get("flagIcon").String()
		//goodsDesc := goodsItemStr.Get("desc").String()
		limitNum := goodsItemStr.Get("limitNum").Int()
		limitTotal := goodsItemStr.Get("limitTotal").Int()
		sellId := goodsItemStr.Get("sellId").Int()
		sellNum := goodsItemStr.Get("sellNum").Int()

		var (
			descConf   []pgPay.DescRespData
			payConf    []pgPay.PayRespData
			detailConf []pgPay.GiftRespDetail
		)

		descConfJson := goodsItemStr.Get("descConf")
		for _, gv := range descConfJson.Array() {
			descItemStr := gv.String()
			descType := gjson.Get(descItemStr, "type").Int()
			descFlagIcon := gjson.Get(descItemStr, "flagIcon").String()
			descDesc := gjson.Get(descItemStr, "desc").String()
			descOrder := gjson.Get(descItemStr, "dorder").Int()
			descDetailJson := gjson.Get(descItemStr, "detail").String()
			detailConfJson := gjson.Parse(descDetailJson)

			for _, dv := range detailConfJson.Array() {
				detailItemStr := dv.String()
				giftGoodsId := gjson.Get(detailItemStr, "goodsId").Int()
				giftGoodsNum := gjson.Get(detailItemStr, "goodsNum").Int()
				detailItem := pgPay.GiftRespDetail{
					GoodsId:  int32(giftGoodsId),
					GoodsNum: int32(giftGoodsNum),
				}
				detailConf = append(detailConf, detailItem)
			}

			descItem := pgPay.DescRespData{
				Type:     int32(descType),
				FlagIcon: descFlagIcon,
				Desc:     descDesc,
				Dorder:   int32(descOrder),
				Detail:   detailConf,
			}

			//giftGoodsId := gjson.Get(giftItemStr, "goodsId").Int()
			//giftGoodsNum := gjson.Get(giftItemStr, "goodsNum").Int()
			//giftItem := pgPay.GiftRespDetail{
			//	GoodsId:  int32(giftGoodsId),
			//	GoodsNum: int32(giftGoodsNum),
			//}
			descConf = append(descConf, descItem)
		}

		payJson := goodsItemStr.Get("pay")

		for _, pv := range payJson.Array() {
			payItemStr := pv.String()
			payId := gjson.Get(payItemStr, "payId").Int()
			oldPayNum := gjson.Get(payItemStr, "oldPayNum").Float()
			payNum := gjson.Get(payItemStr, "payNum").Float()
			pmode := gjson.Get(payItemStr, "pmode").String()
			payItem := pgPay.PayRespData{
				PayId:     int32(payId),
				OldPayNum: float32(oldPayNum),
				PayNum:    float32(payNum),
				Pmode:     pmode,
			}
			payConf = append(payConf, payItem)
		}

		goodsConfItem := pgPay.GoodsRespData{
			Gid:        int32(goodsId),
			Name:       goodsName,
			Icon:       goodsIcon,
			Order:      int32(goodsOrder),
			Desc:       descConf,
			Pay:        payConf,
			LimitNum:   int32(limitNum),
			LimitTotal: int32(limitTotal),
			SellId:     int32(sellId),
			SellNum:    int32(sellNum),
		}

		goodsConf = append(goodsConf, goodsConfItem)
	}
	resp.Data.Goods = goodsConf
	resp.Data.TabId = tabId
	//成功
	resp.Status = 0
	resp.Msg = "请求成功"
	return resp
}

func (this *PayController) Order(rq *pgPay.OrderRequest) *pgPay.OrderResponse {
	resp := new(pgPay.OrderResponse)
	pmode := rq.GetPmode()
	mid := rq.GetMid()
	price := rq.GetPrice()
	tabId := rq.GetTabId()
	//number := rq.GetNumber()
	gid := rq.GetGid()

	if mid <= 0 || tabId <= 0 || gid <= 0 || price <= 0 || pmode <= 0 {
		resp.Status = 1
		resp.Msg = "请求参数有误！"
		return resp
	}

	// 判断是否已经登录
	onlineInfo := this.getOnlineInfo(mid)
	if onlineInfo.Retcode == 1 {
		resp.Status = 1
		resp.Msg = "请先登录！"
		return resp
	}

	// 固定参数
	rurl := "http://bypaycn.boyaa.com/Pay/unifiedOrder/" //"http://bypaycn-debug.oa.com/Pay/unifiedOrder/"
	pamountUnit := "CNY"
	time := time.Now().Unix()
	userIp := onlineInfo.Ip

	// 取得下单信息，转换成数据中心的sid和appid
	cityAppId := onlineInfo.CityApp
	clientType := lib.GetClientId(onlineInfo.CityApp)
	payCenterConf := service.GoodsService.GetPayCenterConf(cityAppId)
	plog.Debug("appid:=======>", cityAppId)
	plog.Debug("conf:=======>", payCenterConf)
	payCenterArr := gjson.Parse(payCenterConf)
	sid := payCenterArr.Get("sid").Int()
	payAppId := payCenterArr.Get("appid").Int()
	signKey := payCenterArr.Get("sign").String()

	// 校验物品信息
	goodsInfoString := service.GoodsService.GetGoodsInfo(int32(tabId), clientType, gid)
	plog.Debug("goods info:", goodsInfoString)
	if goodsInfoString == "" {
		resp.Status = 3
		resp.Msg = "商品ID不存在！"
		return resp
	}
	goodsInfo := gjson.Parse(goodsInfoString)
	goodsName := goodsInfo.Get("name").String()
	goodsPayInfo := goodsInfo.Get("pay")
	// 校验是否支持支付方式
	priceFlag := false
	plog.Debug("goods pay info:", goodsPayInfo)
	for _, v := range goodsPayInfo.Array() {
		payId := v.Get("payId").Int()
		payNum := v.Get("payNum").Float()
		payPmode := v.Get("pmode").String()
		pmodeSplit := strings.Split(payPmode, ";")
		pmodeFlag := false
		for _, pv := range pmodeSplit {
			pmodeNum, _ := strconv.Atoi(pv)
			plog.Debug("pmode:======>", pmodeNum)
			if pmode == int32(pmodeNum) {
				pmodeFlag = true
			}
		}
		plog.Debug("pay pmode:", payPmode)
		plog.Debug("pmode split:", pmodeSplit)
		if payId == -1 && float64(price) == payNum && pmodeFlag {
			priceFlag = true
			break
		}
	}
	if priceFlag == false {
		resp.Status = 3
		resp.Msg = "商品ID与价格或支付方式不匹配！"
		return resp
	}

	shaStr := fmt.Sprintf("appid=%ditem_id=%ditem_name=%smid=%dpamount=%.2fpamount_unit=%spmode=%dsid=%dsitemid=%dtime=%duserip=%s%s", payAppId, gid, goodsName, mid, price, pamountUnit, pmode, sid, mid, time, userIp, signKey)
	plog.Debug("sha================>", shaStr)
	sign := this.getShaValue(shaStr)
	u, _ := url.Parse(rurl)
	q := u.Query()
	q.Set("appid", fmt.Sprintf("%d", payAppId))
	q.Set("pmode", fmt.Sprintf("%d", pmode))
	q.Set("sid", fmt.Sprintf("%d", sid))
	q.Set("sitemid", strconv.FormatInt(mid, 10))
	q.Set("mid", strconv.FormatInt(mid, 10))
	q.Set("pamount", fmt.Sprintf("%.2f", price))
	q.Set("pamount_unit", pamountUnit)
	q.Set("item_id", fmt.Sprintf("%d", gid))
	q.Set("item_name", goodsName)
	q.Set("userip", userIp)
	q.Set("time", strconv.FormatInt(time, 10))
	q.Set("sign", sign)
	u.RawQuery = q.Encode()

	payResp, _ := http.Get(u.String())
	defer payResp.Body.Close()
	callBack, _ := ioutil.ReadAll(payResp.Body)
	plog.Debug("order================>", string(callBack))
	callBackData := gjson.Parse(string(callBack))
	callRet := callBackData.Get("RET").Int()
	if callRet != 0 {
		resp.Status = 3
		callMsg := callBackData.Get("MSG").String()
		resp.Msg = callMsg
		return resp
	}
	resp.Status = 0
	resp.Msg = "下单成功！"
	resp.Data.Pid = callBackData.Get("PID").String()
	resp.Data.Pmode = int32(callBackData.Get("PMODE").Int())
	resp.Data.Mid = int32(callBackData.Get("MID").Int())
	resp.Data.Order = callBackData.Get("ORDER").String()
	resp.Data.Gid = gid
	resp.Data.Pname = goodsName
	resp.Data.Pamount = price
	switch resp.Data.Pmode {
	case 265:
		returnObj := new(entity.AliPayParam)
		returnObj.Pname = goodsName
		returnObj.Pamount = fmt.Sprintf("%0.2f", price)
		returnObj.NotifyUrl = callBackData.Get("NOTIFY_URL").String()
		returnObj.Porder = resp.Data.Order
		alipayConf := service.ServiceConf["pay"]
		plog.Debug(alipayConf)
		returnObj.Partner = alipayConf.String("alipay.seller")
		returnObj.Seller = returnObj.Partner
		returnObj.RsaPrivate = alipayConf.String("alipay.rsa_private")
		returnObj.Udesc = fmt.Sprintf("%0.2f元=%s", price, goodsName)
		extBytes, _ := json.Marshal(returnObj)
		resp.Data.Ext = string(extBytes)
	case 198:
		returnObj := new(entity.UnionPayParam)
		returnObj.Tn = callBackData.Get("tn").String()
		returnObj.Pmode = fmt.Sprintf("%d", resp.Data.Pmode)
		extBytes, _ := json.Marshal(returnObj)
		resp.Data.Ext = string(extBytes)
	case 99:
		returnObj := new(entity.ApplePayParam)
		returnObj.Order = resp.Data.Order
		returnObj.Appstoreid = "com.boyaa.LNQP.180000silver_Tier3" //callBackData.Get("appstoreid").String()
		extBytes, _ := json.Marshal(returnObj)
		resp.Data.Ext = string(extBytes)
	case 431:
		returnObj := new(entity.WeChatParam)
		returnObj.Pmode = fmt.Sprintf("%d", resp.Data.Pmode)
		returnObj.TimeStamp = callBackData.Get("timestamp").String()
		returnObj.PartnerId = callBackData.Get("partnerid").String()
		returnObj.PrepayId = callBackData.Get("prepayid").String()
		returnObj.NonceStr = callBackData.Get("noncestr").String()
		returnObj.PackageValue = callBackData.Get("package").String()
		returnObj.Sign = callBackData.Get("sign").String()
		returnObj.ExtData = resp.Data.Order
		extBytes, _ := json.Marshal(returnObj)
		resp.Data.Ext = string(extBytes)
	}

	return resp

}

func (this *PayController) Report(rq *pgPay.ReportRequest) *pgPay.ReportResponse {
	resp := new(pgPay.ReportResponse)

	pid := rq.GetPid()
	content := rq.GetContent()
	detail := gjson.Parse(content)

	pdealno := detail.Get("pdealno").String()
	receipt := detail.Get("receipt").String()
	sandbox := detail.Get("sandbox").Int()
	bundleId := detail.Get("bundleId").String()

	// 计算签名
	key := "sdajir4%^&@RR@E@Eff0xi-bb6767"
	time := time.Now().Unix()
	preSignStr := fmt.Sprintf("%s%s%d%s", pid, pdealno, time, key)
	sign := lib.GetMd5(preSignStr)

	rurl := "http://bypaycn.boyaa.com/Iphone/callback/"
	u, _ := url.Parse(rurl)
	q := u.Query()
	q.Set("pid", pid)
	q.Set("pdealno", pdealno)
	q.Set("receipt", receipt)
	q.Set("time", fmt.Sprintf("%d", time))
	q.Set("sign", sign)
	if sandbox == 1 { //沙盒账号
		q.Set("test", "test")
	}
	if bundleId != "" { //马甲版本
		q.Set("bypayBid", bundleId)
	}

	u.RawQuery = q.Encode()

	payResp, _ := http.Get(u.String())
	defer payResp.Body.Close()
	callBack, _ := ioutil.ReadAll(payResp.Body)
	callBackData := gjson.Parse(string(callBack))
	callRet := callBackData.Get("ErrorCode").Int()
	switch callRet {
	case 1:
		resp.Status = 0
		resp.Msg = "请求成功！"
	case 6:
		resp.Status = 1
		resp.Msg = "客户端请求超时"
	default:
		resp.Status = 2
		resp.Msg = "请求失败！"
	}
	return resp
}

func (this *PayController) getOnlineInfo(mid int64) *pgOnline.GetOnlineResponse {
	//onlineRequest := new(pgOnline.GetOnlineRequest)
	//onlineRequest.Uid = mid
	//reqBytes, _ := onlineRequest.Marshal()
	onlineResponse := new(pgOnline.GetOnlineResponse)
	//callBack := service.Client.SendAndRecvRespRpcMsg("online.OnlineService.Get", reqBytes, 5000, 0)
	//onlineResponse.Unmarshal(callBack.Body)
	onlineResponse.Ip = "172.20.134.10"
	onlineResponse.CityApp = 203001
	onlineResponse.Retcode = 0
	return onlineResponse
}

func (this *PayController) getShaValue(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	shaV := h.Sum(nil)
	return fmt.Sprintf("%x", shaV)
}

/************************** 发货模块 *********************************/
var (
	rqParam = []string{
		"apple_productid", "do", "ext", "game_item_id", "mid", "pamount",
		"pamount_change", "pamount_rate", "pamount_unit", "pamount_usd", "paychips_v2", "paycoins",
		"payconfid", "payprod_v2", "pc_appid", "pc_rate", "pc_sid", "pc_time",
		"pdealno", "pendtime", "pid", "pmode", "pnum_v2", "pstarttime",
		"sign", "sign_v2", "sign_v3", "sitemid", "time",
	}
)

func (this *PayController) SendMoney(rq *pgPay.SendMoneyRequest) *pgPay.SendMoneyResponse {
	plog.Debug("what the fuck !!!", rq)
	resp := new(pgPay.SendMoneyResponse)

	resp.Status = 0
	resp.Msg = rq.Content
	return resp
	//do := "2222"
	//switch do {
	//case "suc":
	//
	//case "refund":
	//
	//case "sendcbk":
	//
	//default:
	//
	//}

	//sort.Strings(rqParam)
	//var (
	//	signBuff bytes.Buffer
	//)
	//fmt.Println(rq.Form)
	//for _, v := range rqParam {
	//	rqData[v] = rq.Form.Get(v)
	//	if v != "sign_v3" {
	//		signBuff.WriteString(fmt.Sprintf("%s=%s&", v, rqData[v]))
	//	}
	//}
	//signConf := "appkey=DCD023436FC06A91F28C79E64052A925" // 阜新安卓
	//signBuff.WriteString(signConf)
	//signV3 := lib.GetMd5(signBuff.String())
	//if signV3 == rqData["sign_v3"] {
	//	signFlag = true
	//}
	//fmt.Println(signBuff.String(), signV3, rqData["sign_v3"])
}

// 判断订单是否已经发货
func (this *PayController) checkSendStatus(pid int) bool {

	return false
}

// 发货
func (this *PayController) suc(writer http.ResponseWriter, request *http.Request) {

}

// 扣回通知
func (this *PayController) refund(writer http.ResponseWriter, request *http.Request) {

}

// 扣费失败回调
func (this *PayController) sendcbk(writer http.ResponseWriter, request *http.Request) {

}
