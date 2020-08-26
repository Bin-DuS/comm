package ucutil

import (
	"encoding/json"
	"fmt"
	"github.com/Bin-DuS/comm/core/uclog"
	"reflect"
	"sort"
	"strings"
)

func GetInt64ArrayFromMap(data map[string]interface{}, key string) ([]int64, bool) {
	v, ok := data[key]
	if !ok {
		return nil, false
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) <= 0 {
		return nil, false
	}
	result := make([]int64, len(arr))
	for i, val := range arr {
		result[i] = ToInt64(val, 0)
	}
	return result, true
}

func GetUint64ArrayFromMap(data map[string]interface{}, key string) ([]uint64, bool) {
	v, ok := data[key]
	if !ok {
		return nil, false
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) <= 0 {
		return nil, false
	}
	result := make([]uint64, len(arr))
	for i, val := range arr {
		result[i] = ParseUint(val, 0)
	}
	return result, true
}

func GetInt64FromMap(data map[string]interface{}, key string) int64 {
	v, exist := data[key]
	if !exist {
		return 0
	}
	return ToInt64(v, 0)
}

func GetUint64FromMap(data map[string]interface{}, key string) uint64 {
	v, exist := data[key]
	if !exist {
		return 0
	}
	result, _ := ToUint64(v)
	return result
}

func GetInt8FromMap(data map[string]interface{}, key string) int8 {
	v, exist := data[key]
	if !exist {
		return 0
	}
	return int8(ToInt64(v, 0))
}

func GetStringFromMap(data map[string]interface{}, key string) string {
	v, exist := data[key]
	if !exist {
		return ""
	}
	return ToString(v)
}

func GetMapFromMap(data map[string]interface{}, key string) (map[string]interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid param of data nil")
	}
	v, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("not found %s in data: %v", key, data)
	}
	value, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid %s type: %v", key, v)
	}
	return value, nil
}

func GetSliceFromMap(data map[string]interface{}, key string) ([]interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid param of data nil")
	}
	v, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("not found %s in data: %v", key, data)
	}
	value, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid %s type: %v", key, v)
	}
	return value, nil
}

func StructToMap(data interface{}, args ...interface{}) (result map[string]interface{}) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	result = make(map[string]interface{}, t.NumField())

	// take the json struct tag as the default map field name
	tagName := "json"
	if len(args) > 0 {
		if customTag, ok := args[0].(string); ok && customTag != "" {
			tagName = customTag
		}
	}

	for i := 0; i < t.NumField(); i++ {
		mFieldName := t.Field(i).Tag.Get(tagName)
		if mFieldName == "" {
			mFieldName = strings.ToLower(t.Field(i).Name)
		}
		result[mFieldName] = v.Field(i).Interface()
	}

	return
}

func SortStrMapKeys(m map[string]interface{}) []string {

	keys := []string{}

	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

/*
函数功能：interface转为map
函数入参：
v：待转化数据
函数返回值：转化后的map
*/
func InterfaceToMap(v interface{}) map[string]interface{} {
	re := make(map[string]interface{})
	vt := reflect.ValueOf(v)
	if vt.Kind() == reflect.Map {
		for _, key := range vt.MapKeys() {
			value := vt.MapIndex(key)
			re[ToString(key.Interface())] = value.Interface()
		}
	}
	return re
}

func StringToMap(v string) map[string]interface{} {
	re := make(map[string]interface{})
	vt := reflect.ValueOf(v)
	if vt.Kind() == reflect.String {
		err := json.Unmarshal([]byte(v), &re)
		if err != nil {
			uclog.Error("Unmarshal with error: %+v\n", err)
			return nil
		}
	}
	return re
}

//判断value是否在slice中
func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
