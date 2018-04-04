package netutil

import (
	"log"
	"strconv"
	"strings"
)

//addr 格式：ip:port
func GetNetAddr(ip string, port int) (addr string) {
	//ad := make([]string, 2)
	var ad [2]string
	ad[0] = ip
	ad[1] = strconv.Itoa(port)
	addr = strings.Join(ad[0:2], ":")
	return
}

//addr 格式 ip:port
func GetIPandPortByAddr(addr string) (ip string, port int, err error) {
	ipandport := strings.Split(addr, ":")
	ip = ipandport[0]
	port, err = strconv.Atoi(ipandport[1])
	if err != nil {
		log.Fatal("addr is a invlid parameter")
	}
	return

}
