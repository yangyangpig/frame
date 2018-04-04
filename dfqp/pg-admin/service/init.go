package service

import (
	"io/ioutil"
	"github.com/astaxie/beego/config"
	"strings"
	"framework/rpcclient/core"
	"putil/log"
)

var (
	Runmode					  string //项目运行环境
	ServiceConf				  map[string]config.Configer
	Client					  *rpcclient.RpcCall
)

func Init() {
	//初始化配置
	initCfg()
	// 初始化服务对象
	initService()
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
		tmpConf, _ := config.NewConfig("ini", confPath+"/"+v.Name())
		ServiceConf[nameKey] = tmpConf
	}
	Runmode = ServiceConf["app"].String("runmode")
}