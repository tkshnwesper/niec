package main

import (
	"github.com/kataras/iris"
	"niec/db"
	"niec/router"
	"niec/common"
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
	common.Init()
	router.Init()
	db.Init()
	iris.Listen(":8081")
	// iris.Go()
}