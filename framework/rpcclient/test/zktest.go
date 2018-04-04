package main

import (
	"fmt"
	"framework/rpcclient/szk"
	"time"
)

func main() {
	szkclient, err := szk.NewSzkClient()
	if err != nil {
		fmt.Println("new client err!")
	}
	svcName := "arith"
	funcName := "Add"
	typ, ids, err := szkclient.GetSerByNames(svcName, funcName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(typ, ids)

	funcName = "Multiply"
	typ, ids, err = szkclient.GetSerByNames(svcName, funcName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(typ, ids)
	time.Sleep(10000 * time.Second)
}
