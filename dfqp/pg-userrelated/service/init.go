package service

import (
	"framework/rpcclient/core"
	"putil/log"
)

var (
	Client					  *rpcclient.RpcCall
)

//项目初始化
func Init() {
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
