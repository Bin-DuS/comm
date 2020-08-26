package uccontroller

import (
	"bytes"
	"context"
	"encoding/json"
	beegocontext "github.com/astaxie/beego/context"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

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
type fakeResponseWriter struct{}

func (f *fakeResponseWriter) Header() http.Header {
	return http.Header{}
}
func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}
func (f *fakeResponseWriter) WriteHeader(n int) {}
func (ctr *Controller) TestData(body string){
	ctr.Ctx = &beegocontext.Context{
		Request: &http.Request{URL: &url.URL{Scheme: "http", Host: "localhost", Path: "/"}},
		ResponseWriter:&beegocontext.Response{&fakeResponseWriter{},true,http.StatusOK,time.Nanosecond},
	}
	ctr.Ctx.Output = &beegocontext.BeegoOutput{Context: ctr.Ctx}
	ctr.Ctx.Input = &beegocontext.BeegoInput{Context: &beegocontext.Context{Request: ctr.Ctx.Request}}
	ctr.Ctx.Request.Header = http.Header{}
	ctr.Ctx.Request.Header.Set("Content-Type","application/json")
	ctr.Ctx.Input.RequestBody = []byte(body)

}
func (ctr *Controller) ParseRequest(ctx context.Context, v interface{}) (err error) {
	log, _ := uclog.ExtractLogger(ctx)
	contentType := ctr.Ctx.Input.Header("Content-Type")
	body := ctr.Ctx.Input.RequestBody
	query := ctr.Ctx.Request.URL.RawQuery
	// 压缩请求参数中空白字符
	var param string
	dst := new(bytes.Buffer)
	err = json.Compact(dst, body)
	if err != nil {
		param = strings.Replace(string(body), "\n", "", -1)
	} else {
		param = dst.String()
	}
	log.Log_Info("param:%s, query:%s", param, query)
	if len(body) > 0 {
		err = json.Unmarshal(body, v)
		if err != nil {
			// 有可能不是json格式的数据, 所以不能直接返回
			if contentType == "application/json" {
				return
			} else {
				values := ctr.Input()
				if len(values) > 0 {
					err = ctr.ParseForm(v)
				}
			}
		}
	}
	if len(query) > 0 {
		err = ctr.ParseForm(v)
	}
	return
}
// http_outputtype可以为.json(参见beego/mime.go，表示application/json)，或者"application/json; charset=utf-8"
func (ctr *Controller) InitContentType() {
	if beego.AppConfig != nil {
		contentType := beego.AppConfig.DefaultString("http_outputtype", "")
		if contentType != "" {
			if strings.Contains(contentType, "/") {
				ctr.Ctx.Output.Header("Content-Type", contentType)
			} else {
				ctr.Ctx.Output.ContentType(contentType)
			}
		}
	}
}

func (ctr *Controller) CheckClientSession() (err error) {
	ctr.UserID = ucutil.ToUint64V2(ctr.Ctx.Input.Header("User-Id"))
	if ctr.UserID == 0 {
		return ucerror.ErrorUserIDIsZero
	}
	ctr.SessionID = ucutil.ToString(ctr.Ctx.Input.Header("Session-Id"))
	if ctr.SessionID == "" {
		return ucerror.ErrorSessionIDIsEmpty
	}
	return nil
}

func (ctr *Controller) GetUrlParam(key string) string {
	for k, v := range ctr.Input() {
		if strings.ToLower(k) == strings.ToLower(key) {
			if len(v) == 0 {
				return ""
			}
			return v[0]
		}
	}
	return ""
}

func (ctr *Controller) parseParams(v *reflect.Value, src string) bool {
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

func (ctr *Controller) ParseRequestParams(dest interface{}) bool {
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

		paraStr := ctr.GetUrlParam(key)
		if key == "" {
			continue
		}

		item := v.Field(i)
		if !ctr.parseParams(&item, paraStr) {
			return false
		}
	}
	return true
}

func (ctr *Controller) DoGetConfiguration() {
	data := map[string]interface{}{
		"log_level": uclog.GetLogLevelDesc(),
	}
	ctr.Rsp.Data = data
}

func (ctr *Controller) DoSetConfiguration() {
	loglevel := ctr.GetUrlParam("log_level")
	if len(loglevel) == 0 {
		ctr.Rsp.Msg = "invalid log level"
		return
	}
	uclog.SetLogLevel(loglevel)
	data := map[string]interface{}{
		"log_level": uclog.GetLogLevelDesc(),
	}
	ctr.Rsp.Data = data
}
