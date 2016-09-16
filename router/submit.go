package router

import (
    "github.com/kataras/iris"
    "html/template"
    "niec/common"
)

func initSubmit() {
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