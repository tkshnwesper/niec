package main

import (
    "github.com/kataras/iris"
    "github.com/iris-contrib/template/html"
    "fmt"
)

type ErrorContainer struct {
    Title string
    Message string
}

func buildErrorPage(c *iris.Context, err ErrorContainer) {
    c.Render("error.html", struct {
        Title string
        ErrorTitle string
        ErrorMessage string
    } {
        fmt.Sprintf("Error: %v", err.Title),
        err.Title,
        err.Message,
    })
}
    
func InitErrorPages() {
    iris.OnError(iris.StatusNotFound, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "404",
            "Sorry! The page you requested could not be found.",
        })
    })
    
    errTypes := map[string]ErrorContainer {
        "blank-field-error": ErrorContainer {
            "Blank Field(s) Error",
            "Kindly enter all the required fields.",
        },
    }
    
    for s, ec := range errTypes {
        iris.Get(fmt.Sprintf("/error/%v", s), func(c *iris.Context) {
            buildErrorPage(c, ec)
        })(s)
    }
}

func RouterInit() {
    iris.UseTemplate(html.New(html.Config {
        Layout: "layout0.html",
    }))
    
    iris.StaticServe("./static/", "static")
    
    iris.Get("/", func(c *iris.Context) {
        c.Render("index.html", struct{
            Title string
        }{
            "Welcome to Niec!",
        })
    })("landing")
    
    iris.Get("/sign/up", func(c *iris.Context) {
        c.Render("sign.html", struct{
            Title string
            Action string
        }{
            "Niec :: SignUp",
            "SignUp",
        })
    })("signup")
    
    iris.Post("/sign/up", func(c *iris.Context) {
        if res, email, password := getCreds(c); res {
            c.Session().Set("email", email)
            c.Session().Set("password", password)
            c.Render("sign.up.next.html", struct{
                Title string
                Action string
            }{
                "Niec :: SignUp",
                "SignUp",
            })
        } else {
            c.RedirectTo("blank-field-error")
        }
    })
    
    iris.Get("/sign/in", func(c *iris.Context) {
        c.Render("sign.html", struct{
            Title string
            Action string
        }{
            "Niec :: SignIn",
            "SignIn",
        })
    })("signin")
}

func getCreds(c *iris.Context) (bool, string, string) {
    email := c.FormValueString("email")
    password := c.FormValueString("password")
    if email == "" || password == "" {
        return false, "", ""
    }
    return true, email, password
}