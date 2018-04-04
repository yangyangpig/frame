package controllers

import (
	"Proto/sms"
	"time"
	"fmt"
	"PGLibrary"
	"net/http"
	"io/ioutil"
	"crypto/rc4"
	"github.com/tidwall/gjson"
	"github.com/astaxie/beego/logs"
)

type SmsController struct {}

//获取验证码demo
func (this *SmsController) GetCaptcha(req *pgSms.GetCaptchaRequest) *pgSms.GetCaptchaResponse {
	resp := new(pgSms.GetCaptchaResponse)
	if len(req.Phone) <= 0 {
		resp.Status = 1 //失败
		return resp
	}
	timeStamp := time.Now().Unix()
	secret := "dfqpjt$!@iz%s_=a*Aux#23!@#(No_PC"
	source := "1362118431"
	url := "flag=1&phone="+req.Phone+"&source="+source+"&timestamp="+fmt.Sprintf("%d", timeStamp)
	sig := PGLibrary.GetSha1(url+secret)
	url += "&signature="+sig

	rc4Obj, _ := rc4.NewCipher([]byte("by$!Gl_+d#f$%)Sk2=,>zI-l"))
	rs4Str := []byte(url)
	plaintext := make([]byte, len(rs4Str))
	rc4Obj.XORKeyStream(plaintext, rs4Str)
	ps := fmt.Sprintf("%x", plaintext)

	client := &http.Client{}

	ret, err := http.NewRequest("GET", "http://passport-debug.boyaa.com/user/token?ps="+ps, nil)
	if err != nil {
		logs.Error("SmsController.GetCaptcha Err: ", err.Error())
	}
	ret.Header.Set("User-Agent", "Boyaa Agent Alpha 0.0.1")
	re, err := client.Do(ret)
	if err != nil {
		logs.Error("SmsController.GetCaptcha Err: ", err.Error())
	}
	defer re.Body.Close()
	body, err := ioutil.ReadAll(re.Body)
	if err != nil {
		logs.Error("SmsController.GetCaptcha Err: ", err.Error())
	}
	jsonData := gjson.Parse(string(body))
	code := jsonData.Get("code").Int()
	result := jsonData.Get("result").Int()
	logs.Debug("return value: ", jsonData)
	if code != 200 || result != 100 {
		resp.Status = 1 //失败
		return resp
	} else {
		resp.Status = 0 //成功
		return resp
	}
}
