package service

import (
	"github.com/astaxie/beego/orm"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"dfqp/pg-goods/controller"
	"dfqp/pg-goods/entity"
	"dfqp/pg-goods/model"
	"framework/rpcclient/core"
	"github.com/astaxie/beego/config"
	"os"
	"putil/cache"
	"time"
)

var (

	tablePrefix  = make(map[string]string) // 表前缀map

)

// Service 服务接口
type Service interface {
	Init()
	Reload() error
	Run() (int, error)
	Exit()
	Status() map[string]string
}

const (
	configSectionRPC   = "rpc"
	configSectionLOG   = "log"
	configSectionDB    = "db"
	configSectionCache = "cache"
	configSectionRedis = "redis"
)

// PGService pg server
type PGService struct {
	// 服务名称
	name string
	// 服务进程id
	pid int
	// 服务运行环境
	env string
	// 服务所有配置
	config config.Configer
	// 服务状态
	status map[string]string

	rpcServer *rpcclient.RpcCall
}

// SetConfig set config
func (s *PGService) SetConfig(c config.Configer) {
	println(s)
	s.config = c
}

// Init 初始化服务有次序的
func (s *PGService) Init() {
	// TODO 可优化异步初始化，发挥多核优势加速服务初始化
	s.initService()
	// 初始化数据库连接，以及注册好orm mode
	s.initStorage()
	// 初始化所有的model
	s.initModels()
	s.initRPCServer()
}

// init service
func (s *PGService) initService() {
	config, err := s.config.GetSection(configSectionRPC)
	if err != nil {
		panic("init " + s.name + " rpc server error: " + err.Error())
	}

	if env, exist := config["rpc.run_mode"]; exist {
		if _, exist := entity.EnvValue[env]; !exist {
			s.env = entity.EnvProduct
		} else {
			s.env = env
		}
	} else {
		s.env = entity.EnvProduct
	}

	if name, exist := config["rpc.name"]; exist {
		s.name = name
	} else {
		panic("init rpc server error: not exist rpc.name")
	}
}

// 初始化 models
func (s *PGService) initModels() {
	// TODO init models
	o := orm.NewOrm()
	model.GoodsLabel.SetOrm(o)
	model.Goods.SetOrm(o)
	model.Prices.SetOrm(o)
	model.Bag.SetOrm(orm.NewOrm())
	model.ExHistory.SetOrm(orm.NewOrm())
	model.ExHistory.SetEvn(s.env)

}

// init rpc server
func (s *PGService) initRPCServer() {
	var err error
	// 实例化rpc
	client, err := rpcclient.NewRpcCall()
	controller.SetRPCClient(client)

	s.rpcServer = client
	// TODO 优化
	fileName := "/Users/liyufeng/Documents/boyaa-new-hall/src/dfqp/pg-goods/conf/app.conf"
	s.rpcServer.RpcInit(fileName)
	if err != nil {
		panic("create rpc server failed:" + err.Error())
		return
	}

	println("--register rpc request router ...")
	// TODO register rpc request handle functions
	// s.rpcServer.RpcHandleFunc("pgGoods.create", controller.CreateGoods)
	// s.rpcServer.RpcHandleFunc("pgGoods.all", controller.All)
	s.rpcServer.RpcHandleFunc("pgGoods.goodsInfo", controller.GetGoodsInfo)
	// s.rpcServer.RpcHandleFunc("pgGoods.exchangeTypeInfo", controller.ExchangeTypeInfo)
	// userBag := rpcclient.RpcHandleFunc(controller.UserBag)
	// s.rpcServer.RpcHandle("pgGoods.bag", controller.CheckUserOnLine(userBag))
	s.rpcServer.RpcHandleFunc("pgGoods.bag", controller.UserBag)
	s.rpcServer.RpcHandleFunc("pgGoods.use", controller.Use)
	s.rpcServer.RpcHandleFunc("pgGoods.synthesis", controller.Synthesis)
	s.rpcServer.RpcHandleFunc("pgGoods.exchangeRealGoods", controller.ExchangeRealGoods)
	s.rpcServer.RpcHandleFunc("pgGoods.exchangeTelFee", controller.ExchangeTelFee)
	s.rpcServer.RpcHandleFunc("pgGoods.exchangeHistory", controller.ExHistory)
	// s.rpcServer.RpcHandleFunc("pgGoods.bagItem", controller.BagItem)
}

