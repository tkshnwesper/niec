package main

import (
	"github.com/kataras/iris"
	// "github.com/kataras/iris/config"
)

func main() {
	// def := config.DefaultServer()
	// man := config.Server {
	// 	Virtual: true,
	// 	VListeningAddr: ":80",
	// 	ListeningAddr: "",
	// 	VScheme: "/",
	// }
	// conf := def.MergeSingle(man)
	// iris.Servers.Add(conf)
	RouterInit()
	InitDB()
	iris.Listen(":8081")
	// iris.Go()
	
}