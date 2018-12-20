package gothic

/**
 * @desc Response
 * @author zhaojiangwei
 * @date 2018-12-18 11:26
 */

const DefaultGothicResSucc int = 0
const DefaultGothicResMsg string = "succ"

type GothicResponseBody struct {
	Errno int	`json:"Errno"`
	Errmsg string	`json:"Errmsg"`
	Data interface{} `json:"Data"`
}