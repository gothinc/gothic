package logger

import (
	"fmt"
	"github.com/json-iterator/go"
)

/**
 * @desc JsonFormatter
 * @author zhaojiangwei
 * @date 2018-12-14 17:03
 */

type JsonFormatter struct {
}

func NewJsonFormatter() *JsonFormatter{
	return &JsonFormatter{}
}

func (this *JsonFormatter) Format(v ...interface{}) (string, error)  {
	msg := fmt.Sprint(v)
	cont, err := jsoniter.MarshalToString(msg)
	if err != nil{
		return "", err
	}

	return cont, nil
}

func (this *JsonFormatter) FormatFields(fields EntryFields) (string, error)  {
	cont, err := jsoniter.MarshalToString(fields)
	if err != nil{
		return "", err
	}

	return cont, nil
}