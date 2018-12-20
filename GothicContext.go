package gothic

/**
 * @desc Controller
 * @author zhaojiangwei
 * @date 2018-12-20 15:33
 */

type ThreadContext struct{
	Application *GothicApplication
	Controller 	string
	Action 		string

	//用户自定义变量集合
	Params      map[string]interface{}
}
