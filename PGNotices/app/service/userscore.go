package service

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"PGNotices/app/entity"
)

type userScoreService struct {}

const (
	//信任值字段类型
	TYPE_POSITIVE_SCORE = 1
	//Reids的Hash数量(按mid取模分布)
	REDIS_HASH_NUMBER = 10000
)

/**
 * 获取用户的信任值
 * @param int64 $mid 用户ID
 * @return int64 评分
 */
func (u *userScoreService) GetNegativeScore(mid int64) int64 {
	if mid <= 0 {
		return 0
	}
	return u.getScore(mid, TYPE_POSITIVE_SCORE)
}

/**
 * 获取用户的评分(如果Redis中没有会从DB获取并设置缓存)
 * @param int64 $mid 用户ID
 * @param int32 $ntype 字段类型大于0为信任值，其它为罪恶值
 * @return int64 评分
 */
func (u *userScoreService) getScore(mid int64, ntype int32) int64 {
	redisObj, err := NewRedisCache("userscore")
	if err != nil {
		beego.Error("new userscore err:", err.Error())
		return 0
	}
	key := u.getKey(mid, ntype)
	v, err := redis.Int64(redisObj.Do("HGET", key, mid))
	if err != nil && err != redis.ErrNil {
		beego.Error("userscore.Get err:", err.Error())
		return 0
	}
	//redis中查询不到数据
	if err == redis.ErrNil {
		res := u.load(mid)
		var val int32
		if ntype > 0 {
			val = res.PositiveScore
		} else {
			val = res.NegativeScore
		}
		_, err := redisObj.Do("HSET", key, mid, val)
		if err != nil {
			beego.Error("userscore.Get err:", err.Error())
			return 0
		}
	}
	return v
}

/**
 * 生成Redis缓存的Key
 * @param int64 $mid 用户ID
 * @param int32 $type 字段类型大于0为信任值，其它为罪恶值
 * @return string 缓存的KEY
 */
func (u *userScoreService) getKey(mid int64, ntype int32) string {
	newMid := mid % REDIS_HASH_NUMBER
	if ntype > 0 {
		return fmt.Sprintf("p%d", newMid)
	} else {
		return fmt.Sprintf("n%d", newMid)
	}
}

/**
 * 从DB获取用户的字段值
 * @param int64 $mid 用户ID
 * @return entity.Userscore 用户DB中的值
 */
func (u *userScoreService) load(mid int64) entity.Userscore {
	table := tableNameByMid("users.userscore", mid)
	var equipData entity.Userscore
	o.Using("users")
	sql := fmt.Sprintf("SELECT * FROM %s WHERE mid=%d LIMIT 1", table, mid)
	fmt.Println(sql)
	err := o.Raw(sql).QueryRow(&equipData)
	if err != nil {
		beego.Error("userscore.load err:", err.Error())
		return entity.Userscore{}
	}
	return equipData
}