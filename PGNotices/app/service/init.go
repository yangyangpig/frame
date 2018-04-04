package service

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
	"PGNotices/app/entity"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	//"framework/rpcclient/szmq"
	//"framework/rpcclient/szmq"
	"fmt"
)

const (
	MOD = 3
)

var (
	o 				          orm.Ormer   			//mysql实例
	tablePrefix 	          = make(map[string]string)
	NoticesService            *noticesService        // 公告服务
	SystemCacheService        *systemCacheService    // 公告缓存服务
	RegionService		      *regionService         // 地区配置服务
	UserScoreService		  *userScoreService
	//Logger 					  *szmq.SzmqPushClient
)

func Init()  {
	//日志注册
	beego.SetLevel(beego.LevelDebug)
	beego.SetLogFuncCall(true)
	beego.SetLogger(logs.AdapterFile, `{"filename":"logs/debug.log","daily":true}`)

	orm.RegisterDriver("mysql", orm.DRMySQL)

	registerDb("main")
	orm.RegisterModelWithPrefix(tablePrefix["main"] + "_", new(entity.Notices))
	registerDb("users")

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	o = orm.NewOrm()
	orm.RunCommand()

	// 初始化服务对象
	initService()
}

func initService() {
	NoticesService = &noticesService{}
	SystemCacheService = &systemCacheService{}
	RegionService = &regionService{}
	UserScoreService = &userScoreService{}
}

func registerDb(alias string)  {
	newAlias := alias
	if alias == "main" {
		newAlias = "default"
	}
	user := beego.AppConfig.String(beego.BConfig.RunMode+"::"+alias+".mysqluser")
	pass := beego.AppConfig.String(beego.BConfig.RunMode+"::"+alias+".mysqlpass")
	urls := beego.AppConfig.String(beego.BConfig.RunMode+"::"+alias+".mysqlurls")
	dbName := beego.AppConfig.String(beego.BConfig.RunMode+"::"+alias+".mysqldb")
	prefix := beego.AppConfig.String(beego.BConfig.RunMode+"::"+alias+".mysqlprefix")
	orm.RegisterDataBase(newAlias, "mysql", user+":"+pass+"@tcp("+urls+")/"+dbName+"?charset=utf8", 30)
	tablePrefix[alias] = prefix
}

//不分表
func tableName(name string) string {
	arr := strings.Split(name, ".")
	if len(arr) != 2 {
		return ""
	}
	if prefix, ok := tablePrefix[arr[0]]; ok {
		return prefix +"_"+ arr[1]
	}
	return ""
}

//根据mid分表
func tableNameByMid(name string, mid int64) string {
	arr := strings.Split(name, ".")
	if len(arr) != 2 {
		return ""
	}
	if prefix, ok := tablePrefix[arr[0]]; ok {
		m := mid % 10
		return fmt.Sprintf("%s_%s%d", prefix, arr[1], m)
	}
	return ""
}

