package main

import (
    "github.com/kataras/iris"
)

func RouterInit() {
    iris.Get("/", func(c *iris.Context) {
        c.Write("Hello")
    })
}