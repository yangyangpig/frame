package main

import (
	_ "ServerStackMonitorSystem/controllers"
	_ "ServerStackMonitorSystem/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
	//	obj := new(controllers.ServerRegisterDataController)
	//	obj.ShowData()

}
