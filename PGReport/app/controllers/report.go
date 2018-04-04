package controllers

import (
	"github.com/astaxie/beego"
	"PGReport/app/proto"
	"putil/log"
	"framework/rpcclient/szmq"
	"encoding/json"
	"time"
)

// 框架必备
type ReportController struct {
	beego.Controller
}

// 定义上报数据结构，主要用于转换成json
type DevicesUpload struct
{
	Platform string `json:"platform"`
	Osversion string `json:"osversion"`
	DeviceId string `json:"deviceId"`
	Identifier string `json:"identifier"`
	Brand string `json:"brand"`
	BundleVersion string `json:"bundleVersion"`
	BundleShortVersion string `json:"bundleShortVersion"`
	SerialId string `json:"serialId"`
	DeviceName string `json:"deviceName"`
	Model string `json:"model"`
	Manufacturer string `json:"manufacturer"`
	Locale string `json:"locale"`
	Countrt string `json:"countrt"`
	Timezone string `json:"timezone"`
	Acttime int64 `json:"act_time"`
}

// 上报设备信息
func (n *ReportController) Devices(rq *Report.DevicesRequest) *Report.DevicesResponse {
	plog.Debug("request")
	plog.Debug(rq)
	rp := new(Report.DevicesResponse)

	// 处理上报格式
	var uploadData DevicesUpload
	uploadData.Acttime = time.Now().UnixNano()/1e6 // 取得服务器当前时间戳，毫秒
	uploadData.Platform = rq.Platform
	uploadData.Osversion = rq.Osversion
	uploadData.DeviceId = rq.DeviceId
	uploadData.Identifier = rq.Identifier
	uploadData.Brand = rq.Brand
	uploadData.BundleVersion = rq.BundleVersion

	uploadData.BundleShortVersion = rq.BundleShortVersion
	uploadData.SerialId = rq.SerialId
	uploadData.DeviceName = rq.DeviceName
	uploadData.Model = rq.Model
	uploadData.Manufacturer = rq.Manufacturer
	uploadData.Locale = rq.Locale

	uploadData.Countrt = rq.Countrt
	uploadData.Timezone = rq.Timezone

	uploadByte, err := json.Marshal(uploadData)
	if err != nil {
		plog.Debug(err)
	}
	uploadStr := string(uploadByte[:])

	plog.Debug(uploadStr)
	go szmq.Logger.WriteNormalLog("dfqp_devices", uploadStr)	//WriteNormalLog为普通日志上报

	result := "1"
	rp.Result = result
	return rp
}