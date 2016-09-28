package router

import (
    "niec/common"
    "github.com/kataras/iris"
    "html/template"
)

// Markdown is a container for markdown files
type Markdown struct {
    Path, FileName, Title string
    LoggedIn bool
}

// Serve serves markdown pages
func (m Markdown) Serve(c *iris.Context) {
    ok := true
    if m.LoggedIn && !isLoggedIn(c) {
        ok = false
    }
    if ok {
        c.Render("markdown.html", struct {
            Title string
            Property Property
            Text template.HTML
        } {
            m.Title,
            getProperty(c),
            template.HTML(common.GetMarkdown(common.ReadMD(m.Path + m.FileName))),
        })
    } else {
        c.EmitError(iris.StatusUnauthorized)
    }
}

// InitMarkdownPages initializes pages rendered by markdown
func InitMarkdownPages() {
    pages := map[string]Markdown {
        "learn-more": Markdown {
            "",
            "learn.more.md",
            "Learn More",
            false,
        },
    }
    for s, mk := range pages {
        iris.Handle("GET", "/" + mk.Path + s, mk)(s)
    }
}