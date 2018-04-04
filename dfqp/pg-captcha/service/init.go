package service

import (
	"framework/rpcclient/core"
	"putil/log"
	"github.com/astaxie/beego/config"
	"io/ioutil"
	"strings"
)

var (
	Runmode					  string //项目运行环境
	Client 					  *rpcclient.RpcCall
	ByClientService			  *byClientService
	BoyaaSmsService			  *boyaaSmsService
	ServiceConf				  map[string]config.Configer
)

//项目初始化
func Init()  {
	initCfg()
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

// 初始化服务对象
func initService()  {
	var err error
	Client, err = rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		panic("new rpccall fail")
	}

	ByClientService = &byClientService{}
	BoyaaSmsService = &boyaaSmsService{}
}