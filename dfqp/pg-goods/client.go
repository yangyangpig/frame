package main

import (
	"fmt"
	"framework/rpcclient/core"
	"dfqp/proto/goods"
	"time"
	"encoding/json"
)

var client *rpcclient.RpcCall

func init() {
	// logFile := "/Users/liyufeng/Documents/boyaa-new-hall/src/dfqp/pg-goods/log/client.log"
	// plog.SetOutPutFileRotate(logFile, 1*1024*1024, 3, 30)
	// plog.LogRegister("client")
}

// create goods
func createGoods() {
	request := new(ptGoods.CreateGoodsRequest)
	request.Name = "goodsName"
	request.Desc = "goods desc"
	request.Type = uint32(0)
	request.Img = "http://cdn.boyaa.com/goods/1.png"
	request.Conditions = "{\"time\":6}"
	request.Price = uint32(5)
	request.PriceOrg = uint32(3)
	request.GoodsExt = "goods_ext"
	request.ExpireTime = 10000
	request.CreateBy = 1
	request.CreateTime = time.Now().Unix()
	request.Status = ptGoods.UNAVAILABLE

	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	response := client.SendAndRecvRespRpcMsg("pgGoods.create", reqBytes, 5000, 0)
	println("------------send & response ---------------")
	if response.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// 	RPC_RESPONSE_TIMEOUT         = 1 //超时
		// 	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// 	RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// 	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {

	}
}

func getGoodsTypes() {
	request := new(ptGoods.CreateGoodsRequest)
	request.Name = "goodsName"
	request.Desc = "goods desc"
	request.Type = uint32(0)
	request.Img = "http://cdn.boyaa.com/goods/1.png"
	request.Conditions = "{\"time\":6}"
	request.Price = uint32(5)
	request.PriceOrg = uint32(3)
	request.GoodsExt = "goods_ext"
	request.ExpireTime = 10000
	request.CreateBy = 1
	request.CreateTime = time.Now().Unix()
	request.Status = ptGoods.UNAVAILABLE

	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	response := client.SendAndRecvRespRpcMsg("pgGoods.create", reqBytes, 5000, 0)
	println("------------send & response ---------------")
	if response.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// 	RPC_RESPONSE_TIMEOUT         = 1 //超时
		// 	RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// 	RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// 	RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", response.ReturnCode, " return err = ", response.Err)
	} else {

	}
}

func bag() {
	request := new(ptGoods.Mid)
	request.Id = 1
	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.bag"
	println("rpc-client request: ", method)
	response := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & response ---------------")
	if response.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", response.ReturnCode, " return err = ")
		if response.Err != nil {
			println(response.Err.Error())
		}
	} else {

		bagItemsResponse := new(ptGoods.BagItemsResponse)
		bagItemsResponse.Unmarshal(response.Body)
		bytes, _ := json.Marshal(bagItemsResponse)
		println(string(bytes))
	}
}

func testProto() {
	r := new(ptGoods.RPCErrorResponse)
	r.Status = 1
	r.Msg = "测试"
	bytes, _ := r.Marshal()

	r2 := new(ptGoods.GoodsUseResponse)
	if e := r2.Unmarshal(bytes); e != nil {
		println(e.Error())
		return
	}
	println(r2.Status, r2.Msg)
}

// 测试用例 使用鲜花
func testUseFlower() {
	request := new(ptGoods.ExchangeRealGoodsRequest)
	request.Mid = 1
	request.GoodsId = 7
	request.UgId = 4
	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.use"
	println("rpc-client request: ", method)
	rpcResponse := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & rpcResponse ---------------")
	if rpcResponse.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", rpcResponse.ReturnCode, " return err = ")
		if rpcResponse.Err != nil {
			println(rpcResponse.Err.Error())
		}
	} else {

		response := new(ptGoods.GoodsUseResponse)
		response.Unmarshal(rpcResponse.Body)
		bytes, _ := json.Marshal(response)
		println(string(bytes))
	}

}

