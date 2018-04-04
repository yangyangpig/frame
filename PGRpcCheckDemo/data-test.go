package main

import (
	"PGRpcCheckDemo/app/controllers"
	"putil/log"
)

func main() {
	bigPack := new(controllers.BigPack)
	bigPack.Route = 5000
	bigPack.Container = make(map[string]string)
	container := bigPack.CreateBigPack()
	plog.Debug(container)
}
