package router

import (
    "github.com/kataras/iris"
    HTML "github.com/iris-contrib/template/html"
    "github.com/dchest/captcha"
    "html/template"
    "niec/common"
)

// Field holds information about the input fields to be displayed in the view
type Field struct {
    Type string
    Name string
    Placeholder string
}

// Init helps to initialize all the pages required in the site
func Init() {
    iris.UseTemplate(HTML.New(HTML.Config {
        Layout: "layout0.html",
    }))
    
    iris.StaticServe("./static/", "static")
    
    initErrorPages()
    
    initSignPages()
    
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
            template.HTML(common.GetMarkdown(common.ReadMD("learn.more.md"))),
        })
    })("learn-more")
    
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
            Fields := []Field {
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
                Fields,
                buttons,
            })
        }
    })("submit")
    
    iris.Post("/preview", func(c *iris.Context) {
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
                template.HTML(common.GetMarkdown(body)),
            })
        }
    })("preview")
}

var pe = common.Pe