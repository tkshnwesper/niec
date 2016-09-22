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
    "bytes"
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
    stmt, err := db.Prepare("insert into user(email, password, username, dp, bio, created_at, website, verified, verifyhash) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
    a := pe(err)
    hash := common.GenerateRandomString(username + email)
    res, err1 := stmt.Exec(email, hashedPassword, username, dp, bio, getDatetime(), website, false, hash)
    b := pe(err1)
    id, _ := res.LastInsertId()
    var buf bytes.Buffer
    tmpl, err2 := template.New("mail").Parse(common.ReadMD("verify.mail.md"))
    pe(err2)
    err3 := tmpl.Execute(&buf, struct {
        Username, Hash string
        ID int64
    }{
        username,
        hash,
        id,
    })
    pe(err3)
    mail := common.GetMarkdown(buf.String())
    common.MailService.Send("Welcome to Niec!", mail, email)
    return a && b
}

// VerifyEmail verifies the email of the user
func VerifyEmail(id int64, hash string) bool {
    var h string
    err := db.QueryRow("select verifyhash from user where id = ?", id).Scan(&h)
    if pe(err) && h == hash {
        stmt, err2 := db.Prepare("update user set verified = true where id = ?")
        pe(err2)
        _, err3 := stmt.Exec(id)
        pe(err3)
        return true
    }
    return false
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

// getArticlesFromRows returns complete Article objects on a rows input
// it fills in the username
func getArticlesFromRows(rows *sql.Rows) []Article {
    var articles []Article
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
        stmt2, err3 := db.Prepare("select username from user where id = ?")
        pe(err3)
        stmt2.QueryRow(art.UserID).Scan(&art.Username)
        articles = append(articles, art)
    }
    return articles
}

// GetLatestArticles returns a number of recent articles
func GetLatestArticles() []Article {
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article order by created_at desc")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query()
    pe(err2)
    defer rows.Close()
    return getArticlesFromRows(rows)
}

// GetArticle returns the article with the specified id
func GetArticle(id int64) Article {
    var art Article
    var text string
    err := db.QueryRow("select id, title, text, created_at, user_id from article where id = ?", id).Scan(
        &art.ID,
        &art.Title,
        &text,
        &art.CreatedAt,
        &art.UserID,
    )
    art.Text = template.HTML(common.GetMarkdown(text))
    pe(err)
    err2 := db.QueryRow("select username from user where id = ?", art.UserID).Scan(&art.Username)
    pe(err2)
    return art
}

// GetUser returns a user object of the specified username
func GetUser(id int64) User {
    var text string
    var user User
    user.ID = id
    err := db.QueryRow("select username, bio, dp, created_at, website from user where id = ?", id).Scan(
        &user.Username,
        &text,
        &user.DP,
        &user.CreatedAt,
        &user.Website,
    )
    pe(err)
    user.Bio = template.HTML(common.GetMarkdown(text))
    return user
}

// SearchArticles searches in the database for articles that match
// the passed query, and return article objects.
func SearchArticles(query string) []Article {
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article where title like \"%" + query + "%\" or text like \"%" + query + "%\"")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query()
    pe(err2)
    defer rows.Close()
    return getArticlesFromRows(rows)
}

var pe = common.Pe