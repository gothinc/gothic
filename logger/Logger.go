package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"github.com/json-iterator/go"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelAccess
	LevelWarn
	LevelError
)

const(
	DefaultLogRootPath = "./logs"
	DefaultLogPrefix = "gothic-"
	DefaultLogSuffix = ".log"
	DefaultLogLevel = LevelAccess
	DefaultLogTimestampFormat = "2006-01-02 15:04:05"
)

func getLevelName(level LogLevel) string{
	switch level {
	case LevelDebug:
		return "debug"
	case LevelAccess:
		return "access"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return ""
	}
}

func NewDefaultLogger() *GothicLogger{
	return &GothicLogger{
		rootPath: DefaultLogRootPath,
		logLevel: DefaultLogLevel,
		prefix: DefaultLogPrefix,
		suffix: DefaultLogSuffix,
		timestampFormat: DefaultLogTimestampFormat,
		isJson: false,
		disableTime: false,
		formatter: NewTextFormatter(DefaultLogTimestampFormat),
		loggerMap: make(map[string]*log.Logger),
		fdMap: make(map[string]*os.File),
	}
}

//日志类
type GothicLogger struct {
	loggerMap       map[string]*log.Logger
	fdMap           map[string]*os.File
	rootPath        string
	logLevel        LogLevel
	prefix          string
	suffix          string
	timestampFormat string
	isJson          bool
	formatter       GothicLogFormatter
	mu              sync.RWMutex
	disableTime     bool
}

func (this *GothicLogger) SetTimestampFormat(timestampFormat string){
	if this.disableTime{
		this.timestampFormat = ""
		this.formatter.SetTimestampFormat("")
	}else{
		this.timestampFormat = timestampFormat
		this.formatter.SetTimestampFormat(timestampFormat)
	}
}

func (this *GothicLogger) SetDisableTime(disableTime bool){
	this.disableTime = disableTime

	if disableTime{
		this.SetTimestampFormat("")
	}
}

func (this *GothicLogger) SetFormatter(formatter GothicLogFormatter){
	this.formatter = formatter
}

func (this *GothicLogger) SetLogLevel(logLevel LogLevel){
	this.logLevel = logLevel
}

func (this *GothicLogger) SetPrefix(prefix string){
	this.prefix = prefix
}

func (this *GothicLogger) SetSuffix(suffix string){
	this.suffix = suffix
}

func (this *GothicLogger) SetRootPath(rootPath string){
	this.rootPath = rootPath
}

func (this *GothicLogger) SetIsJson(isJson bool){
	this.isJson = isJson
}

func (this *GothicLogger) Format(fields EntryFields) *Entry{
	entry := NewEntry(this, fields)
	return entry
}

func (this *GothicLogger) Debug(v ...interface{}) {
	if this.logLevel > LevelDebug {
		return
	}

	logName := this.prefix + getLevelName(LevelDebug) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Access(v ...interface{}) {
	if this.logLevel > LevelAccess {
		return
	}

	logName := this.prefix + getLevelName(LevelAccess) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Warn(v ...interface{}) {
	if this.logLevel > LevelWarn {
		return
	}

	logName := this.prefix + getLevelName(LevelWarn) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Error(v ...interface{}) {
	if this.logLevel > LevelError {
		return
	}

	logName := this.prefix + getLevelName(LevelError) + this.suffix
	this.writeLog(logName, v...)
}

/**
 * @desc 扩展log，如play.log
 * @author zhaojiangwei
 * @date 18:29 2018/12/12
 * @param logName:日志名，如play
 * @return
 **/
func (this *GothicLogger) Extend(logName string, v ...interface{}) {
	logName = this.prefix + logName + this.suffix
	this.writeLog(logName, v...)
}

//判断文件或目录是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (this *GothicLogger) getLogger(logName string) (*log.Logger, error) {
	filePath := logName
	this.mu.RLock()
	fd, ok := this.fdMap[logName]
	this.mu.RUnlock()
	//如果日志文件没有打开，或者日志名已经变了，就重新打开另外的日志文件
	if !ok || (fd != nil && fd.Name() != filePath) || !PathExist(filePath) {
		this.mu.Lock()
		defer this.mu.Unlock()
		fd, ok = this.fdMap[logName]
		//双重判断，减少重复操作
		if !ok || (fd != nil && fd.Name() != filePath) || !PathExist(filePath) {
			if fd != nil {
				fd.Close()
			}
			fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
			if err != nil {
				return nil, err
			}

			fd.Chmod(0777)
			this.loggerMap[logName] = log.New(fd, "", 0)
			this.fdMap[logName] = fd
			timeNow := time.Now().Format(this.timestampFormat)
			fmt.Println(timeNow, "new logger:", filePath)
		}
	}
	retLogger, ok := this.loggerMap[logName]
	return retLogger, nil
}
func (this *GothicLogger) writeFields(logName string, fields EntryFields) {
	logger, err := this.getLogger(this.rootPath + "/" + logName)
	if err != nil {
		timeNow := time.Now().Format(this.timestampFormat)
		fmt.Println(timeNow, "get logger failed", err)
		return
	}

	msg, err := this.formatter.FormatFields(fields)
	if err != nil{
		fmt.Println(fmt.Sprintf("format msg error[%s], logName[%s], msg[%+v]",
			err.Error(), logName, fields))
	}else{
		logger.Println(msg)
	}
}

func (this *GothicLogger) writeLog(logName string, v ...interface{}) {
	logger, err := this.getLogger(this.rootPath + "/" + logName)
	if err != nil {
		timeNow := time.Now().Format(this.timestampFormat)
		fmt.Println(timeNow, "get logger failed", err)
		return
	}

	msg, err := this.formatter.Format(v...)
	if err != nil{
		fmt.Println(fmt.Sprintf("format msg error[%s], logName[%s], msg[%+v]",
			err.Error(), logName, v))
	}else{
		logger.Println(msg)
	}
}

func (this *GothicLogger) writeJson(logger *log.Logger, v ...interface{}) {
	if len(v) != 1{
		return
	}

	_, err := jsoniter.MarshalToString(v[0])
	if err != nil{
		fmt.Sprintln("marshal json exception: ", err.Error())
	}
}

func CheckLogLevel(level LogLevel) bool{
	if level < LevelDebug || level > LevelError{
		return false
	}

	return true
}