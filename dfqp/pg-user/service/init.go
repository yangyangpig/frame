package service

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/config"
	"strings"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"framework/rpcclient/core"
	"putil/log"
)

var (
	o 				          orm.Ormer   			//mysql实例
	tablePrefix 	          = make(map[string]string)  //表前缀map
	Runmode					  string //项目运行环境
	ServiceConf				  map[string]config.Configer
	Client					  *rpcclient.RpcCall
	tableMap  = []string{"users"}
)

//项目初始化
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
	confPath := "./conf"
	dirList, err := ioutil.ReadDir(confPath)
	if err != nil {
		panic("无有效配置文件, err:"+err.Error())
	}
	ServiceConf = make(map[string]config.Configer)
	for _, v := range dirList {
		nameArr := strings.Split(v.Name(), ".")
		nameKey := nameArr[0]
		// rpc框架配置直接过滤
		if nameKey == "zk" || nameKey == "zmq" {
			continue
		}
		tmpConf, _ := config.NewConfig("ini", confPath+"/"+v.Name())
		ServiceConf[nameKey] = tmpConf
	}
	Runmode = ServiceConf["app"].String("runmode")
}

func initLog()  {
	if Runmode == "dev" {
		logs.SetLevel(logs.LevelDebug)
	}
}

func initDb()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	for i, v :=range tableMap {
		registerDb(i, v)
	}

	if Runmode == "dev" {
		orm.Debug = true
	}
	o = orm.NewOrm()
	orm.RunCommand()
}

// 初始化服务对象
func initService()  {
	var err error
	Client, err = rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		panic("new rpccall fail")
	}
}

//注册db服务
func registerDb(index int, alias string)  {
	newAlias := alias
	//tableMap第一个注册为default
	if index == 0 {
		newAlias = "default"
	}
	user := ServiceConf["db"].String(Runmode+"::mysql."+alias+".user")
	pass := ServiceConf["db"].String(Runmode+"::mysql."+alias+".pwd")
	urls := ServiceConf["db"].String(Runmode+"::mysql."+alias+".host")
	dbName := ServiceConf["db"].String(Runmode+"::mysql."+alias+".db")
	prefix := ServiceConf["db"].String(Runmode+"::mysql."+alias+".prefix")
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

//根据cid模100分表
func tableNameBy(name string, cid int64, mod int) string {
	arr := strings.Split(name, ".")
	if len(arr) != 2 || mod < 2 {
		return ""
	}
	if prefix, ok := tablePrefix[arr[0]]; ok {
		m := cid % int64(mod)
		return fmt.Sprintf("%s_%s%d", prefix, arr[1], m)
	}
	return ""
}