package serverstatus

import (
	"ServerStackMonitorSystem/tool"
	"encoding/json"
	"errors"
	_ "fmt"
	"log"
	"strconv"
)

var (
	zk  *tool.Zookeeper
	err error
)

type ZkServerRegStruct struct {
	Typ     int64   `json:"Typ"`     //服务类型
	Ver     int32   `json:"Ver"`     //服务版本号
	Ids     []int64 `json:"Ids"`     //服务最新的实例id片
	PrevVer int32   `json:"PrevVer"` //上个版本号
	PrevIds []int64 `json:"PrevIds"` //上个版本的实例id片 `json:"PrevIds,omitempty"`当空的时候不进json
}

func init() {
	zk, _ = tool.NewZookeeper()
	//连接zk
	zk, err = zk.ZookeeperConnect()
	if err != nil {
		panic("connect zookeeper is fail!")
	}
}

func GetServerRegisnInfo(path string) (lastData *ZkServerRegStruct, err error) {
	var server *ZkServerRegStruct
	if path == "" {
		return nil, errors.New("path is empty")
	}
	regisnData, _ := GetDataByPath(path)
	err = json.Unmarshal(regisnData, &server)

	lastData = server
	if err != nil {
		log.Println("json  unserialize fail")
		return
	}
	return

}

func GetReginsFuncInfo(path string) (collect []string, err error) {
	var funcInfo []string
	if path == "" {
		return nil, errors.New("path is empty")
	}
	regisnData, _ := GetDataByPath(path)
	err = json.Unmarshal(regisnData, &funcInfo)
	collect = funcInfo
	if err != nil {
		log.Println("json  unserialize fail")
		return
	}
	return
}

func GetDataByPath(path string) (res []byte, err error) {
	//首先判断节点是否存在
	flag := zk.ExistPath(path)
	if !flag {
		return nil, errors.New("node is not exist")
	}
	res, err = zk.GetData(path)
	if err != nil {
		log.Println("get data fail")
		return
	}
	return
}

func GetChildNodeByPath(path string) (res []string, err error) {
	flag := zk.ExistPath(path)
	if !flag {
		return nil, errors.New("node is not exist")
	}
	res, err = zk.GetChildren(path)
	if err != nil {
		log.Println("get data fail")
		return
	}
	return
}

func GetAllData() (map[string][]map[string]interface{}, error) {
	var laste = make(map[string][]map[string]interface{})
	var Showdata []map[string][]map[string]interface{}
	var err error
	var able = true

	childNode, err := GetChildNodeByPath("/rpc")
	if err != nil {
		return laste, err
	}
	for _, value := range childNode {
		var preable = true
		var idmap = make(map[string][]int64)
		var tmp2 []map[string]interface{}
		//获取到子节点所有数据

		childtypePath := "/rpc" + "/" + value + "/" + "typeandids"
		childFunPath := "/rpc" + "/" + value + "/" + "funcs"

		typeData, _ := GetServerRegisnInfo(childtypePath)
		//判断旧版本实例是否可用
		if len(typeData.Ids) > len(typeData.PrevIds) {
			//旧实例id不可用
			preable = false
		}

		idmap["previds"] = typeData.PrevIds
		idmap["ids"] = typeData.Ids

		for k, _ := range idmap {
			if k == "ids" {
				for _, v := range typeData.Ids {
					var agentAble = false
					var tmp = make(map[string]interface{})
					tmp["servname"] = value
					tmp["id"] = v                 //实例id
					tmp["version"] = typeData.Ver //当前版本
					tmp["type"] = typeData.Typ    //之前版本

					//获取该实例下该版本下所有方法
					funcPath := childFunPath + "/" + strconv.Itoa(int(typeData.Ver))
					funcCollect, _ := GetReginsFuncInfo(funcPath)
					tmp["funcCollect"] = funcCollect //服务下所有方法
					tmp["available"] = able          //当前版本一定可用

					//判断在代理是否可到达
					flag := strconv.Itoa(int(typeData.Typ)) + "_" + strconv.Itoa(int(v))
					agentPaht := "/agentNew/agent_list/192.168.202.25:7000/" + flag
					agentFlag := zk.ExistPath(agentPaht)
					if agentFlag {
						agentAble = true
					}
					tmp["agencyAvailable"] = agentAble
					tmp2 = append(tmp2, tmp) //把所有id全部放到一个[]里面

				}
				laste[value] = tmp2
			} else {
				for _, v := range typeData.PrevIds {
					var agentAble = false
					var tmp = make(map[string]interface{})

					tmp["servname"] = value
					tmp["id"] = v                     //实例id
					tmp["version"] = typeData.PrevVer //当前版本
					tmp["type"] = typeData.Typ        //之前版本

					//获取该实例下该版本下所有方法
					funcPath := childFunPath + "/" + strconv.Itoa(int(typeData.PrevVer))
					funcCollect, _ := GetReginsFuncInfo(funcPath)
					tmp["funcCollect"] = funcCollect //服务下所有方法
					tmp["available"] = preable       //之前版本是否可用

					//判断在代理是否可到达
					flag := strconv.Itoa(int(typeData.Typ)) + "_" + strconv.Itoa(int(v))
					agentPaht := "/agentNew/agent_list/192.168.202.25:7000/" + flag
					agentFlag := zk.ExistPath(agentPaht)
					if agentFlag {
						agentAble = true
					}
					tmp["agencyAvailable"] = agentAble
					tmp2 = append(tmp2, tmp) //把所有id全部放到一个[]里面

				}
				laste[value] = tmp2
			}
		}

		Showdata = append(Showdata, laste)
	}
	return laste, nil

}
