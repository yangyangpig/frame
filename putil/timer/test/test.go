package main

import (
	"fmt"
)

import(
	cc "putil/timer"
	"time"
)
type A struct{
	
}


var	timer cc.Timer

	


func init(){
	timer = cc.NewTimer()
}
func (a *A)timenotify(body interface{}){
	t, e := body.([]byte)
	if !e {
		return
	}
	fmt.Println(string(t))
	timer.AddEvent(1, 1000*time.Millisecond, a.timenotify, []byte("hello world"))
}
func main(){
	a := new(A)
	
	timer.AddEvent(1, 1000*time.Millisecond, a.timenotify, []byte("hello world"))
	//timer.Close(1)
	
	
	var controlstr string
	fmt.Scanln(&controlstr)

	if controlstr == "q" || controlstr == "Q" {
		return
	}
	timer.Close(1)
}