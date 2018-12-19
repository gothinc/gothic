package gothic

import (
	"net/http"
	"time"
	"encoding/json"
	"strconv"
)

/**
 * @desc Controller
 * @author zhaojiangwei
 * @date 2018-12-18 10:33
 */

type Controller struct {
	rw           http.ResponseWriter
	r            *http.Request
	startTime    time.Time
	OutputDirect bool //是否直接输出到http
}

func (this *Controller) Init(rw http.ResponseWriter, r *http.Request) {
	this.startTime = time.Now()
	this.rw = rw
	this.r = r
	this.OutputDirect = true
}

func (this *Controller) Destroy(){

}

func (this *Controller) JsonSucc(data interface{}) {
	errInfo := ""
	this.JsonFail(errInfo, data)
}

func (this *Controller) JsonFail(errInfo string, data interface{}) {
	ret := map[string]interface{}{
		"Errno":     0,
		"Errmsg":    "ok",
		"Data":      data,
	}

	content, err := json.MarshalIndent(ret, "", "  ")
	if err != nil {
		panic("server exception")
	}

	this.rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.writeToWriter(content)
}

func (this *Controller) writeToWriter(rb []byte) {
	//this.rw.Header().Set("Content-Length", strconv.Itoa(len(rb)))
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

func (this *Controller) GetStrings(key string) []string {
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