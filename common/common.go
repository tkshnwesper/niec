// Package common includes functions that can cater to all the different
// needs of the entire project. It is a standalone library which does not
// require the other packages in this project
package common

import (
    "log"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
    "io/ioutil"
    // "github.com/iris-contrib/mail"
    mail "github.com/kataras/go-mailer"
    // "net/smtp"
    "math/rand"
    "crypto/md5"
    "time"
    "fmt"
    "strconv"
    "regexp"
)

// PaginationWindow is the number of numbers that you will be able to see at any point in time
// in the pagination
const PaginationWindow = 5

// ArticlesPerPage is the number of articles that are displayed per page.
const ArticlesPerPage = 10

// ConfigObject will be used by modules that require the config
var ConfigObject Config

// MailService provides feature to send mail
var MailService mail.Service

// Authentication stores Auth obj
// var Authentication smtp.Auth

// MailData stores mail data
// var MailData ConfigSMTP

// Pe returns whether an error is real or not and prints output if it is
// false if there is error
// true if there is no error
func Pe(err error) bool {
    if err != nil {
        log.Print(err)
        return false
    }
    return true
}

// GetMarkdown returns the HTML translation of markdown code
func GetMarkdown(s string) string {
    p := bluemonday.NewPolicy()
    // Permits the "dir", "id", "lang", "title" attributes globally
    p.AllowStandardAttributes()

    // Permits the "img" element and it's standard attributes
    p.AllowImages()

    // Permits ordered and unordered lists, and also definition lists
    p.AllowLists()

    // Permits HTML tables and all applicable elements and non-styling attributes
    p.AllowTables()
    p.AllowStandardURLs()
    p.AddTargetBlankToFullyQualifiedLinks(true)
    p.AllowDataURIImages()
    p.RequireNoFollowOnLinks(true)
    p.AllowRelativeURLs(true)
    p.AllowAttrs("href").OnElements("a")
    p.AllowElements("pre", "blockquote", "p", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "strong", "em", "del")
    p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
    return p.Sanitize(string(blackfriday.MarkdownCommon([]byte(s))))
}

// ReadMD reads and returns text from the markdown directory
func ReadMD(name string) string {
    dat, err := ioutil.ReadFile("markdown/" + name)
    if !Pe(err) {
        return ""
    }
    return string(dat)
}

// Init initializes the ConfigObject and Mail
func Init() {
    ConfigObject = parse()
    initMail()
}

func initMail() {
    // MailData = ConfigObject.SMTP
    // data := MailData
    smtp := ConfigObject.SMTP
    var config = mail.Config {
        Host: smtp.Host,
        Username: smtp.Username,
        FromAddr: smtp.Username + "@" + smtp.Host,
        Password: smtp.Password,
        Port: smtp.Port,
        FromAlias: smtp.FromAlias,
    }
    MailService = mail.New(config)
    // Authentication = smtp.PlainAuth("", data.Username, data.Password, data.Host)
}

// SendMail sends mail
// func SendMail(to []string, msg string) {
//     smtp.SendMail(MailData.Host+":25", Authentication, MailData.Username+"@"+MailData.Host, to, []byte(msg))
// }

// StrShuffle shuffles a string that is passed to it
func StrShuffle(str string) string {
    n := 25 + rand.Intn(25)
    var runes []rune
    for _, runeval := range str {  
        runes = append(runes, runeval)
    }
    var shuffle = func() bool {
        if rand.Intn(2) == 1 {
            return true
        }
        return false
    }
    var repeat = func() {
        for i := 0; i < len(runes) - 1; i++ {
            if shuffle() {
                temp := runes[i]
                runes[i] = runes[i + 1]
                runes[i + 1] = temp
            }
        }
        n--
    }
    for n > 0 {
        repeat()
    }
    return string(runes)
}

// GenerateRandomString generates a random string from a hash
func GenerateRandomString(salt string) string {
    m := fmt.Sprintf("%x", md5.Sum([]byte(salt)))
    return StrShuffle(m + fmt.Sprintf("%d", time.Now().Nanosecond()))
}

// Pagination returns an array of integers which will be the page's current pagination
func Pagination(page, paginationWindow, articlesPerPage int, maxart int64) []int {
    if page == 0 {
        page++
    }
    totpage := int(maxart / int64(articlesPerPage))
    if maxart % int64(articlesPerPage) != 0 {
        totpage++
    }
    minthresh := paginationWindow / 2
    maxthresh := totpage - minthresh
    var pages []int
    if page > minthresh && page < maxthresh {
        for i := page - minthresh; i <= page + minthresh; i++ {
            pages = append(pages, i)
        }
    } else if page < minthresh {
        var max int
        if totpage - paginationWindow > 0 {
            max = paginationWindow
        } else {
            max = totpage
        }
        for i := 1; i <= max; i++ {
            pages = append(pages, i)
        }
    } else {
        start := totpage - paginationWindow
        if start < 0 {
            start = 1
        }
        for i := start; i <= totpage; i++ {
            pages = append(pages, i)
        }
    }
    return pages
}

// ValidPagination tells whether the page passed is valid
func ValidPagination(formPage string, count0 int64, articlesPerPage int) (int, bool) {
    var page int
    var err error
    if formPage == "" {
        page = 1
    } else {
        page, err = strconv.Atoi(formPage)
    }
    app, count := int64(articlesPerPage), count0
    numpages := count / app
    if count % app != 0 || count == 0 {
        numpages++
    }
    return page, !(!Pe(err) || page < 1 || int64(page) > numpages)
}

// Increment adds one
func Increment(num int) int {
    return num + 1
}

// Decrement subtracts one
func Decrement(num int) int {
    return num - 1
}

// FormatTime returns a formatted string of datetime
func FormatTime(t string) string {
    pt, err := time.Parse("2006-01-02 15:04:05", t)
    if Pe(err) {
        return pt.Format("Monday, _2 January 2006 15:04:05")
    }
    return ""
}