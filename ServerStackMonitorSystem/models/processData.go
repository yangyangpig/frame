package models

import (
	_ "fmt"
	"reflect"
	"strings"
)

type ProcessData struct {
	Funcname   string
	Servername string
	Param      map[string]string
}

//[]reflect.Value
func Process(p *ProcessData) (res []byte) {
	t := reflect.ValueOf(new(ProcessPb))
	name := strFirstToUpper(p.Servername) + strFirstToUpper(p.Funcname)

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(p.Param)
	s := t.MethodByName(name).Call(params)
	return s[0].Interface().([]byte)
}

func ProcessRes(req []byte, serv string, funame string) string {
	t := reflect.ValueOf(new(ProcessPb))
	name := strFirstToUpper(serv) + strFirstToUpper(funame) + "Repose"
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(req)
	s := t.MethodByName(name).Call(params)
	return s[0].Interface().(string)
}

func strFirstToUpper(str string) (res string) {
	rs := []rune(str)
	len := len(rs)
	first := strings.ToUpper(string(rs[:1]))
	other := string(rs[1:len])
	res = first + other
	return

}
