package main

import (
    "github.com/kataras/iris"
    HTML "github.com/iris-contrib/template/html"
    "fmt"
    "html"
    "crypto/md5"
)

type ErrorContainer struct {
    Title string
    Message string
}

type Field struct {
    Type string
    Name string
    Placeholder string
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
        "password-mismatch": ErrorContainer {
            "Password Mismatch",
            "Your passwords don't match, please try again.",
        },
        "email-already-exists": ErrorContainer {
            "Email Already Exists",
            "This email address already exists on our database.",
        },
        "username-already-taken": ErrorContainer {
            "Username already taken",
            "This username has already been taken by another user.",
        },
    }
    
    for s, ec := range errTypes {
        iris.Get(fmt.Sprintf("/error/%v", s), func(c *iris.Context) {
            buildErrorPage(c, ec)
        })(s)
    }
}

func RouterInit() {
    iris.UseTemplate(HTML.New(HTML.Config {
        Layout: "layout0.html",
    }))
    
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
            if CheckDuplicateEmail(email) {
                c.RedirectTo("email-already-exists")
            } else {
                c.Session().Set("email", email)
                c.Session().Set("password", password)
                c.RedirectTo("signup-next")
            }
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
                "username",
                "Username",
            },
            {
                "password",
                "retype",
                "Retype password",
            },
            {
                "url",
                "dp",
                "Display picture URL (Optional)",
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
    
    iris.Post("/sign/up/next", func(c *iris.Context) {
        bio := html.EscapeString(c.FormValueString("bio"))
        username := html.EscapeString(c.FormValueString("username"))
        dp := html.EscapeString(c.FormValueString("dp"))
        retype := html.EscapeString(c.FormValueString("retype"))
        
        if c.Session().GetString("password") != retype {
            c.RedirectTo("password-mismatch")
        } else if CheckDuplicateUsername(username) {
            c.RedirectTo("username-already-taken")
        } else {
            InsertUser(
                c.Session().GetString("email"), 
                fmt.Sprintf("%x", md5.Sum([]byte(retype))), 
                username, 
                dp, 
                bio,
            )
        }
    })
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