package gothic

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelAccess
	LevelWarn
	LevelError
)

const(
	defaultLogRootPath = "./logs"
	defaultLogSuffix = ".log"
	defaultLogLevel = LevelAccess
	defaultLogTimestampFormat = "2006-01-02 15:04:05"
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

var Logger = &GothicLogger{
	rootPath: defaultLogRootPath,
	suffix: defaultLogSuffix,
	logLevel: defaultLogLevel,
	timestampFormat: defaultLogTimestampFormat,
}

func initLogger(){
	loggerConfig := Config.GetStringMap("log")
	if _, ok := loggerConfig["root"]; ok{
		Logger.rootPath = Config.GetString("log.root")
	}

	Logger.prefix = Config.GetString("log.prefix")

	if _, ok := loggerConfig["suffix"]; ok{
		Logger.suffix = Config.GetString("log.suffix")
	}

	if _, ok := loggerConfig["level"]; ok{
		Logger.logLevel = LogLevel(Config.GetInt("log.level"))
		if !checkLogLevel(Logger.logLevel){
			panic(fmt.Sprintf("invalid log level, config level: %d", Logger.logLevel))
		}
	}

	if _, ok := loggerConfig["timestamp_format"]; ok{
		Logger.timestampFormat = Config.GetString("log.timestamp_format")
	}

	Logger.isJson = Config.GetBool("log.json_format")
	Logger.loggerMap = make(map[string]*log.Logger)
	Logger.fdMap = make(map[string]*os.File)
}

//日志类
type GothicLogger struct {
	loggerMap 		map[string]*log.Logger
	fdMap     		map[string]*os.File
	rootPath  		string
	logLevel  		LogLevel
	prefix    		string
	suffix    		string
	timestampFormat string
	isJson  		bool
	mu        		sync.RWMutex
}

//判断文件或目录是否存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (this *GothicLogger) Debug(v ...interface{}) {
	if this.logLevel > LevelDebug {
		return
	}

	logName := this.prefix + getLevelName(LevelDebug) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Access(v ...interface{}) {
	this.Debug(v)
	if this.logLevel > LevelAccess {
		return
	}

	logName := this.prefix + getLevelName(LevelAccess) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Warn(v ...interface{}) {
	this.Debug(v)
	if this.logLevel > LevelWarn {
		return
	}

	logName := this.prefix + getLevelName(LevelWarn) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Error(v ...interface{}) {
	this.Debug(v)
	if this.logLevel > LevelError {
		return
	}

	logName := this.prefix + getLevelName(LevelError) + this.suffix
	this.writeLog(logName, v...)
}

func (this *GothicLogger) Extend(logName string, v ...interface{}) {
	logName = this.prefix + logName + this.suffix
	this.writeLog(logName, v...)
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

func (this *GothicLogger) writeLog(logName string, v ...interface{}) {
	logger, err := this.getLogger(this.rootPath + "/" + logName)
	if err != nil {
		timeNow := time.Now().Format(this.timestampFormat)
		fmt.Println(timeNow, "log failed", err)
		return
	}

	this.write(logger, v...)
}

func (this *GothicLogger) write(logger *log.Logger, v ...interface{}) {
	if nil == logger {
		return
	}

	if this.isJson{
		this.writeJson(logger, v...)
		return
	}

	msgstr := ""
	for _, msg := range v {
		if msg1, ok := msg.(map[string]interface{}); ok {
			//map每次输出的顺序是随机的，以下保证每次输出的顺序一致，如果map比较大，可能有一定性能损耗
			var keys []string
			for k := range msg1 {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				msgstr = msgstr + fmt.Sprintf("%s=%+v,", k, msg1[k])
			}
		} else {
			msgstr = msgstr + fmt.Sprintf("%+v,", msg)
		}
	}
	msgstr = strings.TrimRight(msgstr, ",")
	timeNow := time.Now().Format(this.timestampFormat)
	logger.Printf("%s %s\n", timeNow, msgstr)
}

func (this *GothicLogger) writeJson(logger *log.Logger, v ...interface{}) {

}

func checkLogLevel(level LogLevel) bool{
	if level < LevelDebug || level > LevelError{
		return false
	}

	return true
}