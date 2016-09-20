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
    stmt, err := db.Prepare("insert into article(created_at, title, text, user_id) values(?, ?, ?, ?)")
    a := pe(err)
    _, err1 := stmt.Exec(getDatetime(), title, body, GetUserID(username))
    b := pe(err1)
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

// GetLatestArticles returns a number of recent articles
func GetLatestArticles() []Article {
    var articles []Article
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article order by created_at desc")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query()
    pe(err2)
    defer rows.Close()
    for rows.Next() {
        var art Article
        var text string
        rows.Scan(
            &art.ID,
            &art.Title,
            &text,
            &art.CreatedAt,
            &art.UserID,
        )
        art.Text = template.HTML(common.GetMarkdown(text))
        err3 := db.QueryRow("select username from user where id = ?").Scan(&art.Username)
        pe(err3)
        articles = append(articles, art)
    }
    return articles
}

var pe = common.Pe