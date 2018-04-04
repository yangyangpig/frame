package service

import (
	"PGLogin/app/entity"
	"github.com/astaxie/beego/orm"
)

type midMapService struct {}

//表名
func (this *midMapService) table() string {
	return tableName("users.idmap")
}

//根据平台id和类型查询id
func (this *midMapService) GetMid(platformId string, platformType int32, region int32) (int64, error) {
	var data entity.Idmap
	o.Using("default")
	err := o.QueryTable(this.table()).Filter("platform_id", platformId).Filter("platform_type", platformType).Filter("region", region).One(&data, "Mid")
	//没有找到记录
	if err == orm.ErrNoRows {
		return 0, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return 0, err
	}
	return data.Mid, nil
}

//插入数据
//@param string platformId 平台帐号ID
//@param int32 platformType 平台帐号类型
//@param int32 region 地区
func (this *midMapService) Insert(platformId string, platformType int32, region int32) (int64, error) {
	o.Using("default")
	midMap := new(entity.Idmap)
	midMap.PlatformId = platformId
	midMap.PlatformType = platformType
	midMap.Region = region
	mid, err := o.Insert(midMap)
	if err != nil {
		return mid, err
	}
	return mid, nil
}