package router

import (
    "github.com/kataras/iris"
    HTML "github.com/iris-contrib/template/html"
    "github.com/dchest/captcha"
    "html/template"
    "niec/common"
    "niec/db"
    "strconv"
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
    
    iris.Get("/", func(c *iris.Context) {
        msg, _ := c.GetFlash("message")
        typ, _ := c.GetFlash("messageType")
        if !isLoggedIn(c) {
            formPage := c.FormValueString("page")
            var page int
            var err error
            if formPage == "" {
                page = 1
            } else {
                page, err = strconv.Atoi(formPage)
            }
            if !pe(err) || page < 0 {
                c.RedirectTo("landing")
            } else {
                pages := common.Pagination(page, common.PaginationWindow, common.ArticlesPerPage, db.GetArticleCount())
                c.Render("home.html", struct {
                    Title string
                    Articles []db.Article
                    Page int
                    Pages []int
                }{
                    "Niec :: Home",
                    db.GetLatestArticles(page),
                    page,
                    pages,
                })
            }
        } else {
            c.Render("index.html", struct {
                Title string
                Message string
                MessageType string
            }{
                "Welcome to Niec!",
                msg,
                typ,
            })
        }
    })("landing")
    
    iris.Get("/learn-more", func(c *iris.Context) {
        c.Render("markdown.html", struct {
            Title string
            Text template.HTML
        } {
            "Learn more",
            template.HTML(common.GetMarkdown(common.ReadMD("learn.more.md"))),
        })
    })("learn-more")
    
    iris.Get("/article/:id", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            msg, _ := c.GetFlash("message")
            art := db.GetArticle(id)
            c.Render("article.html", struct {
                Title, Message string
                Article db.Article
            }{
                art.Title,
                msg,
                art,
            })
        }
    })("article")
    
    iris.Get("/user/:id", func(c *iris.Context) {
       id, err := c.ParamInt64("id")
        if err != nil {
            c.EmitError(iris.StatusNotFound)
        } else {
            user := db.GetUser(id)
            c.Render("user.html", struct {
                Title string
                User db.User
            }{
                user.Username,
                user,
            })
        }
    })("user")
    
    iris.Get("/search", func(c *iris.Context) {
        c.Render("search.html", struct {
            Title string
            Articles []db.Article
        }{
            "Niec :: Search",
            db.SearchArticles(c.FormValueString("query")),
        })
    })("search")
    
    iris.Get("/logout", func(c *iris.Context) {
        c.Session().Clear()
        c.RedirectTo("landing")
    })("logout")
    
    var capHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
    iris.Get("/captcha/*id", iris.ToHandlerFunc(capHandler))("captcha")
    
}

func isLoggedIn(c *iris.Context) bool {
    return c.Session().GetString(common.UserIdentificationAttribute) != ""
}

var pe = common.Pe