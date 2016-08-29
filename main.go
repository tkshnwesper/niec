package main

import (
	"github.com/kataras/iris"
	"fmt"
)

func main() {
	fmt.Println("started server")
	RouterInit()
	InitDB()
	iris.Listen(":8081")
}