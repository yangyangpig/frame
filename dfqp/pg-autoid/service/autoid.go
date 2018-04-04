package service

import (
	"dfqp/pg-autoid/entity"

	"fmt"
	"putil/log"
	"time"

	"github.com/astaxie/beego/orm"
)

type autoidService struct{}

//表名
func (this *autoidService) table() string {
	return tableName("common.autoid_source")
}

//查询所有的btag
func (this *autoidService) GetBtags() []entity.AutoidSource {
	var data []entity.AutoidSource
	o.Using("default")
	//plog.Debug("GetMaxAndStep  GetMaxAndStep  GetMaxAndStep", btag)
	_, err := o.QueryTable(this.table()).All(&data, "btag", "step")
	//没有找到记录
	if err == orm.ErrNoRows {
		plog.Fatal("[autoidService] GetBtags query db err:", err.Error())
		return nil
	}
	if err != nil && err != orm.ErrNoRows {
		plog.Fatal("[autoidService] GetBtags query db err:", err.Error())
		return nil
	}
	return data

	//var noticesData []entity.Notices
	//	o.Using("default")
	//	now := time.Now().Unix()
	//	_, err := o.QueryTable(this.table()).Filter("app_id", appId).Filter("status", 1).Filter("start_time__lte",now).Filter("end_time__gte",now).Limit(30).OrderBy("-weight", "notice_id").All(&noticesData)
	//	if err != nil {
	//		beego.Error("[noticesService] GetList query db err:", err.Error())
	//		return []entity.Notices{}
	//	}
	//	return noticesData
}

//根据btag查询当前maxid和step
func (this *autoidService) GetMaxAndStep(btag string) (int64, int32, error) {
	var data entity.AutoidSource
	o.Using("default")
	plog.Debug("GetMaxAndStep  GetMaxAndStep  GetMaxAndStep", btag)
	err := o.QueryTable(this.table()).Filter("btag", btag).One(&data, "max_id", "step")
	//没有找到记录
	if err == orm.ErrNoRows {
		return 0, 0, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return 0, 0, err
	}
	return data.Maxid, data.Step, nil
}

//更新Max值(用当前max加上step)
func (this *autoidService) ModifyAndGet(btag string) (int64, int32, error) {

	o.Begin()

	sql := "Update " + this.table() + " SET max_id = max_id + step ,update_time = " + fmt.Sprintf("%d", time.Now().Unix()) + " WHERE btag = '" + btag + "'"
	_, err := o.Raw(sql).Exec()
	if err != nil {
		o.Rollback()
		return 0, 0, err
	}
	var data entity.AutoidSource
	err = o.QueryTable(this.table()).Filter("btag", btag).One(&data, "max_id", "step")
	//没有找到记录
	if err == orm.ErrNoRows {
		o.Rollback()
		return 0, 0, nil
	}
	if err != nil && err != orm.ErrNoRows {
		o.Rollback()
		return 0, 0, err
	}
	o.Commit()
	return data.Maxid, data.Step, nil

	//	num, err := o.QueryTable(this.table()).Update(orm.Params{
	//		"max_id": orm.ColValue(orm.ColAdd, "step"),
	//	})
	//	plog.Debug(num, err)

	//	resp := new(pgUser.UpdateUserResponse)
	//	if rq.Mid <= 0 {
	//		resp.Status = 1
	//		return resp
	//	}
	//	var (
	//		userInfo entity.User
	//		err      error
	//	)
	//	userInfo, err = service.UserService.Get(rq.Mid)
	//	if err != nil {
	//		resp.Status = 1
	//		return resp
	//	}
	//	//获取cid
	//	cid, _ := service.CidMapService.GetCid(userInfo.Platform_id, userInfo.Platform_type)
	//	cuser := new(entity.Cuser)
	//	cuser.Cid = cid
	//	cuser.Nick = rq.Nick
	//	cuser.City = rq.City
	//	cuser.Icon = rq.Icon
	//	cuser.Icon_big = rq.IconBig
	//	cuser.Sex = rq.Sex
	//	ret := service.CuserService.Update(cuser)
	//	if !ret {
	//		//失败
	//		resp.Status = 1
	//	}
	//	return resp
}

//插入数据
//@param string platformId 平台帐号ID
//@param int32 platformType 平台帐号类型
//func (this *autoidService) Insert(platformId string, platformType int32) (int64, error) {
//	o.Using("common")
//	cidMap := new(entity.Cidmap)
//	cidMap.PlatformId = platformId
//	cidMap.PlatformType = platformType
//	cid, err := o.Insert(cidMap)
//	if err != nil {
//		return cid, err
//	}
//	return cid, nil
//}
