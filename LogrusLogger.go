package gothic

/*
import (
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
)

const defaultLogTimePattern = "2006-01-02 15:04:05"

var Logger *GothicLog

type GothicLog struct {
	log.Logger
}

type LogWriterHook struct {
	LogLevelPath map[log.Level]string
	Writers map[log.Level]*os.File
}

func newLogWriterHook(level log.Level, path, prefix, suffix string) *LogWriterHook{
	levelPathMap := make(map[log.Level]string)
	levelWriterMap := make(map[log.Level]*os.File)

	for _, v := range log.AllLevels{
		if v == log.FatalLevel || v == log.PanicLevel{
			continue
		}

		if v < level{
			filePath := path + pathSplitSymbol + prefix + suffix
			levelPathMap[v] = filePath

			fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
			if err == nil {
				levelWriterMap[v] = fp
			} else {
				panic(fmt.Sprintf("Failed to log to file, level: %d", level))
			}
		}
	}

	return &LogWriterHook{
		LogLevelPath: levelPathMap,
	}
}

func (hook *LogWriterHook) Fire(entry *log.Entry) error{
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	entry.Writer()

	_, err = hook.Writers[entry.Level].WriteString(line)
	return err
}

func (hook *LogWriterHook) Levels() []log.Level {
	return []log.Level{
		log.TraceLevel,
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
	}
}

func initLogger(){
	formatter := getJsonFormatter()

	logger := log.New()
	logger.SetFormatter(formatter)
	configLevel := Config.GetString("log_level")
	logLevel,_ := log.ParseLevel(configLevel)
	logger.SetLevel(logLevel)

	addWriteHook(logger, configLevel, logLevel)
	logger.Info()

	Logger = &GothicLog{
		*logger,
	}
}

func addWriteHook(logger *log.Logger, levelStr string, levelInt log.Level){
	logPath := Config.GetString("log.log_root")
	logPrefix := Config.GetString("log.log_name")
	logSuffix := Config.GetString("log.log_split") + levelStr

	logger.AddHook(newLogWriterHook(levelInt, logPath, logPrefix, logSuffix))
}

func getJsonFormatter() *log.JSONFormatter{
	logTimePattern := Config.GetString("log_timestamp_pattern")
	if logTimePattern == ""{
		logTimePattern = defaultLogTimePattern
	}

	prettyPrint := Config.GetBool("pretty_print")
	formatter := &log.JSONFormatter{
		TimestampFormat:logTimePattern,
		PrettyPrint: prettyPrint,
	}

	return formatter
}
*/
