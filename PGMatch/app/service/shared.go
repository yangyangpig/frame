package service

import (
	"strings"
	"math"
	"time"
	"runtime/debug"
	"fmt"
	"sort"
	"math/rand"
)

// 大厅版本转换为 version
func HallVerToLang(version string) int {
	hallVerSlice := strings.Split(version, ".")
	i := len(hallVerSlice) - 1
	ver := hallVerSlice[i]
	return StrToInt(ver)
}

// 判断客户端是否为 iOS
func ClientIsIos(appId uint64) bool {
	clientId := getClientId(appId)
	return clientId == 1 || clientId == 2
}

// 判断客户端是否为 Android
func ClientIsAndroid(appId uint64) bool {
	clientId := getClientId(appId)
	return clientId == 1 || clientId == 2
}

// 获取客户端 id
func getClientId(appId uint64) int {
	return int(math.Floor(float64((appId % 100000) / 1000)))
}

// 格式化 2006-01-02 15:04:05 的时间字符串, 返回 Time 类型
func FormatTime(timeStr string) time.Time {
	str := strings.Trim(timeStr, " ")
	hasDate := strings.Contains(timeStr, "-")
	hasTime := strings.Contains(timeStr, ":")
	formatStr := ""
	if hasDate && hasTime {
		formatStr = "2006-01-0215:04:05"
	} else if hasDate {
		formatStr = "2006-01-02"
	} else if hasTime {
		formatStr = "15:04:05"
	}
	if formatStr == "" {
		fmt.Printf("Time Format : Formatting rules do not exist! value : %v \n", timeStr)
		debug.PrintStack()
		return time.Time{}
	}
	retTime, err := time.ParseInLocation(formatStr, str, time.Local)
	if err != nil {
		fmt.Printf("time format is fail!\n")
		debug.PrintStack()
		return time.Time{}
	}
	return retTime
}

type StringSlice []string

// 获取长度
func (c StringSlice) Len() int {
	return len(c)
}

// 如果 index 为 i 的元素小于 index 为 j 的元素，则返回 true，否则返回 false
func (c StringSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Swap 交换索引为 i 和 j 的元素
func (c StringSlice) Less(i, j int) bool {
	return c[i] < c[j]
}

// 返回映射的 key 值排序切片
func KeySortedSlice(m map[string]interface{}) StringSlice {
	// 创建 key 切片
	keySlice := make(StringSlice, len(m))
	for k := range m {
		keySlice = append(keySlice, k)
	}
	if !sort.IsSorted(keySlice) {
		sort.Sort(keySlice)
	}
	return keySlice
}

// 随机一个 int 范围在 min - max
func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}