func testExReal() {
	request := new(ptGoods.ExchangeRealGoodsRequest)
	request.Mid = 1
	request.UgId = 9
	request.GoodsId = 4
	request.RealName = "张三"
	request.Phone = "13692202344"
	request.Addr = "深圳市南山区TCL"
	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.exchangeRealGoods"
	println("rpc-client request: ", method)
	rpcResponse := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & rpcResponse ---------------")
	if rpcResponse.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", rpcResponse.ReturnCode, " return err = ")
		if rpcResponse.Err != nil {
			println(rpcResponse.Err.Error())
		}
	} else {

		response := new(ptGoods.GoodsUseResponse)
		response.Unmarshal(rpcResponse.Body)
		bytes, _ := json.Marshal(response)
		println(string(bytes))
	}
}
func testTelFee() {
	request := new(ptGoods.ExChangeTelFeeRequest)
	request.Mid = 1
	request.UgId = 2
	request.GoodsId = 5
	request.Phone = "13692202344"
	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.exchangeTelFee"
	println("rpc-client request: ", method)
	rpcResponse := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & rpcResponse ---------------")
	if rpcResponse.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", rpcResponse.ReturnCode, " return err = ")
		if rpcResponse.Err != nil {
			println(rpcResponse.Err.Error())
		}
	} else {

		response := new(ptGoods.GoodsUseResponse)
		response.Unmarshal(rpcResponse.Body)
		bytes, _ := json.Marshal(response)
		println(string(bytes))
	}
}

// 合成
func testSynthesis() {
	request := new(ptGoods.SynthesisRequest)
	request.Mid = 1
	request.UgId = 10
	request.GoodsId = 8
	request.Num = 100
	request.TargetId = 5

	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.synthesis"
	println("rpc-client request: ", method)
	rpcResponse := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & rpcResponse ---------------")
	if rpcResponse.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", rpcResponse.ReturnCode, " return err = ")
		if rpcResponse.Err != nil {
			println(rpcResponse.Err.Error())
		}
	} else {

		response := new(ptGoods.GoodsUseResponse)
		response.Unmarshal(rpcResponse.Body)
		bytes, _ := json.Marshal(response)
		println(string(bytes))
	}
}

// 合成
func testExHistory() {
	request := new(ptGoods.ExHistoryRequest)
	request.Mid = 1
	request.New = 1
	request.PreIndex = 0
	request.PageSize = 30

	reqBytes, err := request.Marshal()
	if err != nil {
		fmt.Println(reqBytes)
	}
	method := "pgGoods.exchangeHistory"
	println("rpc-client request: ", method)
	rpcResponse := client.SendAndRecvRespRpcMsg(method, reqBytes, 50000, 0)
	println("------------send & rpcResponse ---------------")
	if rpcResponse.ReturnCode != 0 {
		// rpc返回结果异常
		// RPC_RESPONSE_COMPLET         = 0 //完成
		// RPC_RESPONSE_TIMEOUT         = 1 //超时
		// RPC_RESPONSE_SENDFAILED      = 2 //发送错误
		// RPC_RESPONSE_NETERR          = 3 //发生网络错误
		// RPC_RESPONSE_TARGET_NOTFOUND = 4 //Net层没有发现目标实例
		println("rpc return code = ", rpcResponse.ReturnCode, " return err = ")
			if rpcResponse.Err != nil {
			println(rpcResponse.Err.Error())
		}
	} else {
		response := new(ptGoods.ExHistoryResponse)
		response.Unmarshal(rpcResponse.Body)
		bytes, _ := json.Marshal(response)
		println(string(bytes))
	}
}

func main() {

	var err error
	client, err = rpcclient.NewRpcCall()
	if err != nil {
		println("fatal")
		return
	}
	fileName := "/Users/liyufeng/Documents/boyaa-new-hall/src/dfqp/pg-goods/conf/clientapp.conf"
	err = client.RpcInit(fileName)
	if err != nil {
		println("rpc 初始化异常")
		return
	}

	// 服务基础配置
	err = client.LaunchRpcClient(nil)
	client.WriteNormalLog("test", "here is normal log!")
	client.WriteRealLog("test", "here is real log!")
	client.WriteDebugLog("test", "here is debug log!")

	if err != nil {
		println("lauch failed", err)
		return
	}
	println("LaunchRpcClient succeed!!!")
	// 		"arith"表示服务名
	// 		Add表示调用的方法名
	// 		req_bytes表示请求的参数（经过了protobuf的marshal后）
	// 		5000表示5000毫秒后无响应就超时！

	// createGoods()
	// bag()
	// testUseFlower()
	// testExReal()
	// testTelFee()

	// testSynthesis()
	testExHistory()

}
