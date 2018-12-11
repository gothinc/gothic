package gothic

import (
	"net/http"
)

type GothicContext struct {
	Request *http.Request
	Logid string
	Params map[string]interface{}
}
