package models

import (
	"ServerStackMonitorSystem/static/proto/autoidpro"
	_ "ServerStackMonitorSystem/static/proto/bigPackData"
	"ServerStackMonitorSystem/static/proto/captcha"
	"ServerStackMonitorSystem/static/proto/config"
	"ServerStackMonitorSystem/static/proto/goods"
	"ServerStackMonitorSystem/static/proto/login"
	"ServerStackMonitorSystem/static/proto/mission"
	_ "ServerStackMonitorSystem/static/proto/money"
	"ServerStackMonitorSystem/static/proto/notice"
	_ "ServerStackMonitorSystem/static/proto/online"
	"ServerStackMonitorSystem/static/proto/pay"
	"ServerStackMonitorSystem/static/proto/sign"
	_ "ServerStackMonitorSystem/static/proto/sms"
	"ServerStackMonitorSystem/static/proto/user"
	"ServerStackMonitorSystem/static/proto/userrelated"
	"encoding/json"
	"fmt"
	"reflect"

	"strconv"
	"strings"
)

type ProcessPb struct{}

func (this *ProcessPb) PgSignSignin(s map[string]string) []byte {

	pb := new(pgSign.SigninRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes
}

func (this *ProcessPb) PgSignSigninRepose(req []byte) (res string, err error) {

	pb := new(pgSign.SigninResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-login

func (this *ProcessPb) PgLoginLogin(s map[string]string) []byte {

	pb := new(pgLogin.LoginRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgLoginLoginRepose(req []byte) (res string, err error) {
	pb := new(pgLogin.LoginResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-autoid

func (this *ProcessPb) AutoidmanagerGetId(s map[string]string) []byte {

	pb := new(autoidpro.AutoidRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) AutoidmanagerGetIdRepose(req []byte) (res string, err error) {
	pb := new(autoidpro.AutoidResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-captcha

func (this *ProcessPb) PgCaptchaGetCaptcha(s map[string]string) []byte {

	pb := new(pgCaptcha.GetCaptchaRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgCaptchaGetCaptchaRepose(req []byte) (res string, err error) {
	pb := new(pgCaptcha.GetCaptchaResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

func (this *ProcessPb) PgCaptchaGetVoiceCaptcha(s map[string]string) []byte {

	pb := new(pgCaptcha.GetVoiceCaptchaRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgCaptchaGetVoiceCaptchaRepose(req []byte) (res string, err error) {
	pb := new(pgCaptcha.GetVoiceCaptchaResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-config
func (this *ProcessPb) PgConfigGet(s map[string]string) []byte {

	pb := new(pgConfig.ConfigRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgConfigGetRepose(req []byte) (res string, err error) {
	pb := new(pgConfig.ConfigResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-goods
func (this *ProcessPb) PgGoodsUse(s map[string]string) []byte {

	pb := new(ptGoods.ExchangeRealGoodsRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgGoodsUseRepose(req []byte) (res string, err error) {
	pb := new(ptGoods.GoodsUseResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-mission
func (this *ProcessPb) PGMissionMissionList(s map[string]string) []byte {

	pb := new(pgMission.MissionListRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PGMissionMissionListRepose(req []byte) (res string, err error) {
	pb := new(pgMission.MissionListResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-notice
func (this *ProcessPb) PgNoticeNoticeList(s map[string]string) []byte {

	pb := new(pgNotice.GetListRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgNoticeNoticeListRepose(req []byte) (res string, err error) {
	pb := new(pgNotice.GetListResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-pay
func (this *ProcessPb) PgPaySendMoney(s map[string]string) []byte {

	pb := new(pgPay.SendMoneyRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgPaySendMoneyRepose(req []byte) (res string, err error) {
	pb := new(pgPay.OrderResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-user

func (this *ProcessPb) PgUserGetUserInfo(s map[string]string) []byte {

	pb := new(pgUser.GetUserInfoRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgUserGetUserInfoRepose(req []byte) (res string, err error) {
	pb := new(pgUser.GetUserInfoRequest)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//pg-userrelated
func (this *ProcessPb) PgUserRelatedGet(s map[string]string) []byte {

	pb := new(pgUserRelated.UserRelatedRequest)
	TypeTransverter(pb, s)

	req_bytes, err := pb.Marshal()
	if err != nil {
		panic("pb序列化错误")
	}
	fmt.Println("pb后的数据为", req_bytes)
	return req_bytes

}

func (this *ProcessPb) PgUserRelatedGetRepose(req []byte) (res string, err error) {
	pb := new(pgUserRelated.UserRelatedResponse)
	pb.Unmarshal(req) //转换成json
	pbJson, err := json.Marshal(pb)
	if err != nil {
		fmt.Println("json转化错误")
		return
	}
	res = string(pbJson)
	return
}

//注意这里需要参数名称和转化的名称要相同才不会出错
func TypeTransverter(obj interface{}, s map[string]string) {
	t := reflect.TypeOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		//获取所有结构体的元素
		fileInfo := t.Field(i)
		kind := fileInfo.Type
		tag := fileInfo.Tag //获取tag名称
		name := tag.Get("json")
		if name == "" {
			name = strings.ToLower(fileInfo.Name)
		}
		//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
		name = strings.Split(name, ",")[0]

		fmt.Println("输出正常的字段名", name)
		//		fmt.Println("字段类型为", kind.Kind()) //获取内置基础类型
		//		fmt.Println("字段类型为2", kind)       //获取具体的类型比如[]uint8，但是上面获取到的是slice切片基础类型
		fmt.Println("字段类型为3", kind.String())

		if kind.Kind() == reflect.Struct {
			//结构体的字段是结构体，做特殊处理,暂时还没支持结构体嵌套结构体
			//			fmt.Println("字段类型是结构体")
			//			structObj := new()
			//			structObj := TypeTransverter()

		}
		//TODO其它类型判断，可以封装成一个方法

		name = fileInfo.Name
		if value, ok := s[name]; ok {
			//强制转化指定类型，并且赋值给结构体里面元素
			v := exchangeType(kind.String(), value)
			fmt.Println("转化后的数据类型为", reflect.ValueOf(v).Type())
			if v == nil {
				fmt.Println("转换类型失败", name)
				continue
			}
			reflect.ValueOf(obj).Elem().FieldByName(fileInfo.Name).Set(reflect.ValueOf(v))
		}
	}
}

func exchangeType(t string, v string) interface{} {
	switch t {
	case "int":
		b, error := strconv.Atoi(v)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		return b
	case "int8":
		b, error := strconv.ParseInt(v, 10, 8)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		return int8(b)
	case "int16":
		b, error := strconv.ParseInt(v, 10, 16)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		return int16(b)
	case "int32":
		b, error := strconv.ParseInt(v, 10, 32)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		return int32(b)
	case "int64":
		b, error := strconv.ParseInt(v, 10, 64)
		if error != nil {
			fmt.Println("字符串转换成整数失败")
		}
		return b
	case "uint":
	case "uint8":
	case "uint16":
	case "uint32":
	case "uint64":
	case "uintptr":
	case "float32":
	case "float64":
	case "complex64":
	case "complex128":
	case "array":
	case "chan":
	case "func":
	case "interface":
	case "map":
	case "ptr":
	case "slice":
	case "string":
		return v
	case "struct":
	case "unsafePointer":
	default:
		return nil
	}
	return nil

}
