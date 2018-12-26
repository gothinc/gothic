package gothic

import (
	"sync"
)

/**
 * @desc Controller
 * @author zhaojiangwei
 * @date 2018-12-20 15:33
 */

var threadContextPool = sync.Pool{}

func NewThreadContext(controller, action string) *ThreadContext{
	threadContext, ok := threadContextPool.Get().(*ThreadContext)
	if ok{
		return threadContext
	}

	return &ThreadContext{
		Application: Application,
		Controller: controller,
		Action: action,
		Params: make(map[string]interface{}),
	}
}

func ReleaseThreadContext(context *ThreadContext){
	context.Params = make(map[string]interface{})
	context.Controller = ""
	context.Action = ""
	threadContextPool.Put(context)
}

//协程级别context, 保存单次请求链中的上下文环境
type ThreadContext struct{
	Application *GothicApplication
	Controller 	string
	Action 		string

	//用户自定义变量集合
	Params      map[string]interface{}
}

func (threadContext *ThreadContext) GetParamString(key string) string{
	val, ok := threadContext.Params[key]
	if !ok{
		return ""
	}

	valStr, ok := val.(string)
	if !ok{
		return ""
	}

	return valStr
}

func (threadContext *ThreadContext) Reset(){
	threadContext.Application = Application
	threadContext.Controller = ""
	threadContext.Action = ""
	threadContext.Params = nil
}
