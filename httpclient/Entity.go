package httpclient

/**
 * @desc Entity
 * @author zhaojiangwei
 * @date 2018-12-25 14:44
 */

const(
	//未设置
	ClientCodeUnset = iota - 1
	//成功
	ClientSucc
	//New Request异常
	ClientNewReqFail
	//发起请求异常
	ClientDoReqFail
	//Http返回非200
	ClientHttpCodeFail
	//读取BODY异常
	ClientReadBodyFail
)

type ClientResponse struct{
	//方法返回错误码, 0:succ 其他:fail
	Code 		int			`json:"code"`
	//http状态码
	StatusCode	int			`json:"status_code""`
	//请求耗时, ms
	Cost 		int64		`json:"cost"`
	//错误信息
	Error 		error		`json:"error"`
	//http请求返回数据
	Data 		[]byte		`json:"data"`
}