地址：http://svn.oa.com:24399/flames/go

目录规划如下：

/framework				rpc框架
/github.com				使用到的开源库
/PGLibrary				业务使用到的公共类库，业务自主封装
/PGNotices				公告服务
	/app
		/entity 		DB 表结构(结构体，方便存储数据使用)
		/proto			rpc请求参数和返回值(框架协议必备)
		/controller 	控制器（业务逻辑）
		/service 		model{db,cache}（数据操作）
	/conf  				配置，服务启动相关（也可以抽离到mysql）
/PGSignin				签到服务
	/app
		/entity 		DB 表结构(结构体，方便存储数据使用)
		/proto			rpc请求参数和返回值(框架协议必备)
		/controller 	控制器（业务逻辑）
		/service 		model{db,cache}（数据操作）
	/conf  				配置，服务启动相关（也可以抽离到mysql）

基本执行步骤:
1.linux环境安装go环境。
2.设置GOPATH变量。
3.在$GOPATH/src目录将此目录代码更新下来。
4.服务启动示例，比如公告服务：进入$GOPATH/src目录，直接执行go run PGNotices/main.go。

测试规范，以PGLibrary/ipdata为例：
1.每个功能都需要基准测试和功能测试，ipdata.go对应需要ipdata_test.go测试文件
2.功能测试（执行ipdata_test.go文件的TestIpdata_Find方法）：
    cd $GOPATH/src/PGLibrary/ipdata;
    go test
3.基准测试（执行ipdata_test.go文件的BenchmarkIpdata_Find方法）：
    cd $GOPATH/src/PGLibrary/ipdata;
    go test -bench=. -benchmem