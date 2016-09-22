/*
Package db offers a range of functions to manipulate the database
*/
package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"  // required by sql package to support mysql
    "fmt"
    "niec/common"
    "time"
    "html/template"
)

// Article structure is used for storing data of an article
type Article struct {
    ID int64
    Title string
    Text template.HTML
    CreatedAt string
    UserID int64
    Username string
}

// User structure holds non-sensitive information about an user. Sensitive information to be
// extracted manually using specific functions.
type User struct {
    ID int64
    Username string
    Bio template.HTML
    DP string
    CreatedAt string
    Website string
}

var db *sql.DB

// Init Initializes the database
func Init() {
    config := common.ConfigObject.DB
    var (
        dbName = config.Name
        dbUser = config.User
        dbPass = config.Password
    )
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", dbUser, dbPass, dbName))
	pe(err)
}

// GetDatetime returns a mysql compatible datetime in order to store it in db
func getDatetime() string {
    loc, _ := time.LoadLocation("Asia/Calcutta")
    t, _ := time.ParseInLocation(
        "2006 Jan 02 15:04:05",
        time.Now().Format("2006 Jan 02 15:04:05"),
        loc,
    )
    return fmt.Sprintf(
        "%d-%02d-%02d %02d:%02d:%02d",
        t.Year(),
        t.Month(),
        t.Day(),
        t.Hour(),
        t.Minute(),
        t.Second(),
    )
}

var pe = common.Pe