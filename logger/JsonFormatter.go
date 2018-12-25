package logger

import (
	"github.com/json-iterator/go"
	"time"
)

/**
 * @desc JsonFormatter
 * @author zhaojiangwei
 * @date 2018-12-14 17:03
 */

type JsonFormatter struct {
	TimestampFormat string
}

func NewJsonFormatter(timestampFormat string) *JsonFormatter{
	return &JsonFormatter{
		TimestampFormat: timestampFormat,
	}
}

func (this *JsonFormatter) SetTimestampFormat(format string){
	this.TimestampFormat = format
}

func (this *JsonFormatter) Format(v ...interface{}) (string, error)  {
	ret := make(map[string]interface{}, 2)
	if this.TimestampFormat != ""{
		ret["time"] = time.Now().Format(this.TimestampFormat)
	}

	if len(v) == 1{
		ret["msg"] = v[0]
	}else{
		ret["msg"] = v
	}

	cont, err := jsoniter.MarshalToString(ret)
	if err != nil{
		return "", err
	}

	return cont, nil
}

func (this *JsonFormatter) FormatFields(fields EntryFieldsAny) (cont string, err error)  {
	switch entity := fields.(type) {
	case EntryFields:
		if this.TimestampFormat != "" {
			entity["time"] = time.Now().Format(this.TimestampFormat)
		}
		cont, err = jsoniter.MarshalToString(entity)
	default:
		cont, err = jsoniter.MarshalToString(fields)
	}

	return
}