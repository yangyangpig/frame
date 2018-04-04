package service

import (
	"github.com/astaxie/beego/orm"
	"dfqp/pg-login/entity"
	"dfqp/lib"
	"fmt"
	"putil/log"
)

//实例化
var Platform2cidService = &platform2cidService{}

//平台id和标识和cid的映射
type platform2cidService struct {}

//获取表名
func (this *platform2cidService) table(platformId string, platformType int32) string {
	md5Str := lib.GetMd5(fmt.Sprintf("%s%d", platformId, platformType))
	return tableNameByMd5("users.platform2cid", md5Str, 2)
}

//根据平台id和类型查询cid
func (this *platform2cidService) GetCid(platformId string, platformType int32) (int64, error) {
	cid := Platform2cidCacheService.Get(platformId, platformType)
	if cid > 0 {
		plog.Debug("platform2cidCacheService.Get Resp:", cid)
		return cid, nil
	} else {
		var data entity.Platform2cid
		o.Using("default")
		//分表必须用原生查询
		err := o.Raw("SELECT cid FROM " + this.table(platformId, platformType) + " WHERE platform_id = ? and platform_type = ?", platformId, platformType).QueryRow(&data)
		plog.Debug("platform2cidService.Get Resp:", data, err)
		//没有找到记录
		if err == orm.ErrNoRows {
			return 0, nil
		} else if err != nil {
			return 0, err
		}
		go Platform2cidCacheService.Set(platformId, platformType, data.Cid)
		return data.Cid, nil
	}
}

//插入数据
func (this *platform2cidService) Insert(cid int64, platformId string, platformType int32) bool {
	if cid <= 0 || len(platformId) == 0 {
		return false
	}
	o.Using("default")
	ret, err := o.Raw("INSERT IGNORE INTO " + this.table(platformId, platformType) + " (cid,platform_id,platform_type) VALUES(?, ?, ?)", cid, platformId, platformType).Exec()
	if err != nil {
		return false
	}
	id, err := ret.LastInsertId()
	if id <= 0 || err != nil {
		return false
	}
	return true
}
