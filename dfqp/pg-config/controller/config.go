package controller

import (
	"dfqp/proto/config"
	"strconv"
	"math"
	"reflect"
	"strings"
	"time"
	"dfqp/pg-config/service"
)

const (
	TIMEOUT = 2
)

type ConfigController struct {}

//分发控制器
//@param struct rq 请求参数
func (n *ConfigController) Dispatch(rq *pgConfig.ConfigRequest) *pgConfig.ConfigResponse {
	res := &pgConfig.ConfigResponse{}
	if rq.Appid < 0 || rq.Region < 0 {
		return res
	}
	s1 := int64(rq.S1)
	s1Bin := strconv.FormatInt(s1, 2)
	s2 := int64(rq.S2)
	s2Bin := strconv.FormatInt(s2, 2)
	arr := map[string]string{
		"S1" : s1Bin,
		"S2" : s2Bin,
	}
	methods := make([]string, 0)
	for key, bin := range arr {
		length := len(bin)
		for i:=0; i<length; i++ {
			if bin[i] == 49 {
				code := math.Pow(2, float64(length-i-1))
				s := strconv.FormatFloat(code, 'G', -1, 64)
				method := key+"_"+s
				v := reflect.ValueOf(n)
				mv := v.MethodByName(method)
				if mv.IsValid() {
					methods = append(methods, method)
				}
			}
		}
	}
	//超时处理思路1
	chs := make([]chan bool, len(methods))
	for i, act := range methods {
		chs[i] = make(chan bool)
		go n.async(act, rq, res, chs[i])
	}
	for _, ch := range chs {
		<-ch
	}
	return res
}
//反射处理
func (n *ConfigController) async(method string, rq *pgConfig.ConfigRequest, res *pgConfig.ConfigResponse, ch chan bool) {
	runCh := make(chan bool)
	go func() {
		obj := reflect.ValueOf(n)
		mv := obj.MethodByName(method)
		if mv.IsValid() {
			args := []reflect.Value{reflect.ValueOf(rq)}
			result := mv.Call(args)
			slices := strings.Split(method, "_")
			if len(slices) == 2 {
				rel := reflect.ValueOf(res).Elem()
				mark := slices[0]
				name := rel.FieldByName(mark)
				fl := rel.FieldByName("Flag")
				if name.IsValid() && fl.IsValid() {
					cmd := name.FieldByName("Cmd"+slices[1])
					num := string([]byte(mark)[1:])
					s := fl.FieldByName("S"+num)
					if cmd.IsValid() && s.IsValid() && cmd.CanSet() && s.CanSet() {
						s.Set(reflect.Append(s, reflect.ValueOf("cmd"+slices[1])))
						cmd.Set(result[0])
						runCh <- false
					}
				}
			}
		}
		return
	}()
	select {
	case <- time.After(TIMEOUT*time.Second):
		ch <- true
	case <- runCh:
		ch <- false
	}
}
//access配置
func (n *ConfigController) S1_1(rq *pgConfig.ConfigRequest) pgConfig.AccessDomainResp {
	resp := pgConfig.AccessDomainResp{}
	ret := service.SiteDomainService.Get()
	if len(ret) > 0 {
		resp.Status = 1
		resp.Data = ret
	}
	return resp
}

//access配置
func (n *ConfigController) S1_2(rq *pgConfig.ConfigRequest) pgConfig.AccessDomainResp {
	resp := pgConfig.AccessDomainResp{}
	ret := service.SiteDomainService.Get()
	if len(ret) > 0 {
		resp.Status = 1
		resp.Data = ret
	}
	return resp
}