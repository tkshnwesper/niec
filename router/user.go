package router

import (
    "github.com/kataras/iris"
    "niec/common"
    "niec/db"
)

// InitUserPages initializes the user pages
func InitUserPages() {
    iris.Get("/user/:id", func(c *iris.Context) {
       id, err := c.ParamInt64("id")
       msg, _ := c.GetFlash("message")
       typ, _ := c.GetFlash("messageType")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            user, ok := db.GetUser(id, isLoggedIn(c))
            if !ok {
                c.EmitError(iris.StatusNotFound)
            } else {
                c.Render("user.html", struct {
                    Title string
                    Property Property
                    User db.User
                    Message, MessageType string
                }{
                    user.Username,
                    getProperty(c),
                    user,
                    msg, typ,
                })
            }
        }
    })("user")
    
    iris.Get("/user/:id/article", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            formPage := c.FormValueString("page")
            count := db.GetUserArticleCount(id, isLoggedIn(c))
            page, b := common.ValidPagination(formPage, count, common.ArticlesPerPage)
            if !b {
                c.RedirectTo("user-article", id)
            } else {
                pages := common.Pagination(page, common.PaginationWindow, common.ArticlesPerPage, count)
                c.Render("user.article.html", struct {
                    Title string
                    Property Property
                    Articles []db.Article
                    Page int
                    Pages []int
                    Increment func(int)(int)
                    Decrement func(int)(int)
                    Path, URL string
                }{
                    db.GetUsernameFromID(id) + "'s Articles",
                    getProperty(c),
                    db.GetUserArticles(id, isLoggedIn(c), page),
                    page,
                    pages,
                    common.Increment,
                    common.Decrement,
                    iris.URL("user-article", id),
                    "",
                })
            }
        }
    })("user-article")
    
    iris.Get("/user/:id/draft", func(c *iris.Context) {
       id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else if !isLoggedIn(c) {
            c.EmitError(iris.StatusUnauthorized);
        } else if id != getUserID(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            count := db.GetDraftCount(id)
            formPage := c.FormValueString("page")
            page, b := common.ValidPagination(formPage, count, common.ArticlesPerPage)
            if !b {
                c.RedirectTo("draft", id)
            } else {
                pages := common.Pagination(page, common.PaginationWindow, common.ArticlesPerPage, count)
                c.Render("draft.html", struct {
                    Title string
                    Property Property
                    Articles []db.Article
                    Page int
                    Pages []int
                    Increment func(int)(int)
                    Decrement func(int)(int)
                    Path, URL string
                }{
                    getUsername(c),
                    getProperty(c),
                    db.GetDraftList(getUserID(c), page),
                    page,
                    pages,
                    common.Increment,
                    common.Decrement,
                    iris.URL("draft", id),
                    "",
                })
            }
        }
    })("draft")
    
    iris.Get("/user/:id/edit", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else if !isLoggedIn(c) {
            c.EmitError(iris.StatusUnauthorized)
        } else if id != getUserID(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            dp, website, bio, public := db.FetchForEditProfile(id)
            fields := []Field {
                {
                    "url",
                    "dp",
                    "Display Picture URL",
                    255,
                    dp,
                    false,
                },
                {
                    "url",
                    "website",
                    "Website",
                    255,
                    website,
                    false,
                },
            }
            cb := []Checkbox {
                {
                    "Public profile",
                    "privacy",
                    "public",
                    "If your profile is public, it can be viewed without logging in",
                    "globe",
                    public,
                },
            }
            buttons := []Button {
                {
                    "submit",
                    "submit",
                    "Submit",
                },
                {
                    "reset",
                    "reset",
                    "Reset",
                },
            }
            c.Render("edit.profile.html", struct {
                Title string
                Property Property
                Checkboxes []Checkbox
                Fields []Field
                Bio string
                Buttons []Button
            }{
                "Edit Profile",
                getProperty(c),
                cb,
                fields,
                bio,
                buttons,
            })
        }
    })("edit-profile")
    
    iris.Post("/user/:id/edit", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else if !isLoggedIn(c) {
            c.EmitError(iris.StatusUnauthorized)
        } else if id != getUserID(c) {
            c.EmitError(iris.StatusForbidden)
        } else {
            dp := c.FormValueString("dp")
            website := c.FormValueString("website")
            bio := c.FormValueString("bio")
            var pub = false
            if c.FormValueString("privacy") == "public" {
                pub = true
            }
            if !db.EditProfile(id, dp, website, bio, pub) {
                c.EmitError(iris.StatusInternalServerError)
            } else {
                c.SetFlash("message", "Profile updated successfully!")
                c.SetFlash("messageType", "success")
                c.RedirectTo("user", id)
            }
        }
    })
}