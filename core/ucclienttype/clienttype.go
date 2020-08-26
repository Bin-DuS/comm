package ucclienttype

const (
	iphone  int64 = 1
	android int64 = 2
	pc      int64 = 3
	backend int64 = 4
	ipad    int64 = 5
	web     int64 = 6
	mac     int64 = 7
)

var ClientTypeStr = map[int64]string{
	iphone:  "IPHONE",
	android: "ANDROID",
	pc:      "PC",
	backend: "BACKEND",
	ipad:    "IPAD",
	web:     "WEB",
	mac:     "MAC",
}

/*
 * 神策埋点客户端名称
 */
var ShenceClientTypeStr = map[int64]string{
	iphone:  "ios",
	android: "android",
	pc:      "windows",
	mac:     "macos",
}

func GetShenceOsType(clientType int64) string {
	if _, ok := ShenceClientTypeStr[clientType]; !ok {
		return ShenceClientTypeStr[pc]
	}
	return ShenceClientTypeStr[clientType]
}

func IsAllowed(clientType int64) bool {
	if iphone <= clientType && clientType <= mac {
		return true
	}
	return false
}

func IsIOS(clientType int64) bool {
	if iphone == clientType {
		return true
	}
	return false
}

func IsPC(clientType int64) bool {
	if pc == clientType {
		return true
	}
	return false
}

func IsAndroid(clientType int64) bool {
	if android == clientType {
		return true
	}
	return false
}

func IsMobile(clientType int64) bool {
	return IsIOS(clientType) || IsAndroid(clientType)
}

func IsMac(clientType int64) bool {
	if clientType == mac {
		return true
	}
	return false
}

func ToStr(clientType int64) string {
	val, ok := ClientTypeStr[clientType]
	if ok {
		return val
	}
	return ""
}
