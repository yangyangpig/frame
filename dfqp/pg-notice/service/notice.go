package service

import (
	"dfqp/pg-notice/entity"
	"time"
	"github.com/astaxie/beego"
	"strconv"
)

type noticeService struct{}

//表名
func (this *noticeService) table() string {
	return tableName("main.notice")
}

/**
 * DB中获取公告列表
 * @param int32 appId 应用Id
 * @return []entry.Notice
 */
func (this *noticeService) GetList(appId int32) []entity.Notice {
	var noticeData []entity.Notice
	o.Using("default")
	now := time.Now().Unix()
	_, err := o.QueryTable(this.table()).Filter("app_id", appId).Filter("start_time__lte", now).Filter("end_time__gte", now).Limit(30).OrderBy("-weight", "notice_id").All(&noticeData)
	if err != nil {
		beego.Error("[noticeService] GetList query db err:", err.Error())
		return []entity.Notice{}
	}
	return noticeData
}

/**
 * 从memcache中获取公告列表的缓存
 * @param int32 appId 应用ID
 */
func (this *noticeService) GetMcNotice(appId int32) string {
	app := strconv.Itoa(int(appId))
	key := this.getKey(app)
	result,err := SystemCacheService.Get(key)
	if err != nil {
		beego.Error("[noticeService] GetMcNotice err:", err.Error())
		return ""
	}
	if v, ok := result.(string); ok {
		return v
	}
	return ""
}

/**
 * 获取缓存key
 */
 func (this *noticeService) getKey(key string) string {
 	return "NoticeMc" + key;
 }

 /**
  * 从memcache中获取公告列表时间戳缓存
  * @param int32 appid 应用ID
  */
  func (this *noticeService) GetMcVerNotice(appId int32) int64 {
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
   * 获取缓存版本key
   */
   func (this *noticeService) getVerKey(key string) string {
   		return "NoticeVerMc" + key;
   }