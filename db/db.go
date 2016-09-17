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
    website string,
) bool {
    hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
    stmt, err := db.Prepare("insert into user(email, password, username, dp, bio, created_at, website) values(?, ?, ?, ?, ?, ?, ?)")
    a := pe(err)
    _, err1 := stmt.Exec(email, hashedPassword, username, dp, bio, getDatetime(), website)
    b := pe(err1)
    return a && b
}

// InsertArticle inserts an article into the database
func InsertArticle(username, title, _, body string) bool {
    stmt, err := db.Prepare("insert into article(created_at, title, text) values(?, ?, ?)")
    a := pe(err)
    res, err1 := stmt.Exec(getDatetime(), title, body)
    b := pe(err1)
    if !(a && b) {
        return false
    }
    lid, _ := res.LastInsertId()
    uid := GetUserID(username)
    stmt, err = db.Prepare("insert into map_user_article(user_id, article_id) values(?, ?)")
    a = pe(err)
    res, err1 = stmt.Exec(uid, lid)
    b = pe(err1)
    return a && b
}

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

// GetUserID returns the userID of the user whose username is passed to it
func GetUserID(username string) int64 {
    var id int64
    err := db.QueryRow("select id from user where username = ?", username).Scan(&id)
    pe(err)
    return id
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