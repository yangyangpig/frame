package szmq

//rpc执行时间上报结构
type DcRpcMonitor struct {
	Channel_id       string `json:"channel_id"`       //渠道ID
	Act_time         int64  `json:"act_time"`         //流水时间（精确到毫秒）
	Server_name      string `json:"server_name"`      //服务名称
	Server_id        string `json:"server_id"`        //服务ID
	Func_name        string `json:"func_name"`        //方法或函数名
	Typ              string `json:"type"`             //探针类型
	Exec_microsecond int64  `json:"exec_microsecond"` //用户id
	Ext_data         string `json:"ext_data"`         //消息的业务序列号，主要用于日志bug查找
}
