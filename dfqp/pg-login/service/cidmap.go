package service

import (
	"github.com/astaxie/beego/orm"
	"dfqp/pg-login/entity"
	"errors"
	"putil/log"
)

var CidMapService = &cidMapService{}

type cidMapService struct {}

var (
	paramError = errors.New("参数错误")
)

//表名
func (this *cidMapService) table(cid int64) string {
	return tableNameBy("users.cidmap", cid, 2)
}

//插入数据
//@param int64 cid 用户公共ID
//@param string platformId 平台帐号ID
//@param int32 platformType 平台帐号类型
func (this *cidMapService) Insert(cid int64, platformId string, platformType int32) bool {
	if cid <= 0 || len(platformId) == 0 {
		return false
	}
	o.Using("default")
	ret, err := o.Raw("INSERT INTO " + this.table(cid) + " (cid,platform_id,platform_type) VALUES(?, ?, ?)", cid, platformId, platformType).Exec()
	if err != nil {
		return false
	}
	id, err := ret.LastInsertId()
	if id <= 0 || err != nil {
		return false
	}
	return true
}
//@param string platformId 平台帐号ID
//@param int32 platformType 平台帐号类型
func (this *cidMapService) GetPlatformId(cid int64, platformType int32) (string, error) {
	if cid <= 0 {
		return "", paramError
	}
	ret := CidMapCacheService.Get(cid, platformType)
	if len(ret) > 0 {
		plog.Debug("cidMapCacheService.Get resp:", ret)
		return ret, nil
	} else {
		var data entity.Cidmap
		o.Using("default")
		//分表必须用原生查询
		err := o.Raw("SELECT platform_id FROM " + this.table(cid) + " WHERE cid = ? AND platform_type = ? Limit 1", cid, platformType).QueryRow(&data)
		plog.Debug("cidMapService.Get resp:", data, err)
		//没有找到记录
		if err == orm.ErrNoRows {
			return "", nil
		} else if err != nil {
			return "", err
		}
		go CidMapCacheService.Set(cid, platformType, data.PlatformId)
		return data.PlatformId, nil
	}
}



