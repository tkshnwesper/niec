// Package common includes functions that can cater to all the different
// needs of the entire project. It is a standalone library which does not
// require the other packages in this project
package common

import (
    "log"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
    "io/ioutil"
)

/* UserIdentificationAttribute will be stored in the server session in order to identify the user
It should be something that uniquely identifies the user, for example: user_id, username or email */
const UserIdentificationAttribute = "username"

// Pe returns whether an error is real or not and prints fatal output if it is
func Pe(err error) bool {
    if err != nil {
        log.Fatal(err)
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