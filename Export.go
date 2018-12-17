package gothic

import (
	"github.com/gothinc/gothic/logger"
)

/**
 * @desc Export
 * @author zhaojiangwei
 * @date 2018-12-13 15:56
 */

var Logger = logger.NewDefaultLogger()

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