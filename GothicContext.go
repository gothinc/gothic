package gothic

import (
	"sync"
	"net/http"
)

/**
 * @desc Controller
 * @author zhaojiangwei
 * @date 2018-12-20 15:33
 */

var threadContextPool = sync.Pool{}

func NewThreadContext(request *http.Request, controller, action string) *ThreadContext{
	threadContext, ok := threadContextPool.Get().(*ThreadContext)
	if ok{
		threadContext.Request = request
		threadContext.Controller = controller
		threadContext.Action = action
		return threadContext
	}

	return &ThreadContext{
		Application: Application,
		Request: request,
		Controller: controller,
		Action: action,
		Params: make(map[string]interface{}),
	}
}

func ReleaseThreadContext(context *ThreadContext){
	context.Reset()
}

func (threadContext *ThreadContext) Reset(){
	threadContext.Request = nil
	threadContext.Controller = ""
	threadContext.Action = ""
	threadContext.Params = make(map[string]interface{})
}

//协程级别context, 保存单次请求链中的上下文环境
type ThreadContext struct{
	Application *GothicApplication
	Request     *http.Request
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
