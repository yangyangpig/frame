package routers

import (
	"ServerStackMonitorSystem/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	//beego.AutoRouter(&controllers.MainController{})
	beego.Router("/", &controllers.TestToolController{}, "get:Index")
	beego.Router("/start", &controllers.TestToolController{}, "post:GetResult")
	beego.Router("/servrestatus", &controllers.ServerRegisterDataController{}, "get:ShowData")

}
