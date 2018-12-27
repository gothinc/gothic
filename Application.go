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
	"github.com/gothinc/gothic/httpclient"
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
	basePath string
	configPath string
	configFile string

	//当前模式(开发或线上环境)
	active string

	maxMultipartMemory int64

	//用户自定义全局变量
	definedVariables map[string]interface{}

	serviceContainer map[string]interface{}
}

func NewGothicApplication() *GothicApplication{
	return &GothicApplication{
		maxMultipartMemory: defaultMaxMultipartMemory,
		active: "",
		basePath: ".",
		serviceContainer: make(map[string]interface{}),
	}
}

func (this *GothicApplication) Run(){
	//1. 初始化环境
	this.envInit()

	//2. 启动相关服务
	this.initService()

	//3. 启动服务
	this.startServer()
}

func (this *GothicApplication) initService(){
	this.startHttpClient()
}

func (this *GothicApplication) startHttpClient(){
	key := "application.httpclientpool"
	service := Config.GetStringMap(key)
	httpClientHandlers := map[string]*httpclient.HttpClient{}

	for name, _ := range service{
		if Config.GetBool(key + "." + name + ".disable"){
			continue
		}

		client := this.loadHttpClient(key + "." + name)
		httpClientHandlers[name] = client
	}

	this.serviceContainer["httpclient"] = httpClientHandlers
}

func (this *GothicApplication) loadHttpClient(key string) *httpclient.HttpClient{
	httpClient := httpclient.NewHttpClient(httpclient.HttpPoolSetting{
		MaxIdleConns: Config.GetInt(key + ".max_idle_conn"),
		IdleConnTimeout: Config.GetInt(key + ".idle_conn_timeout"),
		HttpTimeout: Config.GetInt(key + ".http_timeout"),
		DialTimeout: Config.GetInt(key + ".dial_timeout"),
	})

	println("start httpclient " + key + " succ")
	return httpClient
}

func (this *GothicApplication) AddController(controller interface{}){
	gothicHttpApiServer.AddController(controller)
}

func (this *GothicApplication) startServer(){
	port := Config.GetInt("server.http_port")
	if port == 0{
		port = DefaultListenPort
	}

	gothicHttpApiServer.HttpAddr = Config.GetString("server.http_addr")
	gothicHttpApiServer.HttpPort = port
	gothicHttpApiServer.ReadTimeout = Config.GetInt("server.http_read_timeout")
	gothicHttpApiServer.WriteTimeout = Config.GetInt("server.http_write_timeout")
	gothicHttpApiServer.IdleTimeout = Config.GetInt("server.http_idle_timeout")
	gothicHttpApiServer.hanlder.enablePprof = Config.GetBool("server.enable_pprof")

	gothicHttpApiServer.hanlder.Application = this

	gothicHttpApiServer.Run()
}

func (this *GothicApplication) envInit() {
	println("------------------------------------------------------------------------")
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	println(timeNow, "Starting Application...")

	os.Chdir(path.Dir(os.Args[0]))
	this.parseFlag()

	//加载配置
	InvokeSystemHookChain(BeforeConfingLoad)
	loadConfig(this.configPath, this.configFile, this.active)
	InvokeSystemHookChain(AfterConfigLoad)

	//加载日志组件
	InvokeSystemHookChain(BeforLoggerLoad)
	this.initLogger()
	InvokeSystemHookChain(AfterLoggerLoad)
}

func (this *GothicApplication) initLogger(){
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

	this.basePath = *basePath
	this.configPath = this.basePath + pathSplitSymbol + *confPath
	this.configFile = *configFile
	this.active = *active
}