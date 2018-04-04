package controller

import (
	"net/http"
	"io"
	"putil/log"
	"dfqp/pg-http/service"
	"net"
	"time"
	"bytes"
	"mime/multipart"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"dfqp/proto/user"
	"strconv"
	"fmt"
)

type Sizer interface {
	Size() int64
}

type UploadApi struct {
	base
}

//上传头像
func (this *UploadApi) UpAvatar(writer http.ResponseWriter, request *http.Request) {
	if "POST" == request.Method {
		requestTime := fmt.Sprintf("%d", time.Now().Unix())
		request.ParseMultipartForm(5<<20) //5M
		file, header, err := request.FormFile("uploadAvatar")
		if err != nil {
			plog.Debug("get form file fail:", err)
			this.output(writer, 1020, "", "")
			return
		}
		defer file.Close()
		if sizeInterface, ok := file.(Sizer); ok {
			size := sizeInterface.Size()
			if size > 1048576 { //文件大小是否超过1M
				plog.Debug("get file size is:", size)
				this.output(writer, 1021, "", "")
				return
			}
		} else {
			plog.Debug("get file size fail")
			this.output(writer, 1020, "", "")
			return
		}

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, err := w.CreateFormFile("icon", header.Filename)
		_, err = io.Copy(fw, file)
		if err != nil {
			this.output(writer, 1020, "", "")
			return
		}

		cid := request.Form.Get("mid")
		err = w.WriteField("mid", cid)
		if err != nil {
			this.output(writer, 1020, "", "")
			return
		}
		err = w.WriteField("area", service.ServiceConf["app"].String("app.area"))
		if err != nil {
			this.output(writer, 1020, "", "")
			return
		}

		upload := service.ServiceConf["cdn"].String(service.Runmode+"::cdn.upload")
		url := service.ServiceConf["cdn"].String(service.Runmode+"::cdn.url")

		req, err := http.NewRequest("POST", upload, &b)
		if err != nil {
			this.output(writer, 1020, "", "")
			return
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(network, addr, time.Second*5)
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second*60))
					return conn, nil
				},
				ResponseHeaderTimeout:time.Second * 5,
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			plog.Debug("client do err:", err)
			this.output(writer, 1020, "", "")
			return
		}
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			this.output(writer, 1020, "", "")
			return
		}
		plog.Debug("result is:", string(result))
		json := gjson.Parse(string(result))
		flag := json.Get("flag").Int()
		if flag == 1 { //成功
			icon := url + json.Get("icon").String()+"?v="+requestTime
			iconBig := url + json.Get("iconbig").String()+"?v="+requestTime
			//rpc修改user数据库
			newCid, err := strconv.ParseInt(cid, 10, 64)
			if err != nil {
				plog.Debug("cid type transfer fail:", err)
				this.output(writer, 1020, "", "")
				return
			}
			ret := this.modifyUserInfo(newCid, icon, iconBig)
			if ret {
				data := map[string]string{
					"iconBig" : iconBig,
				}
				this.output(writer, 0, data, "")
			} else {
				this.output(writer, 1020, "", "")
			}
			return
		} else {
			this.output(writer, 1020, "", "")
			return
		}
	} else {
		//测试用
		writer.Header().Add("Content-Type", "text/html")
		writer.WriteHeader(200)
		html := `
<form enctype="multipart/form-data" action="/uploadAvatar" method="POST">
    Send this file: <input name="uploadAvatar" type="file" />
<input name="mid" value="224352"/>
    <input type="submit" value="Send File" />
	</form>
`
		io.WriteString(writer, html)
		return
	}
}
//发起rpc请求
func (this *UploadApi) modifyUserInfo(cid int64, icon string, iconBig string) bool {
	req := new(pgUser.ModifyUserInfoRequest)
	req.Mid = cid
	req.Icon = icon
	req.IconBig = iconBig
	reqBytes, err := req.Marshal()
	if err != nil {
		plog.Debug("marshal=========", err)
		return false
	}
	response := service.Client.SendAndRecvRespRpcMsg("pgUser.modifyUserInfo", reqBytes, 2000, 0)
	if response.ReturnCode != 0 {
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
		return false
	} else {
		arithResp := new(pgUser.ModifyUserInfoResponse)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return value  = ", arithResp)
		return true
	}
}
