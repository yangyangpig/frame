package controllers

import (
	"ServerStackMonitorSystem/models/serverstatus"
	_ "fmt"

	"github.com/astaxie/beego"
)

type ServerRegisterDataController struct {
	beego.Controller
}

func (c *ServerRegisterDataController) ShowData() {
	data, _ := serverstatus.GetAllData()
	c.Data["json"] = data
	c.ServeJSON()
}
