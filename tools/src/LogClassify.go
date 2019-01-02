package configure

import "time"

/**
 * @desc 统一定义日志格式
 * @author zhaojiangwei
 * @date 2018-12-24 18:44
 */

//顶级模块
type LogPrimaryModule string

const (
	LogSystem LogPrimaryModule = "system" //系统日志
	LogRpc    LogPrimaryModule = "rpc"    //接口调用相关日志
	LogRedis  LogPrimaryModule = "redis"  //redis日志
	LogDb     LogPrimaryModule = "db"     //数据库日志
	LogBiz    LogPrimaryModule = "biz"    //业务日志
)

//子模块(错误号)
type LogChildModule string

//系统顶级模块下的子模块(错误号100开始)
const (
	LogSystemInfo  LogChildModule = "info"      //普通系统日志
	LogSystemError LogChildModule = "exception" //异常系统日志
)

//接口调用顶级模块下的子模块(错误号200开始)
const (
	LogRpcInfo             LogChildModule = "info"                //普通接口调用日志
	LogRpcTimeout          LogChildModule = "timeout"             //接口调用超时
	LogRpcException        LogChildModule = "response_exception"  //接口请求异常
	LogRpcResContException LogChildModule = "response_data_error" //接口正常，但返回数据非预期
	LogRpcDataInvalid      LogChildModule = "data_invalid"        //接口调用返回数据格式异常
)

//REDIS顶级模块下的子模块(错误号300开始)
const (
	LogRedisInfo        LogChildModule = "info"         //普通REDIS日志
	LogRedisConnTimeout LogChildModule = "conn_timeout" //redis连接超时
	LogRedisRwTimeout   LogChildModule = "rw_timeout"   //redis读写超时
)

//DB顶级模块下的子模块(错误号400开始)
const (
	LogDbInfo    LogChildModule = "info"    //普通DB日志
	LogDbTimeout LogChildModule = "timeout" //db超时
)

//Logic顶级模块下的子模块(错误号400开始)
const (
	LogBizInfo      LogChildModule = "info"      //普通业务日志
	LogBizTimeout   LogChildModule = "timeout"   //业务超时
	LogBizException LogChildModule = "exception" //业务异常
)

type LogFormat struct {
	PrimaryModule LogPrimaryModule       `json:"module"` //主模块
	ChildModule   LogChildModule         `json:"code"`   //子模块
	Time          string                 `json:"time"`
	LogId         string                 `json:"logid"`
	Detail        map[string]interface{} `json:"detail"` //其他信息
}

func NewLogFormat(module LogPrimaryModule, childModule LogChildModule, logId string, detail map[string]interface{}) LogFormat {
	return LogFormat{
		PrimaryModule: module,
		ChildModule:   childModule,
		LogId:         logId,
		Detail:        detail,
		Time:          time.Now().Format("2006-01-02 15:04:05"),
	}
}
