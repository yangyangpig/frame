package service

import (
	"PGNotices/app/entity"
	"strconv"
	"time"
	"github.com/astaxie/beego"
)

type noticesService struct {}

// 表名
func (this *noticesService) table() string {
	return tableName("main.notices")
}

/**
 * db中获取公告列表
 * @param int32 appId 应用ID
 * @return []entry.Notices
 */
func (this *noticesService) GetList(appId int32) []entity.Notices {
	var noticesData []entity.Notices
	o.Using("default")
	now := time.Now().Unix()
	_, err := o.QueryTable(this.table()).Filter("app_id", appId).Filter("status", 1).Filter("start_time__lte",now).Filter("end_time__gte",now).Limit(30).OrderBy("-weight", "notice_id").All(&noticesData)
	if err != nil {
		beego.Error("[noticesService] GetList query db err:", err.Error())
		return []entity.Notices{}
	}
	return noticesData
}

/**
 * 从memcache中获取公告列表缓存
 * @param int32 appId 应用ID
 */
func (this *noticesService) GetMcNotices(appId int32) string {
	app := strconv.Itoa(int(appId))
	key := this.getKey(app)
	result, err := SystemCacheService.Get(key)
	if err != nil {
		beego.Error("[noticesService] GetMcNotices err:", err.Error())
		return ""
	}
	if v, ok := result.(string); ok {
		return v
	}
	return ""
}

/**
 * 从memcache中获取公告列表时间戳缓存
 * @param int32 appId 应用ID
 */
func (this *noticesService) GetMcVerNotices(appId int32) int64 {
	app := strconv.Itoa(int(appId))
	key := this.getVerKey(app)
	result, _ := SystemCacheService.Get(key)
	if v, ok := result.(int64); ok {
		return v
	}
	tm := time.Now().Unix()
	return tm
}

/**
 * 保存公告列表Id缓存到memcache
 * @param int32 appId 应用ID
 * @param string idStr 公告列表ID组
 */
func (this *noticesService) SetMcNotices(appId int32, idStr string) bool {
	app := strconv.Itoa(int(appId))
	key := this.getKey(app)
	ret := SystemCacheService.Set(key, idStr, 0)
	return ret
}

/**
 * 保存公告列表修改时间戳到memcache
 * @param int32 appId 应用ID
 */
func (this *noticesService) SetMcVerNotices(appId int32, now int64) bool {
	app := strconv.Itoa(int(appId))
	key := this.getVerKey(app)
	if now <= 0 {
		now = time.Now().Unix()
	}
	ret := SystemCacheService.Set(key, now, 0)
	return ret
}

/**
 * 获取缓存key
 */
func (this *noticesService) getKey(key string) string {
	return "NoticeMc" + key;
}

/**
 * 获取缓存版本key
 */
func (this *noticesService) getVerKey(key string) string {
	return "NoticeVerMc" + key;
}






