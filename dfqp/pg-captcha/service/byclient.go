package service

import (
	"time"
	"sort"
	"fmt"
	"crypto/rc4"
	"net/http"
	"io/ioutil"
	"putil/log"
	"dfqp/lib"
	"errors"
)

const (
	userAgent = "Boyaa Agent Alpha 0.0.1"
)

var (
	configError = errors.New("配置错误")
)
//byClient sdk
type byClientService struct {
}

//get请求
func (this *byClientService) Get(api string, data map[string]interface{}) (string, error) {
	apiUrl := ServiceConf["sms"].String(Runmode+"::byClient.apiUrl")
	if apiUrl == "" {
		return "", configError
	}
	url := apiUrl + api + "?" + this.Sort(data)
	plog.Debug("url====", url)
	client := &http.Client{}

	ret, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	ret.Header.Set("User-Agent", userAgent)
	re, err := client.Do(ret)
	if err != nil {
		return "", err
	}
	defer re.Body.Close()
	body, err := ioutil.ReadAll(re.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
//排序
func (this *byClientService) Sort(data map[string]interface{}) string {
	newData := make(map[string]interface{})
	newData["source"] = ServiceConf["sms"].String("boyaauc.id")
	timeStamp := time.Now().Unix()
	newData["timestamp"] = timeStamp

	for i, v := range data {
		newData[i] = v
	}

	sortKeys := make([]string, 0)
	for k, _ := range newData {
		sortKeys = append(sortKeys, k)
	}
	sort.Strings(sortKeys)
	var buildStr string
	for i, k := range sortKeys {
		value := fmt.Sprintf("%v", newData[k])
		if value != "" {
			if i != (len(sortKeys) - 1) {
				buildStr = buildStr + k + "=" + value + "&"
			} else {
				buildStr = buildStr + k + "=" + value
			}
		}
	}
	plog.Debug("buildStr===", buildStr)
	sig := lib.GetSha1(buildStr+ServiceConf["sms"].String("boyaauc.secret"))
	buildStr += "&signature=" + sig

	rc4Obj, _ := rc4.NewCipher([]byte(ServiceConf["sms"].String(Runmode+"::byClient.secret")))
	rs4Str := []byte(buildStr)
	plaintext := make([]byte, len(rs4Str))
	rc4Obj.XORKeyStream(plaintext, rs4Str)
	ps := fmt.Sprintf("%x", plaintext)
	baseString := "ps=" + ps
	return baseString
}



