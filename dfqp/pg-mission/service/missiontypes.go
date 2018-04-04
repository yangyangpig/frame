package service

import (
	"dfqp/pg-mission/entity"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/tidwall/gjson"
)

type missionTypesService struct{}

//表名
func (this *missionTypesService) table() string {
	return tableName("main.missiontypes")
}

/**
 * 根据appid查询任务列表数据
 * @param $regionId 地区id
 * @param $mid 	用户id
 * @param $status 状态
 */
func (this *missionTypesService) GetList(regionId int32, mid int64, status int32) ([]entity.Missiontypes, error) {
	orm.Debug = true
	var data []entity.Missiontypes
	var newData []entity.Missiontypes
	//var pogress string
	o.Using("default")
	_, err := o.Raw("SELECT id,region_id,name,`desc`,icon,reward,reward_type,sort_order,conditions,jump_code,status,circle_type FROM "+this.table()+" WHERE region_id = ?  AND status= ? ORDER BY id ", regionId, status).QueryRows(&data)

	for _, val := range data {
		if val.Conditions != "" {
			//t := missionRewardService.getStatus(val.Id, mid, val.Circle_type)
			//if t == 2 && val.Circle_type == 1 {
			//	//已完成的生涯任务不再显示
			//	continue
			//}
			//if t == 2 && val.Circle_type == 0 {
			//	//pogress = "已完成"
			//}
			m, ok := gjson.Parse(val.Conditions).Value().(map[string]interface{})
			if ok {
				for k, v := range (m) {
					//根据数据源获取任务进度
					switch k {
					case "totalPlayTimes": //玩牌局数
						fmt.Println("this is totalplayTimes task,condition:", v)
					case "bindPhoneCount": //绑定手机号
						fmt.Println("this is bindPhoneCount task,condition:", v)
					case "qzoneCount": //qq空间分享
						fmt.Println("this is qzoneCount task,condition:", v)
					case "shareCount": //分享次数
						fmt.Println("this is shareCount task,condition:", v)
					case "inviteCount": //邀请次数
						fmt.Println("this is inviteCount task,condition:", v)
					case "hookedUpFriends": //添加异性朋友
						fmt.Println("this is hookedUpFriends task,condition:", v)
					case "maxWinScore": //最高赢取
						fmt.Println("this is maxWinScore task,condition:", v)
					case "totalWinScore": //累计赢取
						fmt.Println("this is winTimes task,condition:", v)
					case "maxWinTimesInRow": //最高连胜
						fmt.Println("this is winTimes task,condition:", v)
					case "qicqCount": //QQ邀请
						fmt.Println("this is qicqCount task,condition:", v)
					}
				}
			}

			newData = append(newData, val)
		}
	}
	fmt.Println(newData)
	//没有找到记录
	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	return newData, err
}