// initStorage init storage
func (s *PGService) initStorage() {
	// 全局配置orm设置本地时区
	orm.Debug = true
	orm.DefaultTimeLoc = time.UTC
	s.initDB(entity.ModeDefault, entity.ModeGoods)
	// TODO 使用常量
	s.initCache("default", "bag")
	// 向orm注册model
	registerModes()
}

func (s *PGService) initDB(biz ...string) {

	// 获取数据库配置
	config, err := s.config.GetSection(configSectionDB)
	if err != nil {
		panic("load database config failed")
	}
	for _, v := range biz {
		driverField := v + ".db.driver"
		driver, exist := config[driverField]
		if !exist {
			panic("Get " + v + " 's db driver failed")
		}
		hostFiled := v + "." + driver + ".host"
		userFiled := v + "." + driver + ".user"
		pwdFiled := v + "." + driver + ".pwd"
		databaseFiled := v + "." + driver + ".database"
		// tablePrefixFiled := v + "." + driver + ".table_prefix"

		host, exists := config[hostFiled]
		if !exists {
			panic("Get " + v + " 's " + driver + " host failed")
		}
		user, exists := config[userFiled]
		if !exists {
			panic("Get " + v + " 's " + driver + " user failed")
		}
		pwd, exists := config[pwdFiled]
		if !exists {
			panic("Get " + v + " 's " + driver + " pwd failed")
		}
		dbName, exists := config[databaseFiled]
		if !exists {
			panic("Get " + v + " 's " + driver + " dbName failed")
		}

		if driver == "mysql" {
			// 注册mysql
			maxIdle := 10 // 最大空闲连接
			maxConn := 10 // 最大连接数
			dataSource := user + ":" + pwd + "@tcp(" + host + ")/" + dbName + "?charset=utf8"
			println("--register mysql, 业务: " + v + ", dataSource=" + dataSource)
			orm.RegisterDataBase(v, "mysql", dataSource, maxIdle, maxConn)

		}

	}
}

// initCache
func (s *PGService) initCache(biz ...string) {
	// 获取数据库配置
	config, err := s.config.GetSection(configSectionCache)
	if err != nil {
		panic("load database config failed")
	}
	for _, v := range biz {
		driverField := v + ".cache.driver"
		driver, exist := config[driverField]
		if !exist {
			panic("Get " + v + " 's cache driver failed")
		}
		addr := v + "." + driver + ".addr"
		address, exists := config[addr]
		if !exists {
			panic("Get " + v + " 's " + driver + " host failed")
		}

		if driver == "redis" {
			//注册cache
			pgCache := cache.NewPGCache()
			pgCache.SetAddr(address)
			pgCache.Connect()
			cache.RegisterCache(v, pgCache)

		}

	}
}

// Register mode
func registerModes() {
	orm.RegisterModel(new(entity.PgGoodsLabel))
	orm.RegisterModel(new(entity.PgGoods))
	orm.RegisterModel(new(entity.PgBagItem))
	//TODO register models ...

}

// Reload 重新加载
func (s *PGService) Reload() error {

	// reload do something
	return nil
}

// RPCServer return rpc server
func (s *PGService) RPCServer() *rpcclient.RpcCall {
	return s.rpcServer
}

// Run 运行服务
func (s *PGService) Run() (pid int, err error) {
	// dis := new(pgGoodsDispatch)
	err = s.rpcServer.LaunchRpcClient(nil)
	if err != nil {
		panic("launch rpc server failed:" + err.Error())
		return
	}
	s.pid = os.Getpid()
	return s.pid, nil
}

// Exit 退出服务
func (s *PGService) Exit() {
	println(s.name + " service quiting ...")
	//TODO 退出时处理 清理对象
	println(s.name + " service quited")
}

// Status 返回服务状态
func (s *PGService) Status() map[string]string {

	return s.status
}

// Name 返回服务名称
func (s *PGService) Name() string {
	return s.name
}

// NewPGGoods 实例化pg-goods rpc 服务
func NewPGGoods() *PGService {
	return new(PGService)
}
