package PGLibrary

import (
	"regexp"
	"strconv"
	"fmt"
	"strings"
	"math"
	//"bytes"
	//"compress/gzip"
	//"io/ioutil"
	"encoding/base64"
	//"log"
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1>255 || ip2>255 || ip3>255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}

func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}

/**
 * 将字符串类型的版本号转化为整型数字
 * @param string $value 版本号字符串(2.0.1)
 * @return int 版本号的整型数字
 */
func Ver2long(ver string) int {
	verArr := strings.Split(ver, ".")
	defer func() {
		if err := recover(); err != nil {
			//todo
		}
	}()
	_ = verArr[2]
	major, _ := strconv.Atoi(verArr[0])
	minor, _ := strconv.Atoi(verArr[1])
	release, _ := strconv.Atoi(verArr[2])
	return major * 1000000 + minor * 1000 + release
}

/**
 * 获取应用的地区ID
 * @param int $appid 应用ID
 * @return int 地区ID
 */
func GetRegionId(appid int32) int32 {
	return int32(math.Floor(float64(appid / 100000)))
}

/**
 * base64转义
 * @param string str 需要处理的字符串
 * @return string 处理号之后的字符串
 */
func Base64Encode(byteContent []byte) []byte {
	tmpStr := base64.StdEncoding.EncodeToString(byteContent)
	return []byte(tmpStr)
}

/**
 * base64反转义
 * @param string str 需要处理的字符串
 * @return byte 处理号之后的字符串
 */
func Base64Decode(str string) []byte {
	decodeByte, _ := base64.StdEncoding.DecodeString(str)
	return decodeByte
}

/**
 * gzip 压缩
 * @param string str 需要处理的字符串
 * @return string 处理号之后的字符串
 */
func DoGzip(byteContent []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(byteContent)
	w.Flush()
	return b.Bytes()
}
/**
 * gzip 解压
 * @param string str 需要处理的字符串
 * @return byte 处理号之后的字符串
 */
func DoUnGzip(str string) []byte {
	in := bytes.NewReader([]byte(str))
	r2, _ := gzip.NewReader(in)
	defer r2.Close()
	unByte, _ := ioutil.ReadAll(r2)
	return unByte
}