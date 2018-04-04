package byteOrderOp

import (
	"unsafe"
	//"fmt"
)

const INT_SIZE int = int(unsafe.Sizeof(0))

//判断系统中的字节序类型
//返回：true：小尾， false：大尾
func SystemEdian() bool {
	var i int = 0x1
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		//fmt.Println("system edian is little endian")
		return true
	} else {
		//fmt.Println("system edian is big endian")
		return false
	}
}

//int数值网络字节序转换成主机字节序(net收到的字节切片的4字节转成int，并返回)
func Util_ntoh_int16(buf []byte) (value int16, err error) {
	_ = buf[1]
	value = int16(buf[0])<<8 | int16(buf[1])
	return
}

//主机字节序转换成网络字节序(将要发送的int数据转换成网络字节序)
func Util_hton_int16(len int16, buf []byte) {
	_ = buf[1]

	buf[0] = byte(len >> 8)
	buf[1] = byte(len)
}

//int数值网络字节序转换成主机字节序(net收到的字节切片的4字节转成int，并返回)
func Util_ntoh_int32(buf []byte) (value int32, err error) {
	_ = buf[3] //TODO?

	//	if len(buf) < 4 {
	//		err = errors.New("len is not enough")
	//		return
	//	}
	value = int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3])
	return
}

//主机字节序转换成网络字节序(将要发送的int数据转换成网络字节序)
func Util_hton_int32(len int32, buf []byte) {
	_ = buf[3] //TODO?

	buf[0] = byte(len >> 24)
	buf[1] = byte(len >> 16)
	buf[2] = byte(len >> 8)
	buf[3] = byte(len)
}

//int数值网络字节序转换成主机字节序
func Util_ntoh_int64(buf []byte) (value int64, err error) {
	_ = buf[7]
	value = int64(buf[0])<<56 | int64(buf[1])<<48 | int64(buf[2])<<40 | int64(buf[3])<<32 | int64(buf[4])<<24 | int64(buf[5])<<16 | int64(buf[6])<<8 | int64(buf[7])
	return
}

//主机字节序转换成网络字节序
func Util_hton_int64(len int64, buf []byte) {
	_ = buf[7]

	buf[0] = byte(len >> 56)
	buf[1] = byte(len >> 48)
	buf[2] = byte(len >> 40)
	buf[3] = byte(len >> 32)
	buf[4] = byte(len >> 24)
	buf[5] = byte(len >> 16)
	buf[6] = byte(len >> 8)
	buf[7] = byte(len)
}
