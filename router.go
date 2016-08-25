package main

import (
    "github.com/kataras/iris"
    "github.com/iris-contrib/template/html"
)

func RouterInit() {
    // iris.StaticServe("./static/", "static")
    
    iris.UseTemplate(html.New(html.Config {
        Layout: "default-layout.html",
    }))
    
    iris.Get("/", func(c *iris.Context) {
        c.Render("index.html", struct{
            Title string
        }{
            "Hello world!",
        })
    })
}