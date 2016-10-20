package router

import (
    "github.com/kataras/iris"
    "niec/db"
)

// InitArticlePages initializes the article pages
func InitArticlePages() {
    iris.Get("/article/:id", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            msg, _ := c.GetFlash("message")
            art, ok := db.GetArticle(isLoggedIn(c), getUserID(c), id)
            if !ok {
                c.EmitError(iris.StatusNotFound)
            } else {
                var arts []db.Article
                c.Render("article.html", struct {
                    Title, Message, MessageType string
                    Property Property
                    Articles []db.Article
                }{
                    art.Title,
                    msg,
                    "success",
                    getProperty(c),
                    append(arts, art),
                })
            }
        }
    })("article")
    
    iris.Get("/article/:id/raw", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            text, ok := db.GetRaw(isLoggedIn(c), getUserID(c), id)
            if !ok {
                c.EmitError(iris.StatusNotFound)
            } else {
                c.Text(iris.StatusOK, text)
            }
        }
    })("raw-article")
    
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