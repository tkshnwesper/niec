package router

import (
    "niec/common"
    "github.com/kataras/iris"
    "html/template"
    "io/ioutil"
    "strings"
)

// Markdown is a container for markdown files
type Markdown struct {
    Path, FileName, Title, Link string
    Public bool
}

// Serve serves markdown pages
func (m Markdown) Serve(c *iris.Context) {
    ok := true
    var pub = "public"
    if !m.Public {
        pub = "private"
    }
    pub = pub + "/"
    if !m.Public && !isLoggedIn(c) {
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
            template.HTML(common.GetMarkdown(common.ReadMD(pub + m.Path + m.FileName))),
        })
    } else {
        c.EmitError(iris.StatusUnauthorized)
    }
}

func renderDir(n, pre string, pub bool) {
    dirlist, err := ioutil.ReadDir("markdown/" + n + "/" + pre)
    common.Pe(err)
    pages := make(map[string]Markdown)
    for _, i := range dirlist {
        name := i.Name()
        if name[len(name) - 2:] == "md" {
            pages[name] = Markdown {
                pre, name, name[: len(name) - 3], strings.Join(strings.Split(name, " "), ""),  pub,
            }
        } else {
            var exp string
            if pre == "" {
                exp = name + "/"
            } else {
                exp = pre + "/" + name + "/"
            }
            renderDir(n, exp, pub)
        }
    }
    for _, mk := range pages {
        iris.Handle("GET", "/" + mk.Path + mk.Link, mk)(mk.Link)
    }
}

// InitMarkdownPages initializes pages rendered by markdown
func InitMarkdownPages() {
    renderDir("public", "", true)
    renderDir("private", "", false)
}