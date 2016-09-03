package main

import (
    "github.com/kataras/iris"
    HTML "github.com/iris-contrib/template/html"
    "fmt"
    "html"
    "crypto/md5"
    "github.com/dchest/captcha"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
    "html/template"
    "io/ioutil"
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

func (ec ErrorContainer) Serve(c *iris.Context) {
    buildErrorPage(c, ec)
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

func getMarkdown(s string) string {
    return bluemonday.UGCPolicy().Sanitize(string(blackfriday.MarkdownCommon([]byte(s))))
}
    
func InitErrorPages() {

    iris.OnError(iris.StatusNotFound, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "404",
            "Sorry! The page you requested could not be found.",
        })
    })
    
    iris.OnError(iris.StatusInternalServerError, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "503",
            "Our server encountered an internal server error.",
        })
    })
    
    iris.OnError(iris.StatusForbidden, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "403",
            "Please sign in first.",
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
        "invalid-credentials": ErrorContainer {
            "Invalid Credentials",
            "The email or password that you entered are not valid.",
        },
        "user-does-not-exist": ErrorContainer {
            "User does not exist",
            "The user with that email address does not exist on our database.",
        },
        "incorrect-captcha": ErrorContainer {
            "Incorrect Captcha",
            "The captcha that you entered was not correct.",
        },
    }
    
    for s, ec := range errTypes {
        iris.Handle("GET", fmt.Sprintf("/error/%v", s), ErrorContainer {
            Title: ec.Title,
            Message: ec.Message,
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
    
    iris.Get("/learn-more", func(c *iris.Context) {
        c.Render("learn.more.html", struct {
            Title string
            Text template.HTML
        } {
            "Learn more",
            template.HTML(getMarkdown(readMD("learn.more.md"))),
        })
    })("learn-more")
    
    iris.Get("/sign/up", func(c *iris.Context) {
        renderSign(c, "Niec :: SignUp", "SignUp")
    })("signup")
    
    iris.Post("/sign/up", func(c *iris.Context) {
        if res, email, password := getCreds(c); res {
            if CheckEmailExists(email) {
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
        renderSign(c, "Niec :: SignIn", "SignIn")
    })("signin")
    
    iris.Post("/sign/in", func(c *iris.Context) {
        if res, email, password := getCreds(c); res {
            if CheckEmailExists(email) {
                if VerifyCreds(email, password) {
                    c.Session().Set("username", GetUsername(email))
                } else {
                    c.RedirectTo("invalid-credentials")
                }
            } else {
                c.RedirectTo("user-does-not-exist")
            }
        } else {
            c.RedirectTo("blank-field")
        }
    })
    
    iris.Get("/sign/up/next", func(c *iris.Context) {
        capid := captcha.New()
        c.Session().Set("capid", capid)
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
            CapID string
        }{
            "Niec :: SignUp - Next",
            fields,
            capid,
        })
    })("signup-next")
    
    iris.Post("/sign/up/next", func(c *iris.Context) {
        bio := html.EscapeString(c.FormValueString("bio"))
        username := html.EscapeString(c.FormValueString("username"))
        dp := html.EscapeString(c.FormValueString("dp"))
        retype := c.FormValueString("retype")
        cap := c.FormValueString("captcha")
        
        if c.Session().GetString("password") != retype {
            c.RedirectTo("password-mismatch")
        } else if CheckUsernameExists(username) {
            c.RedirectTo("username-already-taken")
        } else if c.Session().GetString("email") == "" {
            c.EmitError(iris.StatusInternalServerError)
        } else if !captcha.VerifyString(c.Session().Get("capid").(string), cap) {
            c.RedirectTo("incorrect-captcha")
        } else {
            if !InsertUser(
                c.Session().GetString("email"),
                fmt.Sprintf("%x", md5.Sum([]byte(retype))),
                username, 
                dp, 
                bio,
            ) {
                c.EmitError(iris.StatusInternalServerError)
            }
        }
    })
    
    var capHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
    iris.Get("/captcha/*id", iris.ToHandlerFunc(capHandler))("captcha")
    
    iris.Get("/submit", func(c *iris.Context) {
        if !isLoggedIn(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            buttons := []Field {
                {
                    "submit",
                    "submit",
                    "Submit",
                },
                {
                    "submit",
                    "preview",
                    "Preview",
                },
            }
            fields := []Field {
                {
                    "text",
                    "title",
                    "Title",
                },
                {
                    "text",
                    "tags",
                    "Tags (comma separated)",
                },
            }
            c.Render("submit.html", struct {
                Title string
                Fields []Field
                Buttons []Field
            } {
                "Submit an article",
                fields,
                buttons,
            })
        }
    })("submit")
    
    iris.Post("/submit", func(c *iris.Context) {
        action := c.FormValueString("action")
        body := c.FormValueString("body")
        // title := c.FormValueString("title")
        // tags := c.FormValueString("tags")
        if action == "preview" {
            c.Render("preview.html", struct {
                Title string
                Text template.HTML
            } {
                "Preview",
                template.HTML(getMarkdown(body)),
            })
        }
    })
}

func isLoggedIn(c *iris.Context) bool {
    return c.Session().Get("username") != ""
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

func renderSign(c *iris.Context, title, action string) {
    fields := []Field {
        {
            "email",
            "email",
            "Email Address",
        },
        {
            "password",
            "password",
            "Password",
        },
    }
    c.Render("sign.html", struct{
        Title string
        Action string
        Fields []Field
    }{
        title,
        action,
        fields,
    })
}

func readMD(name string) string {
    dat, err := ioutil.ReadFile("markdown/" + name)
    if !pe(err) {
        return ""
    }
    return string(dat)
}