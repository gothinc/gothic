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
	"github.com/gothinc/gothic/storage/redis"
	"strings"
	"strconv"
	"io/ioutil"
)

//请求体最大字节数
const defaultMaxMultipartMemory int64 = 32 << 20
const pathSplitSymbol = "/"

const (
	defaultBasePath = "."
	defaultConfPath = "conf"
	defaultConfigNamePre = "app"
	defaultConfigType  = "toml"
)

//配置文件类型
const (
	ConfigTypeToml = "toml"
	ConfigTypeYml = "yml"
)

var Application = NewGothicApplication()

type HttpClientContainer map[string]*httpclient.HttpClient
type RedisClientContainer map[string]*gothicredis.RedisClient

type GothicApplication struct{
	basePath string
	configPath string
	configFile string

	//当前模式(开发或线上环境)
	active string

	//日志类型
	configType string

	maxMultipartMemory int64

	//用户自定义全局变量
	definedVariables map[string]interface{}

	//服务容器(如httpclient连接池, redis连接池等)
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
	defer func() {
		if err := recover(); err != nil{
			this.managePid(false)
			println("application exit")
			fmt.Println(fmt.Sprintf("error: %s", err))
		}
	}()

	os.Chdir(path.Dir(os.Args[0]))
	this.parseFlag()

	//1. 初始化环境
	this.envInit()

	//2. 启动相关服务
	this.initService()

	//3. 生成PID
	this.managePid(true)

	//4. 启动服务
	this.startServer()
}

func (this *GothicApplication) managePid(create bool) {
	pidFile := Config.GetString("application.pid_file")
	if create {
		pid := os.Getpid()
		pidString := strconv.Itoa(pid)
		ioutil.WriteFile(pidFile, []byte(pidString), 0777)
	} else {
		os.Remove(pidFile)
	}
}

func (this *GothicApplication) GetHttpClient(name string) *httpclient.HttpClient{
	container := this.serviceContainer
	handlers, ok := container["httpclient"]
	if !ok{
		return nil
	}

	clients := handlers.(HttpClientContainer)
	if client, ok := clients[name]; ok {
		return client
	}

	return nil
}

func (this *GothicApplication) GetRedisClient(name string) *gothicredis.RedisClient{
	container := this.serviceContainer
	handlers, ok := container["redisclient"]
	if !ok{
		return nil
	}

	clients := handlers.(RedisClientContainer)
	if client, ok := clients[name]; ok {
		return client
	}

	return nil
}

func (this *GothicApplication) initService(){
	this.startHttpClient()
	this.startRedisClient()
}

func (this *GothicApplication) startHttpClient(){
	key := "application.httpclientpool"
	service := Config.GetStringMap(key)
	httpClientHandlers := HttpClientContainer{}

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

func (this *GothicApplication) startRedisClient(){
	key := "application.redisclient"
	service := Config.GetStringMap(key)
	redisClientHandlers := RedisClientContainer{}

	for name, _ := range service{
		if Config.GetBool(key + "." + name + ".disable"){
			continue
		}

		client := this.loadRedisClient(key + "." + name)
		redisClientHandlers[name] = client
	}

	this.serviceContainer["redisclient"] = redisClientHandlers
}

func (this *GothicApplication) loadRedisClient(key string) *gothicredis.RedisClient{
	redisClient := gothicredis.NewRedisClient(&gothicredis.RedisPoolConfig{
		MaxIdle:      Config.GetInt(key + ".max_idle_conn"),
		IdleTimeout:  Config.GetInt(key + ".idle_timeout"),
		Host:       Config.GetString(key + ".host"),
		Port: 		Config.GetString(key + ".port"),
		Password:     Config.GetString(key + ".password"),
		ConnTimeout:  Config.GetInt(key + ".conn_timeout"),
		ReadTimeout:  Config.GetInt(key + ".read_timeout"),
		WriteTimeout: Config.GetInt(key + ".write_timeout"),
	})

	println("start redisclient " + key + " succ")
	return redisClient
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

	//加载配置
	InvokeSystemHookChain(BeforeConfingLoad)
	loadConfig(this.configPath, this.configFile, this.configType, this.active)
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
	basePath := flag.String("b", defaultBasePath, "optional: root path")
	confPath := flag.String("c", defaultConfPath, "optional: config path")
	configType := flag.String("t", defaultConfigType, "optional: config type, yml|toml")
	configFile := defaultConfigNamePre + "." + *configType
	active := flag.String("r", " ", "optional: active mode")

	flag.Parse()

	this.basePath = *basePath
	this.configPath = this.basePath + pathSplitSymbol + *confPath
	this.configType = *configType
	this.configFile = configFile
	this.active = strings.TrimSpace(*active)
}