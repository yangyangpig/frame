package main

import (
	"encoding/json"
	_ "framework/rpcclient/szk"
	"putil/log"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"
)

type zkSt struct {
	value int `json:"Typ"` //服务类型
}

func main() {

	var zkIpPorts []string
	zkConfig, err := config.NewConfig("ini", "./conf/zk.conf")
	zkIpPorts = zkConfig.Strings("zk.zkipport")

	conn, _, cerr := zk.Connect(zkIpPorts, time.Second)
	plog.Debug(cerr)
	if cerr != nil {
		plog.Debug(err)
	}
	setPath(conn)
	//	for i := 0; i < 10; i++ {
	//		setPath()
	//	}
	//	for i := 0; i < 10; i++ {
	//		go setPathv2()
	//	}

}

func setPath(conn *zk.Conn) error {
	var (
		path        = "/zk"
		flags int32 = 0
		acls        = zk.WorldACL(zk.PermAll)
	)
	var typandidsObj zkSt = zkSt{
		value: 0,
	}

	plog.Debug("connect success")
	funcNodeExist, _, _ := conn.Exists(path)
	plog.Debug(funcNodeExist)
	if !funcNodeExist {
		value, _ := json.Marshal(typandidsObj)
		plog.Debug("value is ", value)
		_, err := conn.Create(path, value, flags, acls)
		plog.Debug("create err is ", err)

	} else {
		typandidsBytes, stat, _ := conn.Get(path)
		json.Unmarshal(typandidsBytes, &typandidsObj)

		plog.Debug("typandidsBytes is ", typandidsBytes)
		plog.Debug("typandidsBytes value is ", typandidsObj.value)

		updateValue(&typandidsObj, 1)

		plog.Debug("typandidsBytes is ", typandidsBytes)

		value, _ := json.Marshal(typandidsObj)

		plog.Debug("send zk path value is ", value)
		plog.Debug("stat.Version is ", stat.Version)

		_, errs := conn.Set(path, value, stat.Version)

		plog.Debug("set value return is ", errs)

	}
	return nil

}

func updateValue(typanvalue *zkSt, tmp int) {
	typanvalue.value = typanvalue.value + tmp
	plog.Debug("typandidsBytes is ", typanvalue.value)
}
