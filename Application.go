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

var Application *GothicApplication

func init(){
	Application = NewGothicApplication()
}

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

	//加载日志组件
	initLogger()
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