package router

import (
    "github.com/kataras/iris"
    "html/template"
    "niec/common"
    "niec/db"
)

func initSubmitPages() {
    iris.Get("/submit", func(c *iris.Context) {
        if isLoggedIn(c) {
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
    
    iris.Get("/article/:id/edit", func(c *iris.Context) {
        if isLoggedIn(c) {
            id, err := c.ParamInt64("id")
            if !pe(err) {
                c.EmitError(iris.StatusNotFound)
            } else {
                if getUserID(c) == db.GetArticleUserID(id) {
                    title, text, pub, draft := db.FetchForEdit(id)
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
                    cb := []Checkbox {
                        {
                            "Public", "privacy", "public",
                            "Can be viewed without logging in", "globe", pub,
                        },
                        {
                            "Draft", "draft", "draft",
                            "Drafts are visible only to you", "blackboard", draft,
                        },
                    }
                    c.Render("edit.html", struct {
                        // Do not make Textarea into a template.HTML
                        Property Property
                        Title, Textarea string
                        Fields []Field
                        Buttons []Button
                        Checkboxes []Checkbox
                    }{
                        getProperty(c),
                        "Edit Article",
                        text,
                        Fields,
                        buttons,
                        cb,
                    })
                } else {
                    c.EmitError(iris.StatusForbidden)
                }
            }
        } else {
            c.EmitError(iris.StatusUnauthorized)
        }
    })("edit-article")
    
    iris.Post("/article/:id/edit", func(c *iris.Context) {
        if isLoggedIn(c) {
            id, err := c.ParamInt64("id")
            text := c.FormValueString("text")
            title := c.FormValueString("title")
            action := c.FormValueString("action")
            var draft = false
            if c.FormValueString("draft") == "draft" {
                draft = true
            }
            var pub = false
            if c.FormValueString("privacy") == "public" {
                pub = true
            }
            if !pe(err) {
                c.EmitError(iris.StatusNotFound)
            } else {
                if action == "submit" {
                    if getUserID(c) == db.GetArticleUserID(id) {
                        if db.EditArticle(id, title, text, pub, draft) {
                            msg := "Article updated successfully!"
                            if draft {
                                msg = "Draft updated successfully!"
                            }
                            c.SetFlash("message", msg)
                            c.RedirectTo("article", id)
                        } else {
                            c.EmitError(iris.StatusInternalServerError)
                        }
                    } else {
                        c.EmitError(iris.StatusForbidden)
                    }
                } else if action == "preview" {
                    c.SetFlash("text", text)
                    c.RedirectTo("preview")
                } else if action == "delete" {
                    if !db.DeleteArticle(id) {
                        c.EmitError(iris.StatusInternalServerError)
                    } else {
                        c.SetFlash("message", "Successfully deleted article!")
                        c.SetFlash("messageType", "success")
                        c.RedirectTo("landing")
                    }
                } else {
                    c.EmitError(iris.StatusNotFound)
                }
            }
        } else {
            c.EmitError(iris.StatusUnauthorized)
        }
    })
}