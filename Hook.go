package gothic

type SystemHookPoint int
type ThreadHookPoint int

const (
	BeforeConfingLoad SystemHookPoint = iota
	AfterConfigLoad
	BeforLoggerLoad
	AfterLoggerLoad
)

const(
	BeforeInitController ThreadHookPoint = iota
	BeforeInvokeAction
	AfterInvokeAction
	AfterSendResponse
)

type SystemHookChainType map[SystemHookPoint][]SystemHookFunc
type ThreadHookChainType map[string]map[string]map[ThreadHookPoint][]ThreadHookFunc

//全局框架hook集合
var SystemHookChain SystemHookChainType
//map[Controller]map[Action]HookFunc
var ThreadHookChain ThreadHookChainType

//框架系统级hook方法
type SystemHookFunc func() error
//单个请求链上的hook方法
type ThreadHookFunc func(context *ThreadContext) error

func AddSystemHook(point SystemHookPoint, hookFunc SystemHookFunc)  {
	if SystemHookChain == nil{
		SystemHookChain = make(SystemHookChainType)
	}

	if _, ok := SystemHookChain[point]; !ok{
		SystemHookChain[point] = make([]SystemHookFunc, 0)
	}

	SystemHookChain[point] = append(SystemHookChain[point], hookFunc)
}

func InvokeSystemHookChain(point SystemHookPoint) error{
	if SystemHookChain == nil{
		return nil
	}

	if funs, ok := SystemHookChain[point]; ok{
		var err error
		for _, fun := range funs{
			err = fun()
			if err != nil{
				break
			}
		}

		return err
	}else{
		return nil
	}
}

func AddThreadHook(controller, action string, point ThreadHookPoint, hookFunc ThreadHookFunc){
	if ThreadHookChain == nil{
		ThreadHookChain = make(ThreadHookChainType)
	}

	if _, ok := ThreadHookChain[controller]; !ok {
		ThreadHookChain[controller] = make(map[string]map[ThreadHookPoint][]ThreadHookFunc)
	}

	if _, ok := ThreadHookChain[controller][action]; !ok{
		ThreadHookChain[controller][action] = make(map[ThreadHookPoint][]ThreadHookFunc)
	}

	if _, ok := ThreadHookChain[controller][action][point]; !ok{
		ThreadHookChain[controller][action][point] = make([]ThreadHookFunc, 0)
	}

	ThreadHookChain[controller][action][point] = append(ThreadHookChain[controller][action][point], hookFunc)
}

func InvokeThreadHook(controller, action string, point ThreadHookPoint, context *ThreadContext) error{
	if ThreadHookChain == nil{
		return nil
	}

	if _, ok := ThreadHookChain[controller]; !ok {
		return nil
	}

	if _, ok := ThreadHookChain[controller][action]; !ok{
		return nil
	}

	if funs, ok := ThreadHookChain[controller][action][point]; !ok{
		var err error
		for _, fun := range funs {
			err = fun(context)
			if err != nil{
				break
			}
		}

		return err
	}

	return nil
}