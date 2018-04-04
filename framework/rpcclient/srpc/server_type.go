package srpc

//服务类型
const (
	GAMESER          = 1  //GameSer
	ALLOC            = 2  //Alocc
	ACCESS           = 3  //Access
	MSEVER           = 4  //Money
	USERSVER         = 5  //User
	DISPATCH         = 6  //Dispatch
	LOGIN            = 7  //LOGIN
	PHPAGENT         = 8  //PhpAgent
	ROBOTSEVER       = 9  //RobotServer
	PRIVATEROOMAGENT = 10 //Private room agent
	MATCHSERVER      = 11 //Match Server
	TEAMMGR          = 12 //TEAMMGRMgr
	CONFIGMANAGER    = 13 //Config Manager Server
	LOGSERVER        = 14 //Config Manager Server
	NEWMATCHSERVER   = 15
	UPDATESER        = 16
	NEWBACKSERVER    = 17
	STATSER          = 18
	PROPSMANAGERSER  = 19
	MSGSERVER        = 20 //Msg Server
	UPGRADE_CLIENT   = 21 //用于发送更新请求的客户端

	BACKUPSER     = 50
	COINCACHESER  = 51
	MATCHAGENT    = 52 //match list server
	CONTROLCENTER = 53 //control center
)

//不同的传输对象
const (
	TRANS_P2P       = 0 //点到点
	TRANS_P2G       = 1 //点到组
	TRANS_BROADCAST = 2 //同型广播
	TRANS_BYCMD     = 3 //根据关键字转发
	TRANS_COOK      = 4 //获取消息头cook转发
)

//区分不同的包。比如：在input那里收到的包有注册成功的消息包，rpc响应包等等等等
const (
	DISPATCH_REG_SER                 = 0x101 //表示向Net层注册的包
	DISPATCH_KEEP_ALIVE_SER          = 0x102 //主动向Net层发起心跳包
	GET_GAME_INFO                    = 0x103
	DISPATCH_SENDINFO_SER            = 0x104 //收到Net层发过来的0x104剔除冲突服务
	RESPONSE_DISPATCH_KEEP_ALIVE_SER = 0x105 //Net层发过来的心跳包
	DISPATCH_REG_RSP                 = 0x106 //注册响应包
	NOTIFY_OTHER_SER_RESTART         = 0x107
	DISPATCH_REG_RSP_OLD             = 0x108
	RPC_REQUEST_PACKAGE              = 0x109 //发起rpc 请求包
	RPC_RESPONSE_PACKET              = 0x10A //接收rpc 回复包
	CLIENT_TRANSFORM_PACKET          = 0x10B //发往客户端的转发包装头
	CLIENT_SPLIT_PACKET              = 0x10C //发往客户端的分包
	RPC_NOTIFY_PACKET                = 0x10D //推送消息
	DISPATCH_SER_NOT_FIND            = 0x10E //dispatch回应服务不存在的消息
	RPC_METHOD_NOTIFY_PACKET         = 0x10F //rpc方法调用型推送消息
	DISPATCH_REG_METHOD_REQ          = 0x201 //向Net层注册服务名和方法名的包
	DISPATCH_REG_METHOD_RSP          = 0x202 //接收0x201注册结果的包
)
