package service

import (
	"dfqp/pg-mission/entity"
	"framework/rpcclient/core"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"putil/log"
	"strings"
	"time"
)

var (
	o                   orm.Ormer                 //mysql实例
	tablePrefix         = make(map[string]string) //表前缀map
	runmode             string                    //项目运行环境
	Runmode             string                    //项目运行环境
	MissionTypesService *missionTypesService
	Client              *rpcclient.RpcCall
)

//项目初始化
func Init() {
	//日志注册
	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCall(true)
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/debug.log","daily":true}`)

	iniconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		panic("加载conf配置失败,err:" + err.Error())
	}
	runmode = iniconf.String("runmode")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	registerDb("main")
	registerDb("logs")
	orm.RegisterModelWithPrefix(tablePrefix["main"]+"_", new(entity.Missiontypes))
	o = orm.NewOrm()
	orm.RunCommand()

	// 初始化服务对象
	initService()
}

// 初始化服务对象
func initService() {
	var err error
	Client, err = rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		panic("new rpccall fail")
	}
	MissionTypesService = &missionTypesService{}
}

//注册db服务
func registerDb(alias string) {
	dbConf, err := config.NewConfig("ini", "./conf/db.conf")
	if err != nil {
		panic("加载db配置失败,err:" + err.Error())
	}
	newAlias := alias
	//必须注册一个default
	if alias == "main" {
		newAlias = "default"
	}
	user := dbConf.String(runmode + "::db." + alias + ".mysqluser")
	pass := dbConf.String(runmode + "::db." + alias + ".mysqlpass")
	urls := dbConf.String(runmode + "::db." + alias + ".mysqlurls")
	dbName := dbConf.String(runmode + "::db." + alias + ".mysqldb")
	prefix := dbConf.String(runmode + "::db." + alias + ".mysqlprefix")
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

//按日期分表
func tableNameByDate(name string) string {
	arr := strings.Split(name, ".")
	if len(arr) != 2 {
		return ""
	}

	dateNow := time.Now().Format("20060102")
	if prefix, ok := tablePrefix[arr[0]]; ok {
		return prefix + "_" + arr[1] + dateNow
	}
	return ""
}
