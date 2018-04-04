package service

import (
	"dfqp/pg-mission/entity"
	"fmt"
	"github.com/astaxie/beego/orm"
	//"github.com/tidwall/gjson"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"github.com/garyburd/redigo/redis"
)

type missionRewardService struct{}

//表名
func (this *missionRewardService) table(circle_type int32) string {
	if circle_type == 0 {
		return tableNameByDate("logs.missionreward") //每日任务
	} else {
		return tableName("logs.missionreward") //生涯任务
	}

}

//获取任务key
func (this *missionRewardService) getKey(mid int64, circle_type int32) (string) {
	var key string
	if circle_type == 0 {
		key = "daymission|" + strconv.FormatInt(mid, 10) //每日任务
	} else {
		key = "creermission|" + strconv.FormatInt(mid, 10) //生涯任务
	}
	return key
}

/**
 * 获取任务完成情况(如果Redis中没有会从DB获取并设置缓存)
 * @param int64 $mid 用户ID
 * @param int32 $missionId 任务id
 * @return int32 任务状态
 */
func (this *missionRewardService) getStatus(missionId int64, mid int64, circle_type int32) (int64) {
	redisObj, err := NewRedisCache("missiondata")
	if err != nil {
		beego.Error("new userscore err:", err.Error())
		return 0
	}

	key := this.getKey(mid, circle_type)
	v, err := redis.String(redisObj.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}
	finished := strings.Split(v, ",")
	var flag int64
	flag = 0
	for _, taskId := range finished {
		n, err := strconv.Atoi(taskId)
		if err != nil {
			//fmt.Println("字符串转换成整数失败:", err)
		}
		if n == int(missionId) {
			flag = 2
		}
	}
	if flag == 0 {
		//没有完成，查找数据库
		o.Using("logs")
		var data entity.Missionreward
		err := o.QueryTable(this.table(circle_type)).Filter("mid", mid).Filter("mission", missionId).One(&data, "Status")
		//没有找到记录
		if err == orm.ErrNoRows {
			return 0
		} else {
			value := v + strconv.FormatInt(int64(data.Status), 10) + ","
			redis.String(redisObj.Do("SET", key, value))
			return int64(data.Status)
		}

	} else {
		return flag
	}
}
