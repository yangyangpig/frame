package service

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/config"
	"strings"
	"fmt"
)

var (
	o 				          orm.Ormer   			//mysql实例
	tablePrefix 	          = make(map[string]string)  //表前缀map
	Runmode					  string //项目运行环境
	UserService			      *userService
)

//项目初始化
func Init()  {
	iniconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		panic("加载conf配置失败,err:"+err.Error())
	}
	Runmode = iniconf.String("runmode")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	registerDb("users")

	o = orm.NewOrm()
	orm.RunCommand()

	// 初始化服务对象
	initService()
}

// 初始化服务对象
func initService()  {
	UserService = &userService{}
}

//注册db服务
func registerDb(alias string)  {
	dbConf, err := config.NewConfig("ini", "./conf/db.conf")
	if err != nil {
		panic("加载db配置失败,err:"+err.Error())
	}
	newAlias := alias
	//必须注册一个default
	if alias == "users" {
		newAlias = "default"
	}
	user := dbConf.String(Runmode+"::db."+alias+".mysqluser")
	pass := dbConf.String(Runmode+"::db."+alias+".mysqlpass")
	urls := dbConf.String(Runmode+"::db."+alias+".mysqlurls")
	dbName := dbConf.String(Runmode+"::db."+alias+".mysqldb")
	prefix := dbConf.String(Runmode+"::db."+alias+".mysqlprefix")
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

//根据mid模100分表
func tableNameByMid(name string, mid int64) string {
	arr := strings.Split(name, ".")
	if len(arr) != 2 {
		return ""
	}
	if prefix, ok := tablePrefix[arr[0]]; ok {
		m := mid % 100
		return fmt.Sprintf("%s_%s%d", prefix, arr[1], m)
	}
	return ""
}