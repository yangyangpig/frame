package main

import (
	"fmt"
	"putil/log"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	//  Socket to talk to server
	//	fmt.Println("Connecting to hello world server...")
	//	requester, _ := zmq.NewSocket(zmq.REQ)
	//	defer requester.Close()
	//	err :=requester.Connect("tcp://192.168.96.156:9897")
	//fmt.Println("the connect result is: ",err)

	//	for request_nbr := 0; request_nbr != 10; request_nbr++ {
	//		// send hello
	//		msg := fmt.Sprintf("Hello %d", request_nbr)
	//		fmt.Println("Sending ", msg)
	//		requester.Send(msg, 0)

	//		// Wait for reply:
	//		reply, _ := requester.Recv(0)
	//		fmt.Println("Received ", reply)
	//	}

	bind_to := "tcp://192.168.96.156:9897"
	s_out, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		plog.Fatal("NewSocket 2:", err)
	}

	s_out.SetSndhwm(10000) //设置发送最大队列
	s_out.SetRcvhwm(10000) //设置接收最大队列
	s_out.SetConflate(false)
	err = s_out.Connect(bind_to)
	if err != nil {
		plog.Fatal("s_out.Connect:", err)
	}

	message_count := 20000

	for j := 0; j < message_count; j++ {
		_, err = s_out.Send("R"+fmt.Sprint(j), zmq.DONTWAIT) //zmq.DONTWAIT 表示非阻塞！
		if err != nil {
			plog.Fatal("s_out.Send %d: %v", j, err)
		}
		//time.Sleep(50 * time.Microsecond)
	}

	fmt.Println("send finished!")

	time.Sleep(time.Second * 50000)
}
