package router

import (
    "github.com/kataras/iris"
    "html/template"
    "niec/common"
    "niec/db"
)

func initSubmitPages() {
    iris.Get("/submit", func(c *iris.Context) {
        if !isLoggedIn(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            buttons := []Button {
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
                    255,
                    "",
                },
                // {
                //     "text",
                //     "tags",
                //     "Tags (comma separated)",
                //     255,
                // },
            }
            c.Render("submit.html", struct {
                Title, Textarea string
                Fields []Field
                Buttons []Button
            } {
                "Submit an article",
                "",
                Fields,
                buttons,
            })
        }
    })("submit")
    
    iris.Post("/submit", func(c *iris.Context) {
        if !isLoggedIn(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            text := c.FormValueString("text")
            title := c.FormValueString("title")
            tags := c.FormValueString("tags")
            action := c.FormValueString("action")
            if action == "preview" {
                c.SetFlash("text", text)
                c.RedirectTo("preview")
            } else if action == "submit" {
                if !db.InsertArticle(
                    c.Session().GetString(common.UserIdentificationAttribute), 
                    title, 
                    tags, 
                    text,
                ) {
                    c.EmitError(iris.StatusInternalServerError)
                }
            } else {
                c.EmitError(iris.StatusNotFound)
            }
        }
    })
    
    iris.Get("/preview", func(c *iris.Context) {
        body, _ := c.GetFlash("text")
        c.Render("preview.html", struct {
            Title string
            Text template.HTML
        } {
            "Preview",
            template.HTML(common.GetMarkdown(body)),
        })
    })("preview")
    
    iris.Get("/article/:id/edit", func(c *iris.Context) {
        if isLoggedIn(c) {
            id, err := c.ParamInt64("id")
            if !pe(err) {
                c.EmitError(iris.StatusNotFound)
            } else {
                if db.GetUserID(c.Session().GetString(common.UserIdentificationAttribute)) == db.GetArticleUserID(id) {
                    title, text := db.FetchForEdit(id)
                    buttons := []Button {
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
                            255,
                            title,
                        },
                        // {
                        //     "text",
                        //     "tags",
                        //     "Tags (comma separated)",
                        //     255,
                        // },
                    }
                    c.Render("submit.html", struct {
                        Title, Textarea string
                        Fields []Field
                        Buttons []Button
                    }{
                        "Edit Article",
                        text,
                        Fields,
                        buttons,
                    })
                } else {
                    c.EmitError(iris.StatusForbidden)
                }
            }
        } else {
            c.EmitError(iris.StatusUnauthorized)
        }
    })
    
    iris.Post("/article/:id/edit", func(c *iris.Context) {
        if isLoggedIn(c) {
            id, err := c.ParamInt64("id")
            text := c.FormValueString("text")
            title := c.FormValueString("title")
            if !pe(err) {
                c.EmitError(iris.StatusNotFound)
            } else {
                if db.GetUserID(c.Session().GetString(common.UserIdentificationAttribute)) == db.GetArticleUserID(id) {
                    if db.EditArticle(id, title, text) {
                        c.SetFlash("message", "Article updated successfully!")
                        c.RedirectTo("article", id)
                    } else {
                        c.EmitError(iris.StatusInternalServerError)
                    }
                } else {
                    c.EmitError(iris.StatusForbidden)
                }
            }
        } else {
            c.EmitError(iris.StatusUnauthorized)
        }
    })
}