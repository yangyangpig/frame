package szk

import (
	"encoding/json"
	"errors"
	"fmt"
	//"log"
	"path/filepath"
	"putil/log"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	ZKCONN_CLOSE     int32 = 0 //连接关闭状态
	ZKCONN_CONNECTED int32 = 1 //已连接
)

//var zkIpPorts = []string{"192.168.201.80:2181", "192.168.200.19:2181", "192.168.200.144:2181"}
//var zkIpPorts = []string{"192.168.203.219:2181"}
//flags有4种取值：
//0:永久，除非手动删除
//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
//zk.FlagSequence  = 2:会自动在节点后面添加序号
//3:Ephemeral和Sequence，即，短暂且自动添加序号
//var flags int32 = zk.FlagSequence //zk.FlagSequence时生成的节点上有一段数字串：rpc0000000011
//var flags int32 = 0
//var acls = zk.WorldACL(zk.PermAll)
//var zpathHead = "/zkrpc"
var (
	flags              int32    = 0
	acls                        = zk.WorldACL(zk.PermAll)
	zpathHead                   = "/rpc"
	zkIpPorts          []string //通过配置文件读取
	registerRetryTimes int      //通过配置文件读取

	//zkIpPorts                = []string{"192.168.201.80:2181", "192.168.200.19:2181", "192.168.200.144:2181"}
	//registerRetryTimes       = 10 //多进程并发向zk注册的时候可能会出现失败的情况，需要重试N次
)

//zk客户端结构体
type SzkClient struct {
	ServiceAll  map[string]*zkServiceInfo //zk数据
	conn        *zk.Conn                  //与zk的连接
	connstatus  int32                     //连接状态
	Mutx        sync.Mutex                //并发控制
	haveWatched map[string]bool           //某一类服务的监听标记,只要服务开始监听了就为true
}

//一个服务的结构体，方便根据servicename+functionname找到svrtype+svrid；本地内存使用！
type zkServiceInfo struct {
	serviceName string             //服务名称
	vinfo       *zkSt              //type、id和版本信息
	funcinfo    map[int32][]string //版本对应的方法信息
}

//typeandids数据结构，json化之后存入zk(给zk使用)
type zkSt struct {
	Typ     int64   `json:"Typ"`     //服务类型
	Ver     int32   `json:"Ver"`     //服务版本号
	Ids     []int64 `json:"Ids"`     //服务最新的实例id片
	PrevVer int32   `json:"PrevVer"` //上个版本号
	PrevIds []int64 `json:"PrevIds"` //上个版本的实例id片 `json:"PrevIds,omitempty"`当空的时候不进json
}

//配置初始化
func init() {
	//iniconf, err := config.NewConfig("ini", "/home/JungleYe/go/src/framework/rpcclient/conf/app.conf") //文件的绝对路径ok
	//	plog.Debug("ssssssssssssssssssssssssssssssszk ", beego.BConfig.RunMode)

	iniconf, err := config.NewConfig("ini", "./conf/zk.conf")
	if err != nil {
		plog.Fatal("new config adapter err: ", err)
	}
	registerRetryTimes, _ = iniconf.Int("zk.retrytimes") //读取zk注册重试次数的配置
	zkIpPorts = iniconf.Strings("zk.zkipport")           //读取zk的ip端口配置
	//plog.Debug("ssssssssssssssssssssssssssssssszk ", zkIpPorts[0])
}

