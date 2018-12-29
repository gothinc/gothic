package logic

import (
	"github.com/gothinc/gothic"
)

func NewDemoLogic(context *gothic.ThreadContext) *DemoLogic {
	return &DemoLogic{
		Context: context,
	}
}

type DemoLogic struct {
	Context *gothic.ThreadContext
}

func (demoLogic *DemoLogic) GetMsg() map[string]interface{} {
	gothic.Format(gothic.EntryFields{
		"logtype":   "system",
		"msg":       "demo logic",
		"raw_query": demoLogic.Context.Request.URL.RawQuery,
	}).Access()

	return map[string]interface{}{
		"from":   "gothic",
		"output": "json",
	}
}
