package controller

import (
	"strconv"
	"math"
	"reflect"
	"strings"
	"time"
	"dfqp/pg-userrelated/service"
	"dfqp/proto/user"
	"dfqp/proto/userrelated"
)

const (
	TIMEOUT = 2
)

type UserRelatedController struct {
}

//分发控制器
//@param struct rq 请求参数
func (this *UserRelatedController) Dispatch(rq *pgUserRelated.UserRelatedRequest) *pgUserRelated.UserRelatedResponse {
	res := &pgUserRelated.UserRelatedResponse{}
	s1 := int64(rq.S1)
	s1Bin := strconv.FormatInt(s1, 2)
	arr := map[string]string{
		"S1" : s1Bin,
	}
	methods := make([]string, 0)
	for key, bin := range arr {
		length := len(bin)
		for i:=0; i<length; i++ {
			if bin[i] == 49 {
				code := math.Pow(2, float64(length-i-1))
				s := strconv.FormatFloat(code, 'G', -1, 64)
				method := key+"_"+s
				v := reflect.ValueOf(this)
				mv := v.MethodByName(method)
				if mv.IsValid() {
					methods = append(methods, method)
				}
			}
		}
	}
	//超时处理
	chs := make([]chan bool, len(methods))
	for i, act := range methods {
		chs[i] = make(chan bool)
		go this.async(act, rq, res, chs[i])
	}
	for _, ch := range chs {
		<-ch
	}
	return res
}

func (this *UserRelatedController) async(method string, rq *pgUserRelated.UserRelatedRequest, res *pgUserRelated.UserRelatedResponse, ch chan bool) {
	runCh := make(chan bool)
	go func() {
		obj := reflect.ValueOf(this)
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

//用户个人信息
func (this *UserRelatedController) S1_1(rq *pgUserRelated.UserRelatedRequest) pgUserRelated.UserInfoResp {
	response := pgUserRelated.UserInfoResp{
		Status:1,
	}
	mid := rq.GetMid()
	if mid <= 0 {
		return response
	}
	request := new(pgUser.GetUserInfoRequest)
	request.Mid = mid
	reqBytes, err := request.Marshal()
	if err != nil {
		return response
	}
	resp := service.Client.SendAndRecvRespRpcMsg("pgUser.getUserInfo", reqBytes, 1000, 0)
	if resp.ReturnCode != 0 {
		return response
	}
	pgResp := new(pgUser.GetUserInfoResponse)
	err = pgResp.Unmarshal(resp.Body)
	if err != nil {
		return response
	}
	//成功
	if pgResp.Status == 0 {
		response.Status = 0
		response.Data.City = pgResp.Data.City
		response.Data.Nick = pgResp.Data.Nick
		response.Data.Icon = pgResp.Data.Icon
		response.Data.IconBig = pgResp.Data.IconBig
		response.Data.Sex = pgResp.Data.Sex
		response.Data.Sign = pgResp.Data.Sign
		response.Data.Sign = pgResp.Data.Sign
		response.Data.IconId = pgResp.Data.IconId
	}
	return response
}