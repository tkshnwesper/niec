package router

import (
    "github.com/kataras/iris"
    HTML "github.com/iris-contrib/template/html"
    "github.com/dchest/captcha"
    "niec/common"
    "niec/db"
)

// Button holds information about the buttons to be displayed in the view
type Button struct {
    Type string
    Name string
    Placeholder string
}

// Field holds information about the input fields to be displayed in the view
type Field struct {
    Type string
    Name string
    Placeholder string
    MaxSize int
    Value string
    SROnly bool
}

// Checkbox holds information about checkboxes
type Checkbox struct {
    Title, Name, Value, Tooltip, Glyph string
    Checked bool
}

// Init helps to initialize all the pages required in the site
func Init() {
    iris.UseTemplate(HTML.New(HTML.Config {
        Layout: "layout.html",
    }))
    
    iris.StaticServe("./static/", "static")
    
    initErrorPages()
    
    initSignPages()
    
    initSubmitPages()
    
    InitMarkdownPages()
    
    InitUserPages()
    
    InitArticlePages()
    
    iris.Get("/", func(c *iris.Context) {
        msg, _ := c.GetFlash("message")
        typ, _ := c.GetFlash("messageType")
        // if isLoggedIn(c) {
        formPage := c.FormValueString("page")
        count := db.GetArticleCount(isLoggedIn(c))
        page, b := common.ValidPagination(formPage, count, common.ArticlesPerPage)
        if !b {
            c.RedirectTo("landing")
        } else {
            pages := common.Pagination(page, common.PaginationWindow, common.ArticlesPerPage, count)
            c.Render("home.html", struct {
                Title string
                Property Property
                Articles []db.Article
                Page int
                Pages []int
                Increment func(int)(int)
                Decrement func(int)(int)
                Path, URL, Message, MessageType string
            }{
                "Niec :: Home",
                getProperty(c),
                db.GetLatestArticles(page, isLoggedIn(c)),
                page,
                pages,
                common.Increment,
                common.Decrement,
                "landing",
                "",
                msg,
                typ,
            })
        }
        // } else {
        //     c.Render("index.html", struct {
        //         Title string
        //         Property Property
        //         Message string
        //         MessageType string
        //     }{
        //         "Welcome to Niec!",
        //         getProperty(c),
        //         msg,
        //         typ,
        //     })
        // }
    })("landing")
    
    iris.Get("/search", func(c *iris.Context) {
        formPage := c.FormValueString("page")
        query := c.FormValueString("query")
        count := db.GetSearchCount(isLoggedIn(c), query)
        page, b := common.ValidPagination(formPage, count, common.ArticlesPerPage)
        if !b {
            c.RedirectTo("search")
        } else {
            pages := common.Pagination(page, common.PaginationWindow, common.ArticlesPerPage, count)
            c.Render("search.html", struct {
                Title string
                Property Property
                Articles []db.Article
                Page int
                Pages []int
                Increment func(int)(int)
                Decrement func(int)(int)
                Path, URL string
            }{
                "Niec :: Search",
                getProperty(c),
                db.SearchArticles(isLoggedIn(c), page, query),
                page,
                pages,
                common.Increment,
                common.Decrement,
                "search",
                "query=" + query + "&",
            })
        }
    })("search")
    
    iris.Get("/logout", func(c *iris.Context) {
        c.Session().Clear()
        c.RedirectTo("landing")
    })("logout")
    
    var capHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
    iris.Get("/captcha/*id", iris.ToHandlerFunc(capHandler))("captcha")
    
}

func isLoggedIn(c *iris.Context) bool {
    p := c.Session().Get("property")
    if p != nil {
        return p.(Property).LoggedIn
    }
    return false
}

func getUsername(c *iris.Context) string {
    p := c.Session().Get("property")
    if p != nil {
        return p.(Property).Username
    }
    return ""
}

func getUserID(c *iris.Context) int64 {
    p := c.Session().Get("property")
    if p != nil {
        return p.(Property).UserID
    }
    return 0
}

func getProperty(c *iris.Context) Property {
    p := c.Session().Get("property")
    if p != nil {
        return p.(Property)
    }
    return Property {
        LoggedIn: false,
    }
}

var pe = common.Pe