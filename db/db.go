/*
Package db offers a range of functions to manipulate the database
*/
package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"  // required by sql package to support mysql
    "fmt"
    "crypto/md5"
    "niec/common"
    "time"
)

var db *sql.DB

// Init Initializes the database
func Init() {
    config := parse()
    var (
        dbName = config.DB.Name
        dbUser = config.DB.User
        dbPass = config.DB.Password
    )
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", dbUser, dbPass, dbName))
	pe(err)
}

// CheckEmailExists Checks whether a user with the given email is present in the database
func CheckEmailExists(em string) bool {
    var n int
    err := db.QueryRow("select count(*) from user where email = ?", em).Scan(&n)
    pe(err)
    if n == 0 {
        return false
    }
    return true
}

// CheckUsernameExists Checks whether a user with the give username exists in the database
func CheckUsernameExists(em string) bool {
    var n int
    err := db.QueryRow("select count(*) from user where username = ?", em).Scan(&n)
    pe(err)
    if n == 0 {
        return false
    }
    return true
}

// InsertUser inserts into the table user a new user
func InsertUser(
    email string,
    password string,
    username string,
    dp string,
    bio string,
) bool {
    hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
    stmt, err := db.Prepare("insert into user(email, password, username, dp, bio, created_at) values(?, ?, ?, ?, ?, ?)")
    a := pe(err)
    _, err1 := stmt.Exec(email, hashedPassword, username, dp, bio, getDatetime())
    b := pe(err1)
    return a && b
}

// func InsertArticle

// VerifyCreds verifies whether the email and password match
func VerifyCreds(email string, password string) bool {
    var pass string
    err := db.QueryRow("select password from user where email = ?", email).Scan(&pass)
    pe(err)
    if fmt.Sprintf("%x", md5.Sum([]byte(password))) == pass {
        return true
    }
    return false
}

// GetUsername returns the username of a user when email is passed to it
func GetUsername(email string) string {
    var un string
    err := db.QueryRow("select username from user where email = ?", email).Scan(&un)
    pe(err)
    return un
}

// GetDatetime returns a mysql compatible datetime in order to store it in db
func getDatetime() string {
    t := time.Now()
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