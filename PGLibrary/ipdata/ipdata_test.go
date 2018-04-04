package ipdata

import (
	"testing"
	"fmt"
)
// 基准测试
func BenchmarkIpdata_Find(b *testing.B) {
	ipObj := NewIpdata()
	for i := 0;  i < b.N; i++ {
		ipInfo, _ := ipObj.Find("192.168.200.21")
		fmt.Println(ipInfo)
	}
}
// 功能测试
func TestIpdata_Find(t *testing.T) {
	ipObj := NewIpdata()
	ipInfo, _ := ipObj.Find("192.168.200.21")
	fmt.Println(ipInfo)
}
