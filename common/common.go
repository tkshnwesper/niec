// Package common includes functions that can cater to all the different
// needs of the entire project. It is a standalone library which does not
// require the other packages in this project
package common

import (
    "log"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
    "io/ioutil"
    "github.com/iris-contrib/mail"
    "math/rand"
    "crypto/md5"
    "time"
    "fmt"
)

// UserIdentificationAttribute will be stored in the server session in order to identify the user
// It should be something that uniquely identifies the user, for example: user_id, username or email
const UserIdentificationAttribute = "username"

// ConfigObject will be used by modules that require the config
var ConfigObject Config

// MailService provides feature to send mail
var MailService mail.Service

// Pe returns whether an error is real or not and prints fatal output if it is
func Pe(err error) bool {
    if err != nil {
        log.Print(err)
        return false
    }
    return true
}

// GetMarkdown returns the HTML translation of markdown code
func GetMarkdown(s string) string {
    return bluemonday.UGCPolicy().Sanitize(string(blackfriday.MarkdownCommon([]byte(s))))
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
    var smtp = ConfigObject.SMTP
    var config = mail.Config {
        Host: smtp.Host,
        Username: smtp.Username,
        Password: smtp.Password,
        Port: smtp.Port,
        FromAlias: smtp.FromAlias,
    }
    MailService = mail.New(config)
}

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