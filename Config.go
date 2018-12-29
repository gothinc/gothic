package gothic

/**
 * 配置文件
 * @author zhaojiangwei
 * @date 2018/12/06
 */
import (
	"github.com/spf13/viper"
	"fmt"
)

const includeKeyPre = "includes.files"

var Config = &GothicConfig{}

type GothicConfig struct {
	viper.Viper
}

func GetConfig() *GothicConfig{
	return Config
}

/**
 * 加配置文件
 */
func loadConfig(configPath, configFile, configType, active string){
	println("loading config:" +configFile + "...")

	entity := viper.New()
	entity.SetConfigFile(configPath + pathSplitSymbol + configFile)
	err := entity.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file[%s]: %s \n", configFile, err))
	}

	//处理include包含的配置文件
	includeConfig := entity.GetStringSlice(includeKeyPre)
	for _, v := range includeConfig{
		loadInclude(entity, configPath, v)
	}

	//包含当前环境（线上或开发）的配置
	if active != "" && active != defaultActiveMode {
		loadInclude(entity, configPath, defaultConfigNamePre + "." + active + "." + configType)
	}

	Config = &GothicConfig{
		*entity,
	}
}

/**
 * 加载包含的配置文件
 */
func loadInclude(mainConfig *viper.Viper, configPath, configFile string){
	println("loading config:" + configFile + "...")
	child := viper.New()
	child.SetConfigFile(configPath + pathSplitSymbol + configFile)
	err := child.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file[%s]: %s \n", configFile, err))
	}

	setting := child.AllSettings()
	var includeConfig []string
	if child.InConfig("includes") {
		includeConfig = child.GetStringSlice(includeKeyPre)

		//remove "includes key"
		delete(setting, includeKeyPre)
	}

	if err := mainConfig.MergeConfigMap(setting); err != nil{
		panic(fmt.Errorf("merge config error, file[%s], err[%s]", configFile, err))
	}

	if includeConfig != nil {
		for _, v := range includeConfig {
			loadInclude(mainConfig, configPath, v)
		}
	}
}