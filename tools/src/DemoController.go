package controller

import (
	"demo/src/configure"
	"demo/src/logic"
	"github.com/gothinc/gothic"
)

//示例Controller
type DemoController struct {
	gothic.Controller
}

func (this *DemoController) GetMsgAction() {
	name := this.GetString("name", "")
	if name == "" {
		panic(configure.ERR_INPUT)
	}

	DemoLogic := logic.NewDemoLogic(this.Context)
	data := DemoLogic.GetMsg()

	this.JsonSucc(data)
}
