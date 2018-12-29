package main

import (
	"demo/src/controller"
	"github.com/gothinc/gothic"
)

func main() {
	gothic.Application.AddController(&controller.DemoController{})
	gothic.Application.Run()
}
