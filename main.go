package main

import (
	"github.com/kataras/iris"
)

func main() {
	RouterInit()
	InitDB()
	iris.Listen("192.168.1.4:8081")
}