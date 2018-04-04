package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"PGMatch/app/entity"
	"PGMatch/app/proto"
	"strings"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	//o                  orm.Ormer         // mysql实例
	//tablePrefix        map[string]string //
	FreeTimeService    entity.FreeTimer     // 每日免费次数 redis
	MatchConfService   entity.MatchConfer   // 比赛配置 redis
	SessionService     entity.Sessioner     // session redis
	GenericConfService entity.GenericConfer // 通用配置服务
	MatchInviteService entity.MatchInviter  // 比赛匹配服务
	Response           *proto.ListsResponse // 返回的结构体
)
// redis配置信息
var RedisConf = map[string]map[string]string{}

// 初始函数
func Init() {
	moduleName := "PGMatch"
	loadLogConf(moduleName)
	// 数据库注册
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//registerDb("")
	////orm.RegisterModelWithPrefix(tablePrefix[""]+"_", new(entity.Match))
	//if beego.BConfig.RunMode == "dev" {
	//	orm.Debug = true
	//}
	//o = orm.NewOrm()
	//orm.RunCommand()
	loadRedisConf(moduleName)

	// 初始化服务对象
	initService()
}

// 初始化服务
func initService() {
	MatchConfService = NewMatchConfService()
	FreeTimeService = NewFreeTimeService()
	SessionService = NewSessionService()
	MatchInviteService = NewMatchInviteService()
	GenericConfService = NewGenericConfService()
}

// 加载日志配置
func loadLogConf(moduleName string)  {
	path, _ := os.Getwd()
	pathSlice := strings.Split(path, moduleName)
	num := strings.Split(pathSlice[1], "\\")
	str := ""
	for i:=0; i < len(num) -1; i++ {
		str += "../"
	}
	confPath := str + "logs/" + moduleName + ".log"
	//confPath := pathSlice[0] + moduleName + "\\logs\\" + moduleName + ".log"
	// 日志注册
	l := logs.GetLogger()
	l.Println("this is a message of matchservice")
	// 设置日志写入缓冲区的等级
	logs.SetLevel(beego.LevelDebug)
	// 输出log时能显示输出文件名和行号（非必须）
	logs.SetLogFuncCall(true)
	// 	设置日志记录方式：本地文件记录
	// 	filename 保存的文件名
	//	maxlines 每个文件保存的最大行数，默认值 1000000
	//	maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	//	daily 是否按照每天 logrotate，默认是 true
	//	maxdays 文件最多保存多少天，默认保存 7 天
	//	rotate 是否开启 logrotate，默认是 true
	//	level 日志保存的时候的级别，默认是 Trace 级别
	//	perm 日志文件权限
	jsonConf := `{"filename":"` + confPath + `","daily":true}`
	logs.SetLogger(logs.AdapterFile, jsonConf)
	// 设置单独日志文件
	//	ilename 保存的文件名
	//	maxlines 每个文件保存的最大行数，默认值 1000000
	//	maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	//	daily 是否按照每天 logrotate，默认是 true
	//	maxdays 文件最多保存多少天，默认保存 7 天
	//	rotate 是否开启 logrotate，默认是 true
	//	level 日志保存的时候的级别，默认是 Trace 级别
	//	perm 日志文件权限
	//	separate 需要单独写入文件的日志级别,设置后命名类似 test.error.log
	//errJsonConf := `{"filename":"` + confPath+ `","daily":true,"separate":["error"]}`
	//logs.SetLogger(logs.AdapterMultiFile, errJsonConf)
	// 异步输出
	//logs.Async()
}

// 加载 redis yml 配置文件
func loadRedisConf(moduleName string)  {
	confFileName := "redis.dev.yml"
	path, _ := os.Getwd()
	pathSlice := strings.Split(path, moduleName)
	confPath := pathSlice[0] + moduleName + "\\conf\\" + confFileName
	ymlData, _ := ioutil.ReadFile(confPath)
	err := yaml.Unmarshal(ymlData, &RedisConf)
	if err != nil {
		logs.Error("Load redis genericconf failed! module: %v file: %v", moduleName, confFileName)
		panic(err)
	}
}

// 注册 DataBase
func registerDb(alias string) {

	//newAlias := alias
	//if alias == "main" {
	//	newAlias = "default"
	//}
	//user := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + ".mySqlUser")
	//pass := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + ".mySqlPass")
	//urls := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + ".mySqlUrls")
	//dbName := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + ".mySqlDb")
	//prefix := beego.AppConfig.String(beego.BConfig.RunMode + "::" + alias + ".mySqlPrefix")
	//orm.RegisterDataBase(newAlias, "mysql", user+":"+pass+"@tcp("+urls+")/"+dbName+"?charset=utf8", 30)
	//tablePrefix = make(map[string]string)
	//tablePrefix[alias] = prefix
}

func tableName(name string) string {
	//arr := strings.Split(nameKey, ".")
	//if len(arr) != 2 {
	//	return ""
	//}
	//if prefix, ok := tablePrefix[arr[0]]; ok {
	//	return prefix + "_" + arr[1]
	//}
	return ""
}
