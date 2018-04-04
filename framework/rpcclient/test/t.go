package main

import (
	_ "errors"
	"fmt"
	_ "math/rand"
	"time"

	//"framework/rpcclient/test/tone"
	"github.com/astaxie/beego/config"
)

func main() {
	//	var ids map[int32][]int64 = make(map[int32][]int64)
	//	ids[4] = []int64{1, 2, 3, 4, 5, 6}
	//	ids[5] = []int64{12, 44, 55, 56, 34}
	//	fmt.Println(getRandomId(ids))
	fmt.Println(int(time.Microsecond))
	fmt.Println(time.Millisecond)
	fmt.Println(time.Second)

	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano())
	var Client int32 = 32
	fmt.Println("the main client is:", Client)
	//tttwo.Ech()

	//纳秒转微秒
	var diff int = 2000
	diff = diff / 1e3
	fmt.Println(diff)

	iniconf, _ := config.NewConfig("ini", "../conf/app.conf")
	fmt.Println(iniconf.Int("dev::retrytimes"))
	ipports := iniconf.Strings("dev::zkipport")
	fmt.Println(len(ipports))
	fmt.Println(ipports[1])
	fmt.Println(ipports[2])

}
