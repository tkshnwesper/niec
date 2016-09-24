package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"  // required by sql package to support mysql
    "niec/common"
    "html/template"
)

// InsertArticle inserts an article into the database
func InsertArticle(username, title, _, body string) (int64, bool) {
    stmt, err := db.Prepare("insert into article(created_at, title, text, user_id) values(?, ?, ?, ?)")
    a := pe(err)
    res, err1 := stmt.Exec(getDatetime(), title, body, GetUserID(username))
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
    limit := page * common.ArticlesPerPage
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article order by created_at desc limit ? offset ?")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query(limit, offset)
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
    art.Username = GetUsernameFromID(art.UserID)
    return art
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

// GetArticleUserID returns the user ID of the creator of that article
func GetArticleUserID(id int64) int64 {
    var uid int64
    err := db.QueryRow("select user_id from article where id = ?", id).Scan(&uid)
    pe(err)
    return uid
}

// FetchForEdit fetches the data required by the edit page
func FetchForEdit(id int64) (string, string) {
    var title, text string
    err := db.QueryRow("select title, text from article where id = ?", id).Scan(&title, &text)
    pe(err)
    return title, text
}

// EditArticle writes the edits back to the database
func EditArticle(id int64, title, text string) bool {
    _, err := db.Exec("update article set title = ?, text = ? where id = ?", title, text, id)
    return pe(err)
}

func GetArticleCount() int64 {
    var num int64
    err := db.QueryRow("select count(*) from article where 1").Scan(&num)
    pe(err)
    return num
}