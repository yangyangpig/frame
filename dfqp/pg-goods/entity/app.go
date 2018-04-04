package entity

// EnvDev 开发环境
const EnvDev = "dev"

// EnvTest 测试环境
const EnvTest = "test"

// EnvProduct 生产环境
const EnvProduct = "product"

// EnvValue app 环境类型
var EnvValue = map[string]int{
	"dev":     0,
	"preview": 1,
	"product": 2,
}