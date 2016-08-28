package main

import (
	"github.com/kataras/iris"
)

func main() {
	RouterInit()
	InitDB()
	iris.Listen(":8081")
}