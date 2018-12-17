package gothic

/**
 * 框架入口
 * @author zhaojiangwei
 * @date 2018/12/05
 */
import (
	"time"
	"os"
	"path"
	"flag"
	"fmt"
	"github.com/gothinc/gothic/logger"
)

//请求体最大字节数
const defaultMaxMultipartMemory int64 = 32 << 20
const pathSplitSymbol = "/"

const (
	defaultBasePath = ".."
	defaultConfPath = "conf"
	defaultConfigNamePre = "app"
	defaultConfigType  = ".toml"
)

var Application = NewGothicApplication()

type GothicApplication struct{
	//全局上下文环境
	Context GothicContext

	BasePath string
	ConfigPath string
	ConfigFile string

	//当前模式(开发或线上环境)
	Active string

	//全局框架hook集合
	globalHooks map[HookPoint][]hookFunc

	maxMultipartMemory int64
}

func NewGothicApplication() *GothicApplication{
	return &GothicApplication{
		Context: GothicContext{},
		globalHooks: map[HookPoint][]hookFunc{},
		maxMultipartMemory: defaultMaxMultipartMemory,
		Active: "",
		BasePath: ".",
	}
}

func (this *GothicApplication) AddHook(point HookPoint, hook hookFunc){
	if _, ok := this.globalHooks[point]; !ok{
		this.globalHooks[point] = []hookFunc{}
	}

	this.globalHooks[point] = append(this.globalHooks[point], hook)
}

func (this *GothicApplication) Run(){
	//1. 初始化环境
	this.envInit()

	//2. 扫描controller初始化路由
}

func (this *GothicApplication) envInit() {
	println("------------------------------------------------------------------------")
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	println(timeNow, "Starting Application...")

	os.Chdir(path.Dir(os.Args[0]))
	this.parseFlag()

	//加载配置
	loadConfig(this.ConfigPath, this.ConfigFile, this.Active)

	this.initLogger()
}

func (this *GothicApplication) initLogger(){
	//get logger config
	loggerConfig := Config.GetStringMap("log")
	if loggerConfig == nil || len(loggerConfig) == 0{
		return
	}

	if _, ok := loggerConfig["root"]; ok{
		Logger.SetRootPath(Config.GetString("log.root"))
	}

	if _, ok := loggerConfig["prefix"]; ok{
		Logger.SetPrefix(Config.GetString("log.prefix"))
	}

	if _, ok := loggerConfig["suffix"]; ok{
		Logger.SetSuffix(Config.GetString("log.suffix"))
	}

	if _, ok := loggerConfig["level"]; ok{
		logLevel := logger.LogLevel(Config.GetInt("log.level"))
		if !logger.CheckLogLevel(logLevel){
			panic(fmt.Sprintf("invalid log level, config level: %d", logLevel))
		}
		
		Logger.SetLogLevel(logLevel)
	}

	if _, ok := loggerConfig["json_format"]; ok {
		isJson := Config.GetBool("log.json_format")
		Logger.SetIsJson(isJson)

		if isJson{
			Logger.SetFormatter(logger.NewJsonFormatter(logger.DefaultLogTimestampFormat))
		}
	}

	if _, ok := loggerConfig["timestamp_format"]; ok{
		Logger.SetTimestampFormat(Config.GetString("log.timestamp_format"))
	}

	if _, ok := loggerConfig["disable_time"]; ok{
		Logger.SetDisableTime(Config.GetBool("log.disable_time"))
	}
}

func (this *GothicApplication) parseFlag(){
	basePath := flag.String("b", defaultBasePath, "base path")
	confPath := flag.String("c", defaultConfPath, "config path")
	configFile := flag.String("f", defaultConfigNamePre + defaultConfigType, "config path")
	active := flag.String("r", "", "server region")

	flag.Parse()

	this.BasePath = *basePath
	this.ConfigPath = this.BasePath + pathSplitSymbol + *confPath
	this.ConfigFile = *configFile
	this.Active = *active
}