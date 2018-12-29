package configure

import "github.com/gothinc/gothic"

var (
	ERR_SUCC   = gothic.GothicResponseBody{Errno: 0, Errmsg: "OK"}
	ERR_SYSTEM = gothic.GothicResponseBody{Errno: 100, Errmsg: "system error"}
	ERR_INPUT  = gothic.GothicResponseBody{Errno: 101, Errmsg: "input param error"}
)
