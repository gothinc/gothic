package gothic

import (
	"net/http"
	"strconv"
	"github.com/json-iterator/go"
)

/**
 * @desc Controller
 * @author zhaojiangwei
 * @date 2018-12-18 10:33
 */

type Controller struct {
	rw           http.ResponseWriter
	r            *http.Request
	outputBody   []byte
	OutputDirect bool //是否直接输出到http

	Context  *ThreadContext
}

func (this *Controller) Init(Context *ThreadContext, rw http.ResponseWriter, r *http.Request) {
	this.rw = rw
	this.r = r
	this.Context = Context
	this.OutputDirect = true
}

func (this *Controller) Destruct(){
	this.writeToWriter(this.outputBody)
}

func (this *Controller) OutputString(data string){
	this.rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	this.outputBody = []byte(data)
}

func (this *Controller) JsonSucc(data interface{}) {
	this.JsonFail(DefaultGothicResSucc, DefaultGothicResMsg, data)
}

func (this *Controller) JsonFail(code int, message string, data interface{}) {
	ret := GothicResponseBody{
		Errno: code,
		Errmsg: message,
		Data: data,
	}

	content, err := jsoniter.Marshal(ret)
	if err != nil {
		panic("server exception")
	}

	this.rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.outputBody = content
}

func (this *Controller) writeToWriter(rb []byte) {
	if this.OutputDirect {
		this.rw.Write(rb)
	}
}

func (this *Controller) GetString(key string, defaultValue string) string {
	ret := this.r.FormValue(key)
	if ret == "" {
		ret = defaultValue
	}
	return ret
}

func (this *Controller) GetStringSlice(key string) []string {
	if this.r.Form == nil {
		return []string{}
	}
	vs := this.r.Form[key]
	return vs
}

func (this *Controller) GetInt(key string, defaultValue int64) int64 {
	ret, err := strconv.ParseInt(this.r.FormValue(key), 10, 64)
	if err != nil {
		ret = defaultValue
	}
	return ret
}

func (this *Controller) GetBool(key string, defaultValue bool) bool {
	ret, err := strconv.ParseBool(this.r.FormValue(key))
	if err != nil {
		ret = defaultValue
	}
	return ret
}