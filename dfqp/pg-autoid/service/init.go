package service

import (
	"dfqp/pg-autoid/entity"
	//"fmt"
	"putil/log"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	o             orm.Ormer                 //mysql实例
	tablePrefix   = make(map[string]string) //表前缀map
	Runmode       string                    //项目运行环境
	AutoidService *autoidService
	//CuserService  *cuserService
	//CidMapService *cidMapService
)

//项目初始化
func Init() {
	//日志注册
	logs.SetLevel(logs.LevelDebug)
	//	logs.SetLogFuncCall(true)
	//	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/debug.log","daily":true}`)

	iniconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		panic("加载conf配置失败,err:" + err.Error())
	}
	Runmode = iniconf.String("runmode")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	registerDb("common")
	//非分表的才加这个方法
	orm.RegisterModelWithPrefix(tablePrefix["common"]+"_", new(entity.AutoidSource))

	orm.Debug = true
	o = orm.NewOrm()
	orm.RunCommand()

	// 初始化服务对象
	initService()
}

// 初始化服务对象
func initService() {
	AutoidService = &autoidService{}
	//	CuserService = &cuserService{}
	//	CidMapService = &cidMapService{}
}

//注册db服务
func registerDb(alias string) {
	dbConf, err := config.NewConfig("ini", "./conf/db.conf")
	if err != nil {
		panic("加载db配置失败,err:" + err.Error())
	}
	newAlias := alias
	//必须注册一个default
	if alias == "common" {
		newAlias = "default"
	}
	plog.Debug("xxxxxxxxxxxxxxx" + Runmode + "::db." + alias + ".mysqluser")
	user := dbConf.String(Runmode + "::db." + alias + ".mysqluser")
	pass := dbConf.String(Runmode + "::db." + alias + ".mysqlpass")
	urls := dbConf.String(Runmode + "::db." + alias + ".mysqlurls")
	dbName := dbConf.String(Runmode + "::db." + alias + ".mysqldb")
	prefix := dbConf.String(Runmode + "::db." + alias + ".mysqlprefix")
	plog.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" + user + ":" + pass + "@tcp(" + urls + ")/" + dbName + "?charset=utf8")
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
		return prefix + "_" + arr[1]
	}
	return ""
}
