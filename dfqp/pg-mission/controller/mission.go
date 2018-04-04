package controllers

import (
	"dfqp/pg-mission/service"
	"dfqp/proto/mission"
	"fmt"
	//"math"
	//"reflect"
	//"strconv"
	//"strings"
	//"sync"
	//"time"
)

type MissionController struct{}

//获取任务列表
//@param struct rq 请求参数
func (n *MissionController) MissionList(rq pgMission.MissionListRequest) pgMission.MissionListResponse {
	fmt.Println("############")
	res := pgMission.MissionListResponse{}
	fmt.Println("VIcoterew")
	if rq.App < 0 || rq.Mid < 0 || rq.AreaId < 0 {
		return res
	}

	//获取每日任务列表
	_, err := service.MissionTypesService.GetList(rq.AreaId, rq.Mid, 0)
	// s1 := int64(rq.S1)
	// s1Bin := strconv.FormatInt(s1, 2)
	// s2 := int64(rq.S2)
	// s2Bin := strconv.FormatInt(s2, 2)
	// arr := map[string]string{
	// 	"S1": s1Bin,
	// 	"S2": s2Bin,
	// }
	// methods := make([]string, 0)
	// for key, bin := range arr {
	// 	length := len(bin)
	// 	for i := 0; i < length; i++ {
	// 		if bin[i] == 49 {
	// 			code := math.Pow(2, float64(length-i-1))
	// 			s := strconv.FormatFloat(code, 'G', -1, 64)
	// 			method := key + "_" + s
	// 			v := reflect.ValueOf(n)
	// 			mv := v.MethodByName(method)
	// 			if mv.IsValid() {
	// 				methods = append(methods, method)
	// 			}
	// 		}
	// 	}
	// }
	// res = n.goAsync(methods, rq)
	fmt.Println(err)
	// res.Status = 0
	// res.Msg = "请求成功！"
	// res.Data = "这是返回的数据json字符串"
	return res
}

// //放在不同协程处理
// func (n *ConfigController) goAsync(methods []string, rq pgConfig.ConfigRequest) pgConfig.ConfigResponse {
// 	res := pgConfig.ConfigResponse{}
// 	count := len(methods)
// 	if count > 0 {
// 		var wg sync.WaitGroup
// 		wg.Add(count)
// 		for _, act := range methods {
// 			timer := time.NewTimer(2 * time.Second)
// 			ch := make(chan bool)
// 			go func(tm *time.Timer, cha chan bool) {
// 				for {
// 					select {
// 					case <-tm.C:
// 						cha <- true
// 						return
// 					}
// 				}
// 			}(timer, ch)
// 			go func(method string, tm *time.Timer, cha chan bool) {
// 				obj := reflect.ValueOf(n)
// 				mv := obj.MethodByName(method)
// 				if mv.IsValid() {
// 					args := []reflect.Value{reflect.ValueOf(rq)}
// 					result := mv.Call(args)
// 					slices := strings.Split(method, "_")
// 					if len(slices) == 2 {
// 						rel := reflect.ValueOf(&res).Elem()
// 						mark := slices[0]
// 						name := rel.FieldByName(mark)
// 						fl := rel.FieldByName("Flag")
// 						if name.IsValid() && fl.IsValid() {
// 							cmd := name.FieldByName("Cmd" + slices[1])
// 							num := string([]byte(mark)[1:])
// 							s := fl.FieldByName("S" + num)
// 							if cmd.IsValid() && s.IsValid() && cmd.CanSet() && s.CanSet() {
// 								s.Set(reflect.Append(s, reflect.ValueOf("cmd"+slices[1])))
// 								cmd.Set(result[0])
// 							}
// 						}
// 					}
// 					if tm.Stop() {
// 						cha <- true
// 					}
// 				}
// 			}(act, timer, ch)
// 			go func(cha chan bool) {
// 				<-cha
// 				wg.Done()
// 			}(ch)
// 		}
// 		wg.Wait()
// 	}
// 	return res
// }

// //access配置
// func (n *ConfigController) S1_1(rq pgConfig.ConfigRequest) pgConfig.AccessDomainResp {
// 	resp := pgConfig.AccessDomainResp{}
// 	ret := service.SiteDomainService.Get()
// 	if len(ret) > 0 {
// 		resp.Status = 1
// 		resp.Data = ret
// 	}
// 	return resp
// }
