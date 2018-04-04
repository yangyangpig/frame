package service

import (
	"strconv"
	"log"
	"bytes"
	"encoding/binary"
	"fmt"
)

// int 转 string
// 源码
// Itoa is shorthand for FormatInt(i, 10).
//func Itoa(i int) string {
//	return FormatInt(int64(i), 10)
//}
func IntToStr(i int) string {
	return strconv.Itoa(i)
}

// string -> int
func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("string -> int faile")
	}
	return i
}

// int64 -> string
func Int64ToStr(i64 int64) string {
	return strconv.FormatInt(i64, 10) // 10 表示 10 进制
}

// string -> int64
func StrToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("string -> int64 faile")
	}
	return i
}

// uint32 -> string
func Uint32ToStr(ui32 uint32) string {
	return fmt.Sprint(ui32)
}

// string -> uint32
func StrToUint32(s string) uint32 {
	return uint32(StrToUint64(s))
}

// uint64 -> string
func Uint64ToStr(ui64 uint64) string {
	return strconv.FormatUint(ui64, 10)
}

// string -> uint64
func StrToUint64(s string) uint64 {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalf("string -> uint64 faile")
	}
	return u
}

// []byte -> int32
func ByteToInt32(b []byte) (i32 int32) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.BigEndian, &i32)
	return
}

// []byte -> uint64
func ByteToUint64(b []byte) (u64 uint64) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.BigEndian, &u64)
	return
}

// int32 -> []byte
func Int32ToByte(i32 int32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, i32)
	return bytesBuffer.Bytes()
}

// []byte -> string
func ByteToStr(b []byte) string {
	return string(b)
}

// string -> []byte
func StrToByte(s string) []byte {
	return []byte(s)
}

// string -> float64
func StrToFloat64(s string) float64 {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("string -> uint64 faile")
	}
	return f64
}

// float64 -> string
func Float64ToStr(f64 float64) string {
	return strconv.FormatFloat(f64, 'E', -1, 64)
}
