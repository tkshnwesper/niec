package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"  // required by sql package to support mysql
    "niec/common"
    "html/template"
)

// InsertArticle inserts an article into the database
func InsertArticle(uid int64, title, body string, pub bool) (int64, bool) {
    stmt, err := db.Prepare("insert into article(created_at, title, text, user_id, public) values(?, ?, ?, ?, ?)")
    a := pe(err)
    res, err1 := stmt.Exec(getDatetime(), title, body, uid, pub)
    b := pe(err1)
    id, _ := res.LastInsertId()
    return id, a && b
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
        art.Username = GetUsernameFromID(art.UserID)
        articles = append(articles, art)
    }
    return articles
}

// GetLatestArticles returns a number of recent articles
func GetLatestArticles(page int) []Article {
    offset := (page - 1) * common.ArticlesPerPage
    limit := common.ArticlesPerPage
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article order by created_at desc limit ? offset ?")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query(limit, offset)
    pe(err2)
    defer rows.Close()
    return getArticlesFromRows(rows)
}

// GetArticle returns the article with the specified id
func GetArticle(loggedin bool, id int64) (Article, bool) {
    var art Article
    var text string
    var mid = " "
    if !loggedin {
        mid = " public = true and "
    }
    err := db.QueryRow("select id, title, text, created_at, user_id from article where" + mid + "id = ?", id).Scan(
        &art.ID,
        &art.Title,
        &text,
        &art.CreatedAt,
        &art.UserID,
    )
    art.Text = template.HTML(common.GetMarkdown(text))
    if pe(err) {
        return Article {}, false
    art.Username = GetUsernameFromID(art.UserID)
    return art, true
}

// SearchArticles searches in the database for articles that match
// the passed query, and return article objects.
func SearchArticles(loggedin bool, page int, query string) []Article {
    var mid = " "   // space is intentional
    if !loggedin {
        mid = " public = true and "
    }
    offset := (page - 1) * common.ArticlesPerPage
    limit := common.ArticlesPerPage
    like := "%" + query + "%"
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article where" + mid + "(title like ? or text like ?) limit ? offset ?")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query(like, like, limit, offset)
    pe(err2)
    defer rows.Close()
    return getArticlesFromRows(rows)
}

// GetArticleUserID returns the user ID of the creator of that article
func GetArticleUserID(id int64) int64 {
    var uid int64
    err := db.QueryRow("select user_id from article where id = ?", id).Scan(&uid)
    pe(err)
    return uid
}

// FetchForEdit fetches the data required by the edit page
func FetchForEdit(id int64) (string, string, bool) {
    var title, text string
    var pub bool
    err := db.QueryRow("select title, text, public from article where id = ?", id).Scan(&title, &text, &pub)
    pe(err)
    return title, text, pub
}

// EditArticle writes the edits back to the database
func EditArticle(id int64, title, text string, pub bool) bool {
    _, err := db.Exec("update article set title = ?, text = ?, public = ? where id = ?", title, text, pub, id)
    return pe(err)
}

// GetArticleCount counts the total number of articles present in the article table
func GetArticleCount() int64 {
    var num int64
    err := db.QueryRow("select count(*) from article where 1").Scan(&num)
    pe(err)
    return num
}