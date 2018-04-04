package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	_ "github.com/astaxie/beego/httplib"
)

type res struct {
	Code int
	Msg  string
	Data string
}

func SendDataCli(param []byte, servname string, funcname string) (lastBody string, err error) {

	p := base64.StdEncoding.EncodeToString(param)
	fmt.Println("base64", p)
	postParam := url.Values{
		"params":     {p},
		"servername": {servname},
		"funcname":   {funcname},
	}
	resp, err := http.PostForm("http://www.rpctesttool.com:3001/launch", postParam)
	fmt.Println("resp", resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	resPrt := &res{}
	json.Unmarshal(body, resPrt)
	data, _ := base64.StdEncoding.DecodeString(resPrt.Data) //data是[]byte
	tmp := ProcessRes(data, servname, funcname)
	resPrt.Data = tmp
	if err != nil {
		fmt.Println(err)
		return
	}
	b, _ := json.Marshal(resPrt)
	lastBody = string(b)
	return

}
