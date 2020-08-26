package uccontroller

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Bin-DuS/comm/core/ucerror"
	"github.com/Bin-DuS/comm/core/uclog"
	"github.com/Bin-DuS/comm/core/ucutil"
	"github.com/astaxie/beego"
)

type ResponseParams struct {
	Code           int                      `json:"code"`
	Msg            string                   `json:"msg"`
	RequestID     string                   `json:"request_id"`
	Timestamp      int64                    `json:"timestamp,omitempty"`
	Data           interface{}              `json:"data,omitempty"` // 通用响应
	SuccessData   []map[string]interface{} `json:"success_data"`   // 消息发送响应
	FailedData    []map[string]interface{} `json:"failed_data"`    // 消息发送响应
	ExternalField map[string]interface{}   `json:"-"`              // 除通用响应外增加的字段,在对应server的base中进行处理
}

func (response *ResponseParams) SetErrorRsp(code int, msg string) {
	response.Code = code
	response.Msg = msg
}

type Controller struct {
	beego.Controller
	uclog.UcLog
	MonitorLog *uclog.Monitorlog // 旧的统计日志
	SensorsLog *uclog.Monitorlog // 神策日志
	UserID      uint64
	ResID       uint64
	SessionID  string
	Rsp ResponseParams
}

// http_outputtype可以为.json(参见beego/mime.go，表示application/json)，或者"application/json; charset=utf-8"
func (controller *Controller) InitContentType() {
	if beego.AppConfig != nil {
		contentType := beego.AppConfig.DefaultString("http_outputtype", "")
		if contentType != "" {
			if strings.Contains(contentType, "/") {
				controller.Ctx.Output.Header("Content-Type", contentType)
			} else {
				controller.Ctx.Output.ContentType(contentType)
			}
		}
	}
}

func (controller *Controller) CheckClientSession() (err error) {
	controller.UserID = ucutil.ToUint64V2(controller.Ctx.Input.Header("User-Id"))
	if controller.UserID == 0 {
		return ucerror.ErrorUserIDIsZero
	}
	controller.SessionID = ucutil.ToString(controller.Ctx.Input.Header("Session-Id"))
	if controller.SessionID == "" {
		return ucerror.ErrorSessionIDIsEmpty
	}
	return nil
}

func (controller *Controller) GetUrlParam(key string) string {
	for k, v := range controller.Input() {
		if strings.ToLower(k) == strings.ToLower(key) {
			if len(v) == 0 {
				return ""
			}
			return v[0]
		}
	}
	return ""
}

func (controller *Controller) parseParams(v *reflect.Value, src string) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tempV, err := strconv.ParseInt(src, 10, 0); err == nil {
			v.SetInt(tempV)
			return true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if tempV, err := strconv.ParseUint(src, 10, 0); err == nil {
			v.SetUint(tempV)
			return true
		}

	case reflect.Float32, reflect.Float64:
		if tempV, err := strconv.ParseFloat(src, 10); err == nil {
			v.SetFloat(tempV)
			return true
		}

	case reflect.Bool:
		t := strings.ToLower(src)
		if t == "true" {
			v.SetBool(true)
		} else {
			v.SetBool(false)
		}
		return true

	case reflect.String:
		v.SetString(src)
		return true

	default:
		return false
	}
	return false
}

func (controller *Controller) ParseRequestParams(dest interface{}) bool {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return false
	}

	v := reflect.ValueOf(dest).Elem()
	if v.Kind() != reflect.Struct {
		return false
	}

	realSize := uint32(v.NumField())
	for i := 0; i < int(realSize); i++ {
		tag := v.Type().Field(i).Tag
		key := tag.Get("param")

		if key == "" {
			continue
		}

		paraStr := controller.GetUrlParam(key)
		if key == "" {
			continue
		}

		item := v.Field(i)
		if !controller.parseParams(&item, paraStr) {
			return false
		}
	}
	return true
}

func (controller *Controller) DoGetConfiguration() {
	data := map[string]interface{}{
		"log_level": uclog.GetLogLevelDesc(),
	}
	controller.Rsp.Data = data
}

func (controller *Controller) DoSetConfiguration() {
	loglevel := controller.GetUrlParam("log_level")
	if len(loglevel) == 0 {
		controller.Rsp.Msg = "invalid log level"
		return
	}
	uclog.SetLogLevel(loglevel)
	data := map[string]interface{}{
		"log_level": uclog.GetLogLevelDesc(),
	}
	controller.Rsp.Data = data
}
