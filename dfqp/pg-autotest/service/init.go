package service

import (
	"github.com/astaxie/beego/config"
	"framework/rpcclient/core"
	"putil/log"
	"io/ioutil"
	"strings"
	"dfqp/lang/zh"
	"dfqp/lang/tw"
)

var (
	Runmode				string	//项目运行环境
	ServiceConf			map[string]config.Configer
	LangMap				map[string]map[int]string
	Client				*rpcclient.RpcCall
)

func Init() {
	//初始化配置
	initCfg()
	//注册语言包
	initLang()
	//初始化服务对象
	initService()
}

func initCfg() {
	confPath := "./conf"
	dirList, err := ioutil.ReadDir(confPath)
	if err != nil {
		panic("无效的配置文件, err:" + err.Error())
	}
	ServiceConf = make(map[string]config.Configer)
	for _, v := range dirList {
		nameArr := strings.Split(v.Name(), ".")
		nameKey := nameArr[0]
		tmpConf, _ := config.NewConfig("ini", confPath + "/" + v.Name())
		ServiceConf[nameKey] = tmpConf
	}
	Runmode = ServiceConf["app"].String("runmode")
}

func initLang() {
	LangMap = make(map[string]map[int]string)
	LangMap["zh"] = zh.LangMap
	LangMap["tw"] = tw.LangMap
}

func Lang(code int) string {
	lang := ServiceConf["app"].String("app.lang")
	if lang == "" {
		lang = "zh"
	}
	if v, ok := LangMap[lang][code]; ok {
		return v
	} else {
		return ""
	}
}

//初始化服务对象
func initService() {
	var err error
	Client, err = rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		panic("new rpccall fail")
	}

}