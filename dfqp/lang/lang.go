package lang

import (
	"dfqp/lang/zh"
	"dfqp/lang/tw"
)

var langMap map[string]map[int]string

const ZH = 0
const TW = 1

var types = map[int]string{
	0: "ZH",
	1: "TW",
}

func init() {
	langMap = make(map[string]map[int]string)
	langMap["ZH"] = zh.LangMap
	langMap["TW"] = tw.LangMap
}

// Msg get a message by code
func Msg(code int, langType int) string {
	langTypeString, ok := types[langType]
	if !ok {
		langTypeString, _ = types[0]
	}

	subLang, ok := langMap[langTypeString]
	if !ok {
		return ""
	}

	if msg, ok := subLang[code]; ok {
		return msg
	}
	return ""
}
