package main

import (
    "github.com/kataras/iris"
    "github.com/iris-contrib/template/html"
)

func RouterInit() {
    iris.UseTemplate(html.New(html.Config {
        Layout: "layout0.html",
    }))
    
    iris.StaticServe("./static/", "static")
    
    iris.Get("/", func(c *iris.Context) {
        c.Render("index.html", struct{
            Title string
        }{
            "Hello world!",
        })
    })
}