//辅助方法：判断方法是否存在，并根据方法找合适的版本(返回的是版本对应ids的map)
func getProperFuncs(funcname string, zkserviceinfoPtr *zkServiceInfo) map[int32][]int64 {
	var v2ids = make(map[int32][]int64)
	if zkserviceinfoPtr.vinfo.Ver > 0 && len(zkserviceinfoPtr.vinfo.Ids) > 0 { //最新版本有实例
		version := zkserviceinfoPtr.vinfo.Ver
		//检查方法存在,其版本对应的实例返回(考虑查询性能，可将slice转为map)
		for i := 0; i < len(zkserviceinfoPtr.funcinfo[version]); i++ {
			//plog.Debug(zkserviceinfoPtr.funcinfo[version][i])
			if funcname == zkserviceinfoPtr.funcinfo[version][i] {
				v2ids[version] = zkserviceinfoPtr.vinfo.Ids //选中(新版本实例中有这个方法)
			}
		}
	}

	//老版本中是否有该方法
	if zkserviceinfoPtr.vinfo.PrevVer > 0 && len(zkserviceinfoPtr.vinfo.PrevIds) > 0 { //最新版本有实例
		version := zkserviceinfoPtr.vinfo.PrevVer
		//检查方法存在,其版本对应的实例返回(考虑查询性能，可将slice转为map)
		for i := 0; i < len(zkserviceinfoPtr.funcinfo[version]); i++ {
			if funcname == zkserviceinfoPtr.funcinfo[version][i] {
				v2ids[version] = zkserviceinfoPtr.vinfo.PrevIds //选中(老版本实例中有这个方法)
			}
		}
	}
	return v2ids
}

//实例化szkclient
func NewSzkClient() (szkclient *SzkClient, err error) {
	szkclient = new(SzkClient)
	if szkclient == nil {
		return nil, errors.New("alloc memory failed")
	}
	szkclient.ServiceAll = make(map[string]*zkServiceInfo)
	szkclient.haveWatched = make(map[string]bool)
	return
}

//连接zk(并发线程存在时也只连接一次)
func (szkclient *SzkClient) szkconnect(servers []string, sessionTimeout time.Duration) (err error) {
	//本地内存找不到就从zk找，防止并发连接
	szkclient.Mutx.Lock()
	defer szkclient.Mutx.Unlock()

	if szkclient.connstatus == ZKCONN_CLOSE {
		//当zk连不上的时候，会每隔一秒钟调用一次！
		conn, _, cerr := zk.Connect(servers, sessionTimeout)
		if cerr != nil {
			err = cerr
			return //连接出错
		}
		szkclient.connstatus = ZKCONN_CONNECTED //存储连接状态
		szkclient.conn = conn                   //存储连接
	}
	return
}

/*
*	根据serviceName查找svrtype
 */
func (szkclient *SzkClient) GetSertypeByServiceName(svcName string) (svrtype int64, v2ids map[int32][]int64, err error) {
	//	//查找本地内存
	//	if _, ok := szkclient.ServiceAll[svcName]; ok {
	//		return
	//	}

	//	//连接
	//	err = szkclient.szkconnect(zkIpPorts, time.Second)
	//	if err != nil {
	//		return
	//	}

	//	//从zk中查找方法（TODO：多线程并发写内存的情况）
	//	spath := zpathHead + "/" + svcName + "/typeandids" //服务存放地址
	//	serviceExists, _, _ := szkclient.conn.Exists(spath)
	//	if serviceExists {
	//		tmpsvrandids, _, terr := szkclient.conn.Get(spath)
	//		if terr != nil {
	//			err = terr
	//			return //查询出错
	//		}
	//		var jsd_typeandids []int64
	//		json.Unmarshal(tmpsvrandids, &jsd_typeandids)
	//		_ = jsd_typeandids[0]
	//		rttype = jsd_typeandids[0]
	//		wholeids := jsd_typeandids[1:] //wholeids与rtids的区别在于，后者是前者的一个子集！假设某个服务新增了一个方法！则新方法的ids只能到方法层去找！

	//		//有其他方法存在（比如说之前只存了Arith的A方法，当访问Arith的B方法时就会触发此处）
	//		var funcmap map[string][]int64 = make(map[string][]int64)
	//		_, sok := szkclient.ServiceAll[svcName]
	//		if sok {
	//			funcmap = szkclient.ServiceAll[svcName].funcs
	//		}

	//		//zk查到之后，放到内存map：
	//		var zst zkServiceInfo
	//		zst.serviceName = svcName
	//		zst.svrtype = rttype
	//		zst.svrid = wholeids
	//		zst.funcs = funcmap
	//		szkclient.ServiceAll[svcName] = &zst

	//		//锁控制，防止一类服务的多个监听！
	//		szkclient.Mutx.Lock()
	//		//如果还没有监听，就开始监听
	//		if !(sok && szkclient.haveWatched[svcName]) {
	//			fmt.Println("I am here!!!!!!!!!!!!!!!!!!!!!!!")
	//			go szkclient.watchNodeDataChange(spath)  //奇怪：注释掉该行代码后
	//			go szkclient.watchChildrenChanged(spath) //监控服务节点数据变化和子节点变化「增删」
	//			szkclient.haveWatched[svcName] = true    //添加监听标记
	//		}
	//		szkclient.Mutx.Unlock()

	//	}
	return
}

/*
*	根据serviceName和functionName获取serverid、servertype和groupid(连接、监听都只建立一个！操作内存也加锁)
 */
func (szkclient *SzkClient) GetSerByNames(svcName string, funcName string) (svrtype int64, v2ids map[int32][]int64, err error) { //, group_ids []int64
	//本地内存查找含有该方法的实例！
	plog.Debug("已经进来zk")
	if zkserviceinfoPtr, ok := szkclient.ServiceAll[svcName]; ok {
		v2ids = getProperFuncs(funcName, zkserviceinfoPtr)
		svrtype = szkclient.ServiceAll[svcName].vinfo.Typ
		return
	}

	//连接
	err = szkclient.szkconnect(zkIpPorts, time.Second*2)
	if err != nil {
		return
	}

	//从zk中查找方法(步骤是：先把zk数据加载到内存，然后从内存去查找！)
	//锁控制，防止一类服务的多个监听！
	szkclient.Mutx.Lock()
	spath := zpathHead + "/" + svcName + "/typeandids" //服务存放地址
	zkserviceinfoPtr, memok := szkclient.ServiceAll[svcName]
	if !memok {

		//获取vinfo部分数据
		svcExists, _, terr := szkclient.conn.Exists(spath)
		if terr != nil || !svcExists {
			//如果节点不存在或者出错了直接返回！
			plog.Fatal("check err happened: ", terr)
			err = terr
			szkclient.Mutx.Unlock()
			return
		}
		tmptypeandids, _, terr := szkclient.conn.Get(spath)
		if terr != nil {
			err = terr
			szkclient.Mutx.Unlock()
			return
		}
		var jsd_vinfo *zkSt
		//plog.Debug("===================1", string(tmptypeandids))
		jsonerr := json.Unmarshal(tmptypeandids, &jsd_vinfo) //json解析到结构体中
		if jsonerr != nil {
			plog.Fatal("=================2", jsonerr)
		}
		//plog.Debug("=================3", jsd_vinfo)

		//funcinfo部分
		var funcinfo map[int32][]string = make(map[int32][]string)
		var versions = []int32{jsd_vinfo.Ver, jsd_vinfo.PrevVer}
		for _, v := range versions { //每个版本的方法都存一次
			if v > 0 {
				fpath1 := zpathHead + "/" + svcName + "/funcs/" + strconv.Itoa(int(v))
				//plog.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", fpath1)
				funcExists, _, terr := szkclient.conn.Exists(fpath1)
				if terr == nil && funcExists {
					tmpfuncs, _, terr := szkclient.conn.Get(fpath1)
					//plog.Debug("tmpfuncs is: ", tmpfuncs)
					if terr != nil {
						plog.Fatal("get funcs err: ", terr)
						err = terr
						szkclient.Mutx.Unlock()
						return //查询出错
					}
					var jsd_funcs []string
					jerr := json.Unmarshal(tmpfuncs, &jsd_funcs)
					if jerr != nil {
						plog.Fatal("funcs decode err: ", jerr)
						err = jerr
						szkclient.Mutx.Unlock()
						return //查询出错
					}
					//plog.Debug("the jsd funcs is: ", jsd_funcs)
					funcinfo[v] = jsd_funcs //待存储的方法部分
				} else {
					plog.Info("terr is: ", terr, "funcExists flag is: ", funcExists)
				}
			}
		}

		//zk查到之后，放到内存map：
		var zst zkServiceInfo
		zst.serviceName = svcName
		zst.vinfo = jsd_vinfo
		zst.funcinfo = funcinfo
		szkclient.ServiceAll[svcName] = &zst
		zkserviceinfoPtr = &zst
	}
	szkclient.Mutx.Unlock()
	v2ids = getProperFuncs(funcName, zkserviceinfoPtr)
	svrtype = szkclient.ServiceAll[svcName].vinfo.Typ

	//锁控制，防止一类服务的多个监听！
	szkclient.Mutx.Lock()
	//如果还没有监听，就开始监听
	if !(szkclient.haveWatched[svcName]) {
		fmt.Println("I am watcher!!!!!!!!!!!!!!!!!!!!!!!")
		go szkclient.watchNodeDataChange(spath)  //奇怪：注释掉该行代码后
		go szkclient.watchChildrenChanged(spath) //监控服务节点数据变化和子节点变化「增删」
		szkclient.haveWatched[svcName] = true    //添加监听标记
	}
	szkclient.Mutx.Unlock()

	return
}

//rpc启动添加节点的操作

//监听某个节点(此处主要是改值和子节点变化)
//func watchDemoNode(path string, conn *zk.Conn) {
//	//wg.Add(1)
//	//创建
//	//watchNodeCreated(path, conn)
//	//改值
//	go watchNodeDataChange(path, conn)
//	//子节点变化「增删」
//	go watchChildrenChanged(path, conn)
//	//删除节点
//	//watchNodeDeleted(path, conn)
//	//wg.Done()
//}

//func watchNodeCreated(path string, conn *zk.Conn) {
//	log.Println("watchNodeCreated")
//	for {
//		_, _, ch, _ := conn.ExistsW(path)
//		e := <-ch
//		log.Println("ExistsW:", e.Type, "Event:", e)
//		if e.Type == zk.EventNodeCreated {
//			log.Println("NodeCreated ")
//			return
//		}
//	}
//}
//func watchNodeDeleted(path string, conn *zk.Conn) {
//	log.Println("watchNodeDeleted")
//	for {
//		_, _, ch, _ := conn.ExistsW(path)
//		e := <-ch
//		log.Println("ExistsW:", e.Type, "Event:", e)
//		if e.Type == zk.EventNodeDeleted {
//			log.Println("NodeDeleted ")
//			return
//		}
//	}
//}

