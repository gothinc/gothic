package gothic

import (
	"github.com/gothinc/gothic/logger"
	"github.com/gothinc/gothic/httpclient"
)

/**
 * @desc Export
 * @author zhaojiangwei
 * @date 2018-12-13 15:56
 */

var Logger = logger.NewDefaultLogger()

type EntryFields = logger.EntryFields
type EntryFieldsAny = logger.EntryFieldsAny

func Format(fields EntryFieldsAny) *logger.Entry{
	return Logger.Format(fields)
}

func Debug(v ...interface{}) {
	Logger.Debug(v...)
}

func Access(v ...interface{}) {
	Logger.Access(v...)
}

func Warn(v ...interface{}) {
	Logger.Warn(v...)
}

func Error(v ...interface{}) {
	Logger.Error(v...)
}

func Extend(logName string, v ...interface{}) {
	Logger.Extend(logName, v...)
}

func AddDefinedVariable(key string, value interface{}){
	Application.definedVariables[key] = value
}

func GetDefinedVariable(key string) interface{}{
	if val, ok := Application.definedVariables[key]; ok{
		return val
	}
	return nil
}


func GetHttpClient(name string) *httpclient.HttpClient{
	container := Application.serviceContainer
	handlers, ok := container["httpclient"]
	if !ok{
		return nil
	}

	clients := handlers.(map[string]*httpclient.HttpClient)
	if client, ok := clients[name]; ok {
		return client
	}

	return nil
}