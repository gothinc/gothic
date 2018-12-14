package logger

/**
 * @desc Export
 * @author zhaojiangwei
 * @date 2018-12-13 15:56
 */

var gothicLogger = newGothicLogger()

func Debug(v ...interface{})  {
	gothicLogger.Debug(v)
}

func Access(v ...interface{})  {
	gothicLogger.Access(v)
}

func Warn(v ...interface{})  {
	gothicLogger.Warn(v)
}

func Error(v ...interface{})  {
	gothicLogger.Error(v)
}

func Extend(logName string, v ...interface{})  {
	gothicLogger.Extend(logName, v)
}