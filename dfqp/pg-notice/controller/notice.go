package controller

import (
	"fmt"
	"dfqp/proto/notice"
	"math"
	"time"
	"encoding/json"
	"strconv"
	"dfqp/lib"
	"dfqp/pg-notice/service"
	"dfqp/pg-notice/entity"
	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
	"strings"
	"dfqp/lib/ipdata"
	//"framework/rpcclient/core"
	"putil/log"
)

type NoticeController struct {}

/**
 * 获取公告列表
 */

func (n *NoticeController) NoticeList(rq pgNotice.GetListRequest) *pgNotice.GetListResponse {
	plog.Debug("BobinSun--------noticeController")
	appid := rq.App
	mid := rq.Mid
	hallver := rq.HallVersion
	version := lib.Ver2long(rq.Version)
	areaId := rq.AreaId
	cli_ver := rq.CliVer
	svr_ver := cli_ver
	region := lib.GetRegionId(appid)

	rp := new(pgNotice.GetListResponse)
	var (
		err error
		result []pgNotice.GetListResData
		noticeData []entity.Notice
		idList []interface{}
		index int
		isTrue bool
	)
	week := int(time.Now().Weekday())
	now := time.Now().Unix()
	dt := time.Now().Format("2006-01-02")
	nkey := fmt.Sprintf("goni%s%d", dt, appid)
	//从memcache中获取缓存
	cacheData, _ := service.SystemCacheService.Get(nkey)
	isTrue = true
	if v, ok := cacheData.(string); ok {
		err = json.Unmarshal([]byte(v),&noticeData)
		if err != nil {
			isTrue = false
			beego.Error("unmarshal memcache data err: ", err.Error())
		}
	} else {
		isTrue = false
	}
	if !isTrue {
		noticeData = service.NoticeService.GetList(appid)
		if len(noticeData) > 0 {
			noticeStr, err := json.Marshal(noticeData)
			if err != nil {
				beego.Error("marshal notice data err: ", err.Error())
			} else {
				service.SystemCacheService.Set(nkey,noticeStr, 7030)
			}
		}
	}

	if len(noticeData) > 0 {
		for _, value := range noticeData {
			mids := gjson.Parse(value.Mids).Array()
			beego.Debug("mids: ",mids)
			if len(mids) > 0 {
				isTrue = false
				for _, m := range mids {
					if mid == m.Int() {
						isTrue = true
					}
				}
				if !isTrue {
					continue
				}
			}
			conditions := gjson.Parse(value.Conditions)
			//最低版本号判断
			minVersion := conditions.Get("min_version").String()
			if len(minVersion) > 0 {
				index = strings.Index(minVersion, ".")
				if index > 0 {
					beego.Debug("ver2long minversion: ", lib.Ver2long(minVersion))
					if version < lib.Ver2long(minVersion) {
						continue
					}
				} else {
					ver, err := strconv.Atoi(minVersion)
					if err != nil {
						beego.Debug("atoi minVersion fail: ", err.Error())
						continue
					}
					if ver > 0 && hallver < int64(ver) {
						continue
					}
				}
			}
			//最高版本号判断
			maxVersion := conditions.Get("max_version").String()
			if len(maxVersion) > 0 {
				index = strings.Index(minVersion, ".")
				if index > 0 {
					beego.Debug("ver2long maxversion: ", lib.Ver2long(maxVersion))
					if version < lib.Ver2long(maxVersion) {
						continue
					}
				} else {
					ver, err := strconv.Atoi(maxVersion)
					if err != nil {
						beego.Debug("atoi maxVersion fail: ", err.Error())
						continue
					}
					if ver > 0 && hallver < int64(ver) {
						continue
					}
				}
			}
			//地区ID判断
			if areaId > 0 {
				cities := conditions.Get("cities").Array()
				beego.Debug("cities value is: ", cities)
				if len(cities) > 0 {
					isTrue = false
					for _, vv := range cities {
						if vv.Int() == int64(areaId) {
							isTrue = true
							break
						}
					}
					if !isTrue {
						continue
					}
				}
			}
			//信誉值积分判断
			isLogined := conditions.Get("isLogined").Int()
			trustvalue := conditions.Get("trustvalue").Int()
			if isLogined == 1 && trustvalue > 0 && mid > 0 {
				positiveScore := service.UserScoreService.GetNegativeScore(mid)
				beego.Debug("positiveScore is: ", positiveScore)
				if positiveScore < trustvalue {
					continue
				}
			}
			//判断是否是本省ip
			ischeckip := conditions.Get("ischeckip").Int()
			if ischeckip == 1 {
				if !n.isRegionIp(region) {
					continue
				}
			}
			//兼容老版本时间段显示，900以上的都在客户端判断
			weekArr := conditions.Get("week").Array()
			var weekTmp []int32
			beego.Debug("week value is: ", weekArr)
			if len(weekArr) > 0 {
				isTrue = false
				for _, wk := range weekArr {
					weekTmp = append(weekTmp, int32(wk.Int()))
					if wk.Int() == int64(week) {
						isTrue = true
					}
				}
				if !isTrue && (hallver < 900) {
					continue
				}
			}
			var pertimeTmp []pgNotice.GetListPertime
			pertimeArr := conditions.Get("pertime").Array()
			timeflag := 0
			beego.Debug("pertime value is: ", pertimeArr)
			if len(pertimeArr) > 0 {
				for _, v := range pertimeArr {
					stime := v.Get("stime").String()
					etime := v.Get("etime").String()
					pertimeTmp = append(pertimeTmp, pgNotice.GetListPertime{
						Stime: stime,
						Etime: etime,
					})
					if hallver < 900 {
						if len(stime) > 0 {
							tm, err := time.Parse("2006-01-02 15:04:05", dt + " " + stime)
							if err != nil {
								beego.Error("stime parse err: ",err.Error())
								continue
							}
							sTime := tm.Unix()
							beego.Debug("stime value is: ", sTime)
							if sTime < 0 || sTime > now {
								continue
							}
						}
						if len(etime) > 0 {
							tm, err := time.Parse("2006-01-02 15:04:05", dt + " " + stime)
							if err != nil {
								beego.Error("etime parse err: ",err.Error())
								continue
							}
							eTime := tm.Unix()
							beego.Debug("etime value is: ", eTime)
							if eTime < 0 || now > eTime {
								continue
							}
							timeflag = 1
							break
						}
					}
					if timeflag == 0 && hallver < 900 {
						continue
					}
				}
				//公告类型为富文本公告时
				if value.NoticeType == 1 {
					urlPrefix := beego.AppConfig.String("site_url")
					beego.Debug("urlPrefix", urlPrefix)
					url := fmt.Sprintf("%s?action=notice.getdetail%vid=%d%vapp%d", urlPrefix, "&", value.NoticeId, "&", value.AppId)
					value.Content = url
				}

				idList = append(idList, value.NoticeId)
				if hallver < 900 {
					stime, err := strconv.ParseInt(value.StartTime, 10, 64)
					if err != nil {
						continue
					}
					tm := time.Unix(stime, 0)
					nstime := tm.Format("2006/01/02")
					value.StartTime = nstime
				}
				sendtype := int32(conditions.Get("sendtype").Int())
				poptype := int32(conditions.Get("poptype").Int())
				result = append(result, pgNotice.GetListResData{
					NoticeId:value.NoticeId,
					NoticeType:value.NoticeType,
					AppId:appid,
					Weight:value.Weight,
					Title:value.Title,
					Content:value.Content,
					StartTime:value.StartTime,
					EndTime:value.EndTime,
					Conditions:pgNotice.GetListConditions{
						Sendtype:sendtype,
						Poptype:poptype,
						//IsLogined: int32(isLogined),
						//Week:weekTmp,
						//Pertime:pertimeTmp,
					},
				})
			}
		}
	}
	svrVerCache := service.NoticeService.GetMcVerNotice(appid)
	svr_ver = int64(math.Max(float64(svr_ver), float64(svrVerCache)))
	if int64(cli_ver) < svr_ver {
		rp.Isrefresh = 1
		//if hallver < 900 {
		//	for _, v := range idList {
		//		if vv, ok := v.(int32); ok {
		//			rp.Idlist = append(rp.Idlist, vv)
		//		}
		//	}
		//}
		rp.SvrVer = svr_ver
		//rp.Data = result
	} else {
		rp.Isrefresh = 0
		//if hallver < 900 {
		//	for _, v := range idList {
		//		if vv, ok := v.(int32); ok {
		//			rp.Idlist = append(rp.Idlist, vv)
		//		}
		//	}
		//}
		rp.SvrVer = svr_ver
		//rp.Data = []pgNotice.GetListResData{}
	}
	rp.Svrtime = now
	plog.Debug("BobinSun=========================")
	//return new(pgNotice.GetListResponse)
	return &pgNotice.GetListResponse{}
}

//是否是本省ip
func (n *NoticeController) isRegionIp(region int32) bool {
	ipObj := ipdata.NewIpdata()
	ip := ""
	ipInfo, err := ipObj.Find(ip)
	if err != nil {
		return false
	}
	if _, ok := ipInfo["province"]; ok {
		provinceData := service.RegionService.GetProvince()
		jsonData := gjson.Parse(provinceData)
		city2Province := jsonData.Get("city2province")
		province := jsonData.Get("province")
		newRegion := strconv.Itoa(int(region))
		nameZh := province.Get(newRegion).Get("name_zh").String()
		var pid string
		//说明是省份id
		if len(nameZh) <= 0 {
			pid = city2Province.Get(newRegion).String()
			nameZh = province.Get(pid).Get("name_zh").String()
		}
		if nameZh == "其他" {
			cityNameZh := province.Get(pid).Get("cities").Get(newRegion).Get("name_zh").String()
			return strings.Contains(cityNameZh, "四川")
		} else {
			return strings.Contains(nameZh, "江西")
		}
	}
	return false
}