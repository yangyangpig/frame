package service

import (
	"strings"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	o					orm.Ormer 					// Mysql实例
	tablePrefix			= make(map[string]string)	// 表前缀map
	runmode				string 						// 项目运行环境
	ConfigService		*configService
	SigninService		*signinService
	tableMap			= []string{"main"}			// 数据库 map
)

// 项目初始化
func Init() {
	initCfg()
	//日志注册
	initLog()
	// 初始化db
	initDb()
	// 初始化服务对象
	initService()
}

func initCfg()  {
	iniconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		panic("加载conf配置失败,err:"+err.Error())
	}
	runmode = iniconf.String("runmode")
}

func initLog()  {
	if runmode == "dev" {
		logs.SetLevel(logs.LevelDebug)
		logs.SetLogFuncCall(true)
		logs.SetLogger(logs.AdapterFile, `{"filename":"log/debug.log","daily":true}`)
	}
}

func initDb()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	for i, v := range tableMap {
		registerDb(i, v)
	}

	if runmode == "dev" {
		orm.Debug = true
	}
	o = orm.NewOrm()
	orm.RunCommand()
}

func initService() {
	ConfigService = &configService{}
	SigninService = &signinService{}
}

// 注册db服务
func registerDb(index int, alias string) {
	dbConf, err := config.NewConfig("ini", "./conf/db.conf")
	if err != nil {
		panic("加载db配置失败,err:" + err.Error())
	}
	newAlias := alias
	//必须注册一个default
	if index == 0 {
		newAlias = "default"
	}
	user := dbConf.String(runmode+"::db."+alias+".mysqluser")
	pass := dbConf.String(runmode+"::db."+alias+".mysqlpass")
	urls := dbConf.String(runmode+"::db."+alias+".mysqlurls")
	dbName := dbConf.String(runmode+"::db."+alias+".mysqldb")
	prefix := dbConf.String(runmode+"::db."+alias+".mysqlprefix")
	orm.RegisterDataBase(newAlias, "mysql", user+":"+pass+"@tcp("+urls+")/"+dbName+"?charset=utf8", 30)
	tablePrefix[alias] = prefix
}

// 不分表
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
