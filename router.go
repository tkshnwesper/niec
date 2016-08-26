package main

import (
    "github.com/kataras/iris"
    "github.com/iris-contrib/template/html"
    "fmt"
    // "encoding/json"
    // "github.com/iris-contrib/template/amber"
)

type ErrorContainer struct {
    Title string
    Message string
}

type Field struct {
    Type string
    Name string
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
        "blank-field": ErrorContainer {
            "Blank Field(s)",
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
    // iris.UseTemplate(amber.New()).Directory("./templates", ".html")
    
    iris.StaticServe("./static/", "static")
    
    InitErrorPages()
    
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
            c.RedirectTo("signup-next")
        } else {
            c.RedirectTo("blank-field")
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
    
    iris.Get("/sign/up/next", func(c *iris.Context) {
        fields := []Field {
            {
                "text",
                "Name",
            },
        }
        c.Render("sign.up.next.html", struct{
            Title string
            Fields []Field
        }{
            "Niec :: SignUp - Next",
            fields,
        })
    })("signup-next")
}

func getCreds(c *iris.Context) (bool, string, string) {
    c.Session().Clear()
    email := c.FormValueString("email")
    password := c.FormValueString("password")
    if email == "" || password == "" {
        return false, "", ""
    }
    return true, email, password
}