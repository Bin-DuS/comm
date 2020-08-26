package ucerror

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/config"
)

const (
	ErrUsualBegin = 10000 +iota
	ErrUserIDIsZero
	ErrSessionIDIsEmpty
	ErrUsualEnd
)
var (
	ErrorUsualBegin = Errorf(ErrUsualBegin,"Usual error begin")
	ErrorUserIDIsZero = Errorf(ErrUserIDIsZero,"UserID is 0")
	ErrorSessionIDIsEmpty = Errorf(ErrSessionIDIsEmpty,"SessionID Is Empty")
	ErrorUsualEnd = Errorf(ErrUsualEnd,"Usual error end")
)

type errorCodeAndString struct {
	code int
	text string
	data interface{}
}

func (e *errorCodeAndString) Error() string {
	return e.text
}

func Error(code int) error {
	return Errorf(code, Desc(code))
}

func Errorf(code int, format string, a ...interface{}) error {
	return &errorCodeAndString{code: code, text: fmt.Sprintf(format, a...)}
}

func ErrorCode(err error) int {
	if err == nil {
		return S_Ok
	}
	if e, ok := err.(*errorCodeAndString); ok {
		return e.code
	}

	return S_Error
}

func ErrorCodeWithDefault(err error, defaultCode int) int {
	if err == nil {
		return defaultCode
	}
	if e, ok := err.(*errorCodeAndString); ok {
		return e.code
	}

	return defaultCode
}

func ErrorData(err error) interface{} {
	if err == nil {
		return nil
	}
	if e, ok := err.(*errorCodeAndString); ok {
		return e.data
	}

	return nil
}

func ErrorDataWithDefault(err error, defaultValue interface{}) interface{} {
	if err == nil {
		return defaultValue
	}
	if e, ok := err.(*errorCodeAndString); ok && e.data != nil {
		return e.data
	}

	return defaultValue
}

func SetErrorData(err error, d interface{}) {
	if e, ok := err.(*errorCodeAndString); ok {
		e.data = d
	}
}

func IsError(code int) bool {
	return code != S_Ok
}

func Desc(code int, err ...interface{}) string {
	// 返回错误信息，如果没有则返回err.Error()
	var msg string
	if len(err) > 0 {
		switch t := err[0].(type) {
		case error:
			msg = t.Error()
		case string:
			msg = t
		}
	}
	if len(msg) <= 0 {
		if v, ok := errMap[code]; ok {
			msg = v
		}
	}
	return msg
}

var DefaultModuleBits uint = 24
var DefaultSection = "errorcode"

func NewCode(mode int, code int) int {
	if mode >= 0 && code >= 0 {
		return mode<<DefaultModuleBits | code
	} else {
		return code
	}
}

// 从配置文件中加载错误码信息,配置文件需要使用section,并且全部小写,默认为errorcode
func Load(filename string, section ...string) {
	iniConfig, err := config.NewConfig("ini", filename)
	if err != nil {
		return
	}
	var sect string = DefaultSection
	if len(section) > 0 {
		sect = section[0]
	}
	confMap, err := iniConfig.GetSection(sect)
	if err != nil {
		return
	}
	for key, value := range confMap {
		code, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		errMap[code] = value
	}
}
