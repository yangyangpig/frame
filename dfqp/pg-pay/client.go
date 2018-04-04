package main

import (
	"dfqp/proto/pay"
	"fmt"
	"framework/rpcclient/core"
	"putil/log"
	"runtime"
)

type MyClientDispatch struct {
	client *rpcclient.RpcCall
}

func (dis *MyClientDispatch) RpcRequest(req *rpcclient.RpcRecvReq, body []byte) {
	plog.Debug(req)
	plog.Debug(body)

	responsestr := "hi I'm Server return your data"
	dis.client.SendPacket(req, []byte(responsestr))
}

type orderTyep struct {
	Pdealno  string `json:"pdealno"`
	Receipt  string `json:"receipt"`
	Sandbox  int    `json:"sandbox"`
	BundleId string `json:"bundleId"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	quitFlag := make(chan bool)
	//实例化rpc
	client, err := rpcclient.NewRpcCall()
	if err != nil {
		plog.Fatal("fatal")
		return
	}

	err = client.RpcInit("./conf/clientapp.conf")
	if err != nil {
		plog.Debug("rpc 初始化异常")
		return
	}

	//服务基础配置
	mydisp := new(MyClientDispatch)
	mydisp.client = client
	err = client.LaunchRpcClient(mydisp)
	client.WriteNormalLog("test", "here is normal log!")
	client.WriteRealLog("test", "here is real log!")
	client.WriteDebugLog("test", "here is debug log!")
	//err = client.LaunchRpcClient("192.168.100.126", 9081, mydisp)
	if err != nil {
		plog.Fatal("lauch failed", err)
		return
	}
	plog.Debug("LaunchRpcClient succeed!!!")
	//超载包后面加一个正常包
	json := `{"apple_productid":"0","do":"suc","ext":"","game_item_id":"5904","mid":"20549022","pamount":"0.5","pamount_change":"0","pamount_rate":"0.155446","pamount_unit":"CNY","pamount_usd":"0.08","paychips_v2":"0","paycoins":"0","payconfid":"0","payprod_v2":"0","pc_appid":"1770","pc_rate":"0.155446","pc_sid":"7","pc_time":"1521215991","pdealno":"4200000054201803160048617074","pendtime":"1521215999","pid":"882316184","pmode":"431","pnum_v2":"1","pstarttime":"1521215991","request_ip":"127.0.0.1:54728","sign":"eba52662704224f5f5edc9bf5f5a5cd1","sign_v2":"f03cf56615fcdfabecba079156468760","sign_v3":"a0f315e2b298a1907dcedca7d4fceb47","sitemid":"","time":"1521215999"}`
	request := new(pgPay.SendMoneyRequest)
	request.Content = json
	//request.Pid = "000773000099BYORDFLG003152137300"
	//reportData := new(orderTyep)
	//reportData.Pdealno = "1000000382669314"
	//reportData.Receipt = "ewoJInNpZ25hdHVyZSIgPSAiQTRNUUxvZUdLN2RTUk1vMmZuQVZjNm1RRWJHWUp4S01QbDdDcytoQkdYRUdWYVlDeFltL2pFQjBiUVp1bDBGdHBpbHZNZHVMeVBxSmt0SE5KMEtGQzFWdzdWUGg4VW9VR0QwVGNoUmc1ak5HMnBmS0FBaGU1Z3BDTWMwZTJLY01ld2F3VzhVOWR4SkpkczZhR2lsazBOYnc3NHBGWWhRaGE4eHhROWlJUUs4OG5JZTFhazJ4V0FVV2RqSFZsZlVZUUNPc3RpNHNmL3N3bFdPdkZmbWhkSWZkb3c3M0syNWZvK1oySUtvZVFIQTdlZWRWbFdZYTZDNDFHaE5LanVrZVZxVjhwTEJKRnFnakdqeW1SSlc1Ty96czg1N0thQ1YzU0NVUDIwWnJ3ZnBSRGx5QVdoRTBESzJHRXd6NFJpMzRuWDg1bytSVFpmcjZyNUxqd1BlMkVsOEFBQVdBTUlJRmZEQ0NCR1NnQXdJQkFnSUlEdXRYaCtlZUNZMHdEUVlKS29aSWh2Y05BUUVGQlFBd2daWXhDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFLREFwQmNIQnNaU0JKYm1NdU1Td3dLZ1lEVlFRTERDTkJjSEJzWlNCWGIzSnNaSGRwWkdVZ1JHVjJaV3h2Y0dWeUlGSmxiR0YwYVc5dWN6RkVNRUlHQTFVRUF3dzdRWEJ3YkdVZ1YyOXliR1IzYVdSbElFUmxkbVZzYjNCbGNpQlNaV3hoZEdsdmJuTWdRMlZ5ZEdsbWFXTmhkR2x2YmlCQmRYUm9iM0pwZEhrd0hoY05NVFV4TVRFek1ESXhOVEE1V2hjTk1qTXdNakEzTWpFME9EUTNXakNCaVRFM01EVUdBMVVFQXd3dVRXRmpJRUZ3Y0NCVGRHOXlaU0JoYm1RZ2FWUjFibVZ6SUZOMGIzSmxJRkpsWTJWcGNIUWdVMmxuYm1sdVp6RXNNQ29HQTFVRUN3d2pRWEJ3YkdVZ1YyOXliR1IzYVdSbElFUmxkbVZzYjNCbGNpQlNaV3hoZEdsdmJuTXhFekFSQmdOVkJBb01Da0Z3Y0d4bElFbHVZeTR4Q3pBSkJnTlZCQVlUQWxWVE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBcGMrQi9TV2lnVnZXaCswajJqTWNqdUlqd0tYRUpzczl4cC9zU2cxVmh2K2tBdGVYeWpsVWJYMS9zbFFZbmNRc1VuR09aSHVDem9tNlNkWUk1YlNJY2M4L1cwWXV4c1FkdUFPcFdLSUVQaUY0MWR1MzBJNFNqWU5NV3lwb041UEM4cjBleE5LaERFcFlVcXNTNCszZEg1Z1ZrRFV0d3N3U3lvMUlnZmRZZUZScjZJd3hOaDlLQmd4SFZQTTNrTGl5a29sOVg2U0ZTdUhBbk9DNnBMdUNsMlAwSzVQQi9UNXZ5c0gxUEttUFVockFKUXAyRHQ3K21mNy93bXYxVzE2c2MxRkpDRmFKekVPUXpJNkJBdENnbDdaY3NhRnBhWWVRRUdnbUpqbTRIUkJ6c0FwZHhYUFEzM1k3MkMzWmlCN2o3QWZQNG83UTAvb21WWUh2NGdOSkl3SURBUUFCbzRJQjF6Q0NBZE13UHdZSUt3WUJCUVVIQVFFRU16QXhNQzhHQ0NzR0FRVUZCekFCaGlOb2RIUndPaTh2YjJOemNDNWhjSEJzWlM1amIyMHZiMk56Y0RBekxYZDNaSEl3TkRBZEJnTlZIUTRFRmdRVWthU2MvTVIydDUrZ2l2Uk45WTgyWGUwckJJVXdEQVlEVlIwVEFRSC9CQUl3QURBZkJnTlZIU01FR0RBV2dCU0lKeGNKcWJZWVlJdnM2N3IyUjFuRlVsU2p0ekNDQVI0R0ExVWRJQVNDQVJVd2dnRVJNSUlCRFFZS0tvWklodmRqWkFVR0FUQ0IvakNCd3dZSUt3WUJCUVVIQWdJd2diWU1nYk5TWld4cFlXNWpaU0J2YmlCMGFHbHpJR05sY25ScFptbGpZWFJsSUdKNUlHRnVlU0J3WVhKMGVTQmhjM04xYldWeklHRmpZMlZ3ZEdGdVkyVWdiMllnZEdobElIUm9aVzRnWVhCd2JHbGpZV0pzWlNCemRHRnVaR0Z5WkNCMFpYSnRjeUJoYm1RZ1kyOXVaR2wwYVc5dWN5QnZaaUIxYzJVc0lHTmxjblJwWm1sallYUmxJSEJ2YkdsamVTQmhibVFnWTJWeWRHbG1hV05oZEdsdmJpQndjbUZqZEdsalpTQnpkR0YwWlcxbGJuUnpMakEyQmdnckJnRUZCUWNDQVJZcWFIUjBjRG92TDNkM2R5NWhjSEJzWlM1amIyMHZZMlZ5ZEdsbWFXTmhkR1ZoZFhSb2IzSnBkSGt2TUE0R0ExVWREd0VCL3dRRUF3SUhnREFRQmdvcWhraUc5Mk5rQmdzQkJBSUZBREFOQmdrcWhraUc5dzBCQVFVRkFBT0NBUUVBRGFZYjB5NDk0MXNyQjI1Q2xtelQ2SXhETUlKZjRGelJqYjY5RDcwYS9DV1MyNHlGdzRCWjMrUGkxeTRGRkt3TjI3YTQvdncxTG56THJSZHJqbjhmNUhlNXNXZVZ0Qk5lcGhtR2R2aGFJSlhuWTR3UGMvem83Y1lmcnBuNFpVaGNvT0FvT3NBUU55MjVvQVE1SDNPNXlBWDk4dDUvR2lvcWJpc0IvS0FnWE5ucmZTZW1NL2oxbU9DK1JOdXhUR2Y4YmdwUHllSUdxTktYODZlT2ExR2lXb1IxWmRFV0JHTGp3Vi8xQ0tuUGFObVNBTW5CakxQNGpRQmt1bGhnd0h5dmozWEthYmxiS3RZZGFHNllRdlZNcHpjWm04dzdISG9aUS9PamJiOUlZQVlNTnBJcjdONFl0UkhhTFNQUWp2eWdhWndYRzU2QWV6bEhSVEJoTDhjVHFBPT0iOwoJInB1cmNoYXNlLWluZm8iID0gImV3b0pJbTl5YVdkcGJtRnNMWEIxY21Ob1lYTmxMV1JoZEdVdGNITjBJaUE5SUNJeU1ERTRMVEF6TFRFeklESXhPakV6T2pVNUlFRnRaWEpwWTJFdlRHOXpYMEZ1WjJWc1pYTWlPd29KSW5WdWFYRjFaUzFwWkdWdWRHbG1hV1Z5SWlBOUlDSXhPR0ptTjJGbFlXSXdaVFU1TVRsaE1qQXlZV0l3WTJJd01qRTJNamRtTlRKbVpERmtaV05qSWpzS0NTSnZjbWxuYVc1aGJDMTBjbUZ1YzJGamRHbHZiaTFwWkNJZ1BTQWlNVEF3TURBd01ETTRNalkyT1RNeE5DSTdDZ2tpWW5aeWN5SWdQU0FpTVNJN0Nna2lkSEpoYm5OaFkzUnBiMjR0YVdRaUlEMGdJakV3TURBd01EQXpPREkyTmprek1UUWlPd29KSW5GMVlXNTBhWFI1SWlBOUlDSXhJanNLQ1NKdmNtbG5hVzVoYkMxd2RYSmphR0Z6WlMxa1lYUmxMVzF6SWlBOUlDSXhOVEl4TURBd09ETTVNekF4SWpzS0NTSjFibWx4ZFdVdGRtVnVaRzl5TFdsa1pXNTBhV1pwWlhJaUlEMGdJa1kzTkRrMFFUWTFMVGxDT0RJdE5FWXhNeTA0UmtJd0xUazNPVFZGTlRNeU9ERXpOaUk3Q2draWNISnZaSFZqZEMxcFpDSWdQU0FpWTI5dExtSnZlV0ZoTGt4T1VWQXVNVGd3TURBd2MybHNkbVZ5WDFScFpYSXpJanNLQ1NKcGRHVnRMV2xrSWlBOUlDSXhNekE1TkRJek16TTBJanNLQ1NKaWFXUWlJRDBnSW1OdmJTNWliM2xoWVM1c2JuRndRVkJRU1VRaU93b0pJbWx6TFdsdUxXbHVkSEp2TFc5bVptVnlMWEJsY21sdlpDSWdQU0FpWm1Gc2MyVWlPd29KSW5CMWNtTm9ZWE5sTFdSaGRHVXRiWE1pSUQwZ0lqRTFNakV3TURBNE16a3pNREVpT3dvSkluQjFjbU5vWVhObExXUmhkR1VpSUQwZ0lqSXdNVGd0TURNdE1UUWdNRFE2TVRNNk5Ua2dSWFJqTDBkTlZDSTdDZ2tpYVhNdGRISnBZV3d0Y0dWeWFXOWtJaUE5SUNKbVlXeHpaU0k3Q2draWNIVnlZMmhoYzJVdFpHRjBaUzF3YzNRaUlEMGdJakl3TVRndE1ETXRNVE1nTWpFNk1UTTZOVGtnUVcxbGNtbGpZUzlNYjNOZlFXNW5aV3hsY3lJN0Nna2liM0pwWjJsdVlXd3RjSFZ5WTJoaGMyVXRaR0YwWlNJZ1BTQWlNakF4T0Mwd015MHhOQ0F3TkRveE16bzFPU0JGZEdNdlIwMVVJanNLZlE9PSI7CgkiZW52aXJvbm1lbnQiID0gIlNhbmRib3giOwoJInBvZCIgPSAiMTAwIjsKCSJzaWduaW5nLXN0YXR1cyIgPSAiMCI7Cn0"
	//reportData.Sandbox = 1
	//reportData.BundleId = "com.boyaa.lnqpAPPID"
	//tmpData,_ := json.Marshal(reportData)
	//request.Content = string(tmpData)

	req_bytes, err := request.Marshal()
	if err != nil {
		fmt.Println(req_bytes)
	}
	//		"arith"表示服务名
	//		Add表示调用的方法名
	//		req_bytes表示请求的参数（经过了protobuf的marshal后）
	//		5000表示5000毫秒后无响应就超时！
	response := client.SendAndRecvRespRpcMsg("pgPay.sendMoney", req_bytes, 5000, 0)
	if response.ReturnCode != 0 {
		//rpc返回结果异常
		//RPC_RESPONSE_COMPLET         = 0 //完成
		//	RPC_RESPONSE_TIMEOUT         = 1 //超时
		//	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		//	RPC_RESPONSE_NETERR          = 3 //发生网络错误
		//	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		plog.Debug("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {
		arithResp := new(pgPay.OrderResponse)
		arithResp.Unmarshal(response.Body)
		plog.Debug("return value  = ", arithResp)
	}

	<-quitFlag
}
