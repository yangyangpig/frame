package service

import (
	"PGLogin/app/entity"
	"github.com/astaxie/beego/orm"
)

type cidMapService struct {}

//表名
func (this *cidMapService) table() string {
	return tableName("common.cidmap")
}

//根据平台id和类型查询cid
func (this *cidMapService) GetCid(platformId string, platformType int32) (int64, error) {
	var data entity.Cidmap
	o.Using("common")
	err := o.QueryTable(this.table()).Filter("platform_id", platformId).Filter("platform_type", platformType).One(&data, "Cid")
	//没有找到记录
	if err == orm.ErrNoRows {
		return 0, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return 0, err
	}
	return data.Cid, nil
}

//插入数据
//@param string platformId 平台帐号ID
//@param int32 platformType 平台帐号类型
func (this *cidMapService) Insert(platformId string, platformType int32) (int64, error) {
	o.Using("common")
	cidMap := new(entity.Cidmap)
	cidMap.PlatformId = platformId
	cidMap.PlatformType = platformType
	cid, err := o.Insert(cidMap)
	if err != nil {
		return cid, err
	}
	return cid, nil
}