//服务启动时向zk注册！(重点：节点存在与否，版本信息如何，防并发修改！！)
func (szkclient *SzkClient) NodeRegister(servicename string, funcnames []string, svrtype int64, svrid int64, version int32) error {
	defer plog.CatchPanic()
	//plog.Debug("HHHHHHHHHHHHHHHHHHHHHH zk begin HHHHHHHHHHHHHHHHHHHHHHH")
	//数据初始化
	spath := zpathHead + "/" + servicename + "/typeandids"
	funcpath := zpathHead + "/" + servicename + "/funcs/" + strconv.Itoa(int(version)) //eg: "/rpc/arith/funcs/14"
	var typandidsObj zkSt = zkSt{
		Typ:     0,
		Ver:     0,
		Ids:     []int64{},
		PrevVer: 0,
		PrevIds: []int64{},
	}
	funcnamesbytes, err := json.Marshal(funcnames)
	if err != nil {
		//json编码异常
		return err
	}

	//连接
	err = szkclient.szkconnect(zkIpPorts, time.Second*2)
	if err != nil {
		plog.Fatal(err)
		return err
	}

	//带重试的设置zk！
	for i := 1; i <= registerRetryTimes; i++ {
		//初始化必须的目录（创建）
		szkclient.initPath(servicename)
		//funcs节点不存在(不存在就创建，存在就不处理！)并发创建，失败重试！
		funcNodeExist, _, err := szkclient.conn.Exists(funcpath)
		if err != nil {
			plog.Warn(funcpath, err)
			return err
		}
		if !funcNodeExist {
			_, err_create := szkclient.conn.Create(funcpath, funcnamesbytes, flags, acls)
			if err_create != nil {
				plog.Debug(err_create)
				//return err_create
				continue //可能是并发的创建导致的失败！retry
			}
		}
		//funcs节点修改完之后，开始修改typeandids节点
		svrNodeExist, _, err := szkclient.conn.Exists(spath)
		if err != nil {
			plog.Fatal(err)
			return err
		}
		//plog.Debug("HHHHHHHHHHHHHHHHHHHHHHendHHHHHHHHHHHHHHHHHHHHHHH")
		if !svrNodeExist {
			//节点不存在，创建节点
			szkclient.updateData(&typandidsObj, svrtype, svrid, version) //修改(添加)数据
			typandidsBytes, err := json.Marshal(typandidsObj)
			if err != nil {
				plog.Fatal(err)
				return err
			}

			_, err_create := szkclient.conn.Create(spath, typandidsBytes, flags, acls)
			if err_create != nil {
				plog.Fatal(err_create)
				//return err_create
				continue //可能是并发的创建导致的失败！retry
			} else {
				return nil //设置成功
			}
		} else {
			//节点存在，修改节点
			typandidsBytes, stat, err := szkclient.conn.Get(spath) //stat中存储了版本信息
			if err != nil {
				plog.Fatal("get zk data err: ", err)
				return err
			}
			json.Unmarshal(typandidsBytes, &typandidsObj)                //json反序列化到结构体
			szkclient.updateData(&typandidsObj, svrtype, svrid, version) //修改(添加)数据
			typandidsBytes, err = json.Marshal(typandidsObj)             //准备写入zk的数据
			if err != nil {
				plog.Fatal("marshal data err: ", err)
				return err
			}
			_, err = szkclient.conn.Set(spath, typandidsBytes, stat.Version) //json数据写入zk
			if err != nil {
				plog.Fatal("set zk data error :", err)
				//return err //写入失败，可能是版本号冲突了，或者其他情况
				continue //可能是并发的创建导致的失败！retry
			} else {
				return nil
			}
		}
	}

	//retry了N次之后还是失败的，那就记录日志！(重要点)
	//plog.Debug("retry " + strconv.Itoa(registerRetryTimes) + " times err!")
	return errors.New("retry " + strconv.Itoa(registerRetryTimes) + " times err!")
}

//初始化zk目录(确保各级目录已经创建好了！)
func (szkclient *SzkClient) initPath(servicename string) {
	path := zpathHead
	path1 := zpathHead + "/" + servicename
	path2 := zpathHead + "/" + servicename + "/funcs"
	//plog.Debug(path, path1, path2)
	//创建/rpc目录
	nodeExist, _, err := szkclient.conn.Exists(path)
	if !nodeExist && err == nil {
		_, err := szkclient.conn.Create(path, []byte("null"), flags, acls)
		if err != nil {
			plog.Warn(path, err)
		}
	}

	//创建/rpc/arith目录
	nodeExist, _, err = szkclient.conn.Exists(path1)
	if !nodeExist && err == nil {
		_, err := szkclient.conn.Create(path1, []byte("null"), flags, acls)
		if err != nil {
			plog.Warn(path1, err)
		}
	}

	//创建/rpc/arith/funcs目录
	nodeExist, _, err = szkclient.conn.Exists(path2)
	if !nodeExist && err == nil {
		_, err := szkclient.conn.Create(path2, []byte("null"), flags, acls)
		if err != nil {
			plog.Warn(path2, err)
		}
	}
}

//修改typeandids节点的结构体数据
func (szkclient *SzkClient) updateData(typandidsObj *zkSt, svrtype int64, svrid int64, version int32) {
	if typandidsObj.Typ == 0 {
		typandidsObj.Typ = svrtype
	}
	if version > typandidsObj.Ver { //大于最新版本
		typandidsObj.PrevIds = typandidsObj.Ids
		typandidsObj.PrevVer = typandidsObj.Ver
		typandidsObj.Ids = []int64{svrid}
		typandidsObj.Ver = version
	} else if version == typandidsObj.Ver { //等于最新版本
		//避免重复加入
		var svridExist = false
		for _, v := range typandidsObj.Ids {
			if svrid == v {
				svridExist = true
				break
			}
		}
		if !svridExist {
			typandidsObj.Ids = append(typandidsObj.Ids, svrid)
		}

	} else if version == typandidsObj.PrevVer { //等于上一个版本
		//避免重复加入
		var svridExist = false
		for _, v := range typandidsObj.PrevIds {
			if svrid == v {
				svridExist = true
				break
			}
		}
		if !svridExist {
			typandidsObj.PrevIds = append(typandidsObj.PrevIds, svrid)
		}

	} else {
		//比上一个版本还要小，不更新！(TODO：加日志)
	}
}

