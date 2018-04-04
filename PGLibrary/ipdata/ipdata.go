package ipdata

import (
	"bytes"
	"encoding/binary"
	"bufio"
	"strings"
	"strconv"
	"os"
	"PGLibrary"
)

var instance map[string]map[string]string

type Ipdata struct {}

func (this *Ipdata) Find(ipStr string) (map[string]string, error) {
	if v, ok := instance[ipStr]; ok {
		return v, nil
	}
	var (
		indexOffset int32
		indexLength int16
		length int32
		len2 int32
		err error
		fp *os.File
		buf []byte
		buffer *bytes.Buffer
	)
	fp, err = os.Open("./ipv4data.datx")
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	r := bufio.NewReader(fp)
	buf = make([]byte, 4)
	_, err = fp.Read(buf)
	if err != nil {
		return nil, err
	}
	buffer = bytes.NewBuffer(buf)
	binary.Read(buffer, binary.BigEndian, &length)

	buf = make([]byte, length, length)
	r.Read(buf)

	arr := strings.Split(ipStr, ".")
	ip0, _:= strconv.Atoi(arr[0])
	ip1, _:= strconv.Atoi(arr[1])
	offset := (ip0*256 + ip1)*4

	buffer = bytes.NewBuffer(buf[offset:offset+4])
	binary.Read(buffer, binary.LittleEndian, &len2)
	maxCompLen := length - 262144 - 4

	buffer = bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, PGLibrary.Ip2long(ipStr))
	for start := len2 * 9 + 262144; start < maxCompLen; start += 9 {
		if bytes.Compare(buf[start:start+4], buffer.Bytes()) >= 0 {
			tmpByte := []byte{};
			tmpByte = append(tmpByte, buf[start+4:start+7]...)
			tmpByte = append(tmpByte, "\x00"...)
			binary.Read(bytes.NewBuffer(tmpByte), binary.LittleEndian, &indexOffset)
			binary.Read(bytes.NewBuffer(buf[start+7:start+9]), binary.BigEndian, &indexLength)
			break
		}
	}
	result := make(map[string]string)
	if indexOffset == 0 {
		return result, nil
	}
	fp.Seek(int64(length+indexOffset)-262144, 0)
	buf = make([]byte, indexLength)
	r.Read(buf)
	strArr := strings.Split(string(buf), "\t")
	result["country"] = strArr[0]
	result["province"] = strArr[1]
	result["city"] = strArr[2]
	result["operator"] = strArr[3]
	instance = make(map[string]map[string]string)
	instance[ipStr] = result
	return result, nil
}

func NewIpdata() *Ipdata  {
	return &Ipdata{}
}
