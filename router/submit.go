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
            c.EmitError(iris.StatusUnauthorized)
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
                    true,
                },
                // {
                //     "text",
                //     "tags",
                //     "Tags (comma separated)",
                //     255,
                // },
            }
            cb := []Checkbox {
                {
                    "Public", "privacy", "public",
                    "Can be viewed without logging in", "globe", false,
                },
                {
                    "Draft", "draft", "draft",
                    "Drafts are visible only to you", "blackboard", false,
                },
            } 
            c.Render("submit.html", struct {
                Title, Textarea string
                Property Property
                Fields []Field
                Buttons []Button
                Checkboxes []Checkbox
            } {
                "Submit an Article",
                "",
                getProperty(c),
                Fields,
                buttons,
                cb,
            })
        }
    })("submit")
    
    iris.Post("/submit", func(c *iris.Context) {
        if !isLoggedIn(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            text := c.FormValueString("text")
            title := c.FormValueString("title")
            // tags := c.FormValueString("tags")
            action := c.FormValueString("action")
            privacy := c.FormValueString("privacy")
            var draft = false
            if c.FormValueString("draft") == "draft" {
                draft = true
            }
            if action == "preview" {
                c.SetFlash("text", text)
                c.RedirectTo("preview")
            } else if action == "submit" {
                var pub = false
                if privacy == "public" {
                    pub = true
                }
                id, success := db.InsertArticle(
                    getUserID(c),
                    title,
                    // tags,
                    text,
                    pub,
                    draft,
                )
                if !success {
                    c.EmitError(iris.StatusInternalServerError)
                } else {
                    msg := "Article submitted successfully!"
                    if draft {
                        msg = "Draft saved successfully!"
                    }
                    c.SetFlash("message", msg)
                    c.RedirectTo("article", id)
                }
            } else {
                c.EmitError(iris.StatusNotFound)
            }
        }
    })
    
    iris.Get("/preview", func(c *iris.Context) {
        body, err := c.GetFlash("text")
        if err != nil {
            c.EmitError(iris.StatusNoContent)
        } else {
            c.Render("preview.html", struct {
                Title string
                Property Property
                Text template.HTML
            } {
                "Preview",
                getProperty(c),
                template.HTML(common.GetMarkdown(body)),
            })
        }
    })("preview")
}