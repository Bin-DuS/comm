package ucclienttype

//ProductTypeStr 产品类型BI映射
var ProductTypeStr = map[string]string{
	"beeteam": "beeteam",
	"tb":      "timebook",
	"uc":      "bee",
	"meeting": "meeting",
}

//ProductTypeToStr 转换映射关系
func ProductTypeToStr(productType string) string {
	val, ok := ProductTypeStr[productType]
	if ok {
		return val
	}
	return ""
}
