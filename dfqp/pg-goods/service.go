package main

import (
	"runtime"
	"dfqp/pg-goods/service"
	"github.com/astaxie/beego/config"
)

var serviceConfig config.Configer

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	println("--loading config...")
	configFile := "/Users/liyufeng/Documents/boyaa-new-hall/src/dfqp/pg-goods/conf/pg-goods.conf"
	var err error
	serviceConfig, err = config.NewConfig("ini", configFile)
	if err != nil {
		panic(err.Error())
	}
	println("--init logger ...")
	// logFile := serviceConfig.String("log::logFile")
	// rumMode := serviceConfig.String("rpc::rpc.run_mode")
	// plog.SetOutPutFileRotate(logFile, 1*1024*1024, 3, 30)
	// plog.LogRegister("debug")
}

var quitFlag chan bool

// PGGoods RPC Server
func main() {
	quitFlag = make(chan bool)
	println("--create pGService ...")
	pgGoodsService := service.NewPGGoods()
	defer pgGoodsService.Exit()
	pgGoodsService.SetConfig(serviceConfig)
	pgGoodsService.Init()
	println("--pGService ready go ...")
	pid, err := pgGoodsService.Run()
	if err != nil {
		panic(err)
	}
	println(pgGoodsService.Name(), "is running, pid=", pid)

	// 控制退出
	_ = <-quitFlag // 收消息退出！
}
