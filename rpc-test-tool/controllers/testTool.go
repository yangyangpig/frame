package controllers

import (
	"encoding/json"
	"fmt"
	"rpc-test-tool/models"

	"github.com/astaxie/beego"
)

type TestToolController struct {
	beego.Controller
}

type ReqStruct struct {
	ServeName string
	FuncName  string
	Args      []ParmStruct
}
type ParmStruct struct {
	Name  string
	Value string
}

func (c *TestToolController) Index() {

	c.TplName = "index.tpl"
}

func (c *TestToolController) GetResult() {
	var ob ReqStruct
	var servname string
	var funcname string
	var para []ParmStruct

	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	servname = ob.ServeName
	funcname = ob.FuncName
	para = ob.Args

	t := assembleData(servname, funcname, para)

	res := models.Process(t) //返回的对应的pb转化过后的字节切片
	lastres, err := models.SendDataCli(res, t.Servername, t.Funcname)
	if err != nil {
		lastres = "获取远端服务失败"
	}

	c.Data["json"] = lastres
	c.ServeJSON()

}

//组装数据
func assembleData(servName string, fucName string, param []ParmStruct) (res *models.ProcessData) {
	var container = make(map[string]string)
	res = new(models.ProcessData)
	//模拟表单数据
	res.Servername = servName
	res.Funcname = fucName

	for _, v := range param {
		//v是一个结构体
		container[v.Name] = v.Value
	}

	res.Param = container
	return
}
