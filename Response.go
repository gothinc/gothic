package gothic

/**
 * @desc Response
 * @author zhaojiangwei
 * @date 2018-12-18 11:26
 */

type Response interface {
	Serialize() interface{}
}

const DefaultResSucc int = 0
const DefaultResMsg string = "succ"

type GothicResponse struct {
	Errno int	`json:"Errno"`
	Errmsg string	`json:"Errmsg"`
	Data interface{} `json:"Data"`
}

func (this *GothicResponse) Serialize() interface{}{
	return this
}

func (this *GothicResponse) GenSucc(data interface{}) interface{}{
	this.Errno = DefaultResSucc
	this.Errmsg = DefaultResMsg
	this.Data = data
	return this
}