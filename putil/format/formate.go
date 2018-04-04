package format

import "fmt"
import "strings"


func ByteSlice2HexString(DecimalSlice []byte) string {
    var sa = make([]string, 0)
    for index, v := range DecimalSlice {
		if index == len(DecimalSlice)- 1{
			sa = append(sa, fmt.Sprintf("%#0x\n", v))	
		}else{
			sa = append(sa, fmt.Sprintf("%#0x,", v))	
		}
        
    }
    ss := strings.Join(sa, "")
    return ss
}