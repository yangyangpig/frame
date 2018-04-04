package service

import (
	"github.com/astaxie/beego/config"
	//"github.com/astaxie/beego/logs"
	//"putil/log"
	"framework/rpcclient/core"
	"io/ioutil"
	"strings"
)

var (
	Runmode      string //项目运行环境
	GoodsService *goodsService
	RedisService *redisService
	ServiceConf  map[string]config.Configer
	Client       *rpcclient.RpcCall
)

//项目初始化
func init() {
	initCfg()
	//日志注册
	//initLog()
	//注册服务对象
	initService()
}

func initService() {
	GoodsService = &goodsService{}
}

// 初始化所有配置
func initCfg() {
	confPath := "./conf"
	dirList, err := ioutil.ReadDir(confPath)
	if err != nil {
		panic("无有效配置目录, err:" + err.Error())
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

//func initLog()  {
//	if Runmode == "dev" {
//		logs.SetLevel(logs.LevelDebug)
//	}
//}
