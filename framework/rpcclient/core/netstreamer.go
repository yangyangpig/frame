package rpcclient

import (
	. "putil/perror"
)

type netstreamer interface {
	SetReconTimes(times int)
	Connect(lIp string, lPort int, ip string, port int) (err error)
	ReadAtLeast(data []byte, size int) (err Perror)
	Write(data []byte) (err Perror)
	Reconnect() (err error)
	GetRomteAddr() (ip string, port int, err error)
	GetLocalAddr() (ip string, port int, err error)
	Close()
}
