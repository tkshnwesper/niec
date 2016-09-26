package router

import (
    "niec/common"
    "github.com/kataras/iris"
    "html/template"
)

// Markdown is a container for markdown files
type Markdown struct {
    FileName string
}

// Serve serves markdown pages
func (m Markdown) Serve(c *iris.Context) {
    c.Render("markdown.html", struct {
        Title string
        Property Property
        Text template.HTML
    } {
        "Learn more",
        getProperty(c),
        template.HTML(common.GetMarkdown(common.ReadMD("learn.more.md"))),
    })
}

// InitMarkdownPages initializes pages rendered by markdown
func InitMarkdownPages() {
    pages := map[string]Markdown {
        "learn-more": Markdown {
            "learn-more.html",
        },
    }
    for s, mk := range pages {
        iris.Handle("GET", "/" + s, mk)(s)
    }
}