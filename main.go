package main

import (
	"github.com/kataras/iris"
)

func main() {
	RouterInit()
	iris.Listen(":8081")
}