//监听节点的数据发生了变化！
func (szkclient *SzkClient) watchNodeDataChange(path string) {
	for {
		_, _, ch, _ := szkclient.conn.GetW(path)
		e := <-ch
		//dealEventHappened(e, path, conn)
		//数据变动通知内存调整
		plog.Info("GetW('"+path+"'):", e.Type, "Event:", e)
	}
}

//监听该节点内容有变动，或者子节点有增删！（EventNodeChildrenChanged/EventNodeDataChanged）
func (szkclient *SzkClient) watchChildrenChanged(path string) {
	for {
		_, _, ch, _ := szkclient.conn.ChildrenW(path)
		e := <-ch
		//plog.Debug("ChildrenW:", e.Type, "Event:", e)
		//数据变动通知内存调整
		szkclient.dealEventHappened(e, path)
	}
}

//处理事件发生
func (szkclient *SzkClient) dealEventHappened(e zk.Event, spath string) {
	svcName := filepath.Base(filepath.Dir(spath)) //获取服务名
	//funcHeadPath := zpathHead + "/" + svcName + "/funcs"
	//3表示EventNodeDataChanged 4表示EventNodeChildrenChanged 重置内存数据
	if e.Type == 3 || e.Type == 4 { //获取vinfo部分数据
		svcExists, _, terr := szkclient.conn.Exists(spath)
		if terr != nil || !svcExists {
			//如果节点不存在或者出错了直接返回！
			plog.Info("service not exists or err happened: ", terr)

			return
		}
		tmptypeandids, _, terr := szkclient.conn.Get(spath)
		//plog.Debug(string(tmptypeandids))
		if terr != nil {
			return
		}

		var jsd_vinfo *zkSt
		json.Unmarshal(tmptypeandids, &jsd_vinfo) //json解析到结构体中
		//plog.Debug(*jsd_vinfo)
		//funcinfo部分
		var funcinfo map[int32][]string = make(map[int32][]string)
		var versions = []int32{jsd_vinfo.Ver, jsd_vinfo.PrevVer}
		if jsd_vinfo.Ver == 0 && jsd_vinfo.PrevVer == 0 {
			return //zk中的版本都为空了，调用方就不更新了！
		}
		//plog.Debug("the versions are:", versions)
		for _, v := range versions { //每个版本的方法都存一次
			fpath1 := zpathHead + "/" + svcName + "/funcs/" + fmt.Sprintf("%d", v)
			//plog.Debug("the path is: ", fpath1)
			funcExists, _, terr := szkclient.conn.Exists(fpath1)
			//plog.Debug(v, " the exist flag is", funcExists)
			if terr == nil && funcExists {
				tmpfuncs, _, terr := szkclient.conn.Get(fpath1)
				if terr != nil {

					return //查询出错
				}
				//plog.Debug("the tmp funcs are:", tmpfuncs)
				var jsd_funcs []string
				json.Unmarshal(tmpfuncs, &jsd_funcs)
				funcinfo[v] = jsd_funcs //待存储的方法部分
			}
		}

		//zk查到之后，放到内存map：
		var zst zkServiceInfo
		zst.serviceName = svcName
		zst.vinfo = jsd_vinfo
		zst.funcinfo = funcinfo
		szkclient.ServiceAll[svcName] = &zst
		//zkserviceinfoPtr = &zst
		//fmt.Println(zst.serviceName)
		//fmt.Println(*zst.vinfo)
		//fmt.Println(zst.funcinfo)
	}

	fmt.Println(szkclient.ServiceAll[svcName])
}
