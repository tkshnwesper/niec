package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"  // required by sql package to support mysql
    "niec/common"
    "html/template"
)

// InsertArticle inserts an article into the database
func InsertArticle(uid int64, title, body string, pub, draft bool) (int64, bool) {
    stmt, err := db.Prepare("insert into article(created_at, title, text, user_id, public, draft) values(?, ?, ?, ?, ?, ?)")
    a := pe(err)
    res, err1 := stmt.Exec(getDatetime(), title, body, uid, pub, draft)
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


///////////////////////////////////////////////////////////////////////////////////////////
// Similar queries in GetLatestArticles and GetArticleCount
///////////////////////////////////////////////////////////////////////////////////////////

// GetLatestArticles returns a number of recent articles
func GetLatestArticles(page int) ([]Article) {
    offset := (page - 1) * common.ArticlesPerPage
    limit := common.ArticlesPerPage
    prep := "select id, title, text, created_at, user_id from article where draft = false order by created_at desc limit ? offset ?"
    stmt, err := db.Prepare(prep)
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query(limit, offset)
    pe(err2)
    defer rows.Close()
    return getArticlesFromRows(rows)
}


// GetArticleCount counts the total number of articles present in the article table
func GetArticleCount() int64 {
    var num int64
    err := db.QueryRow("select count(*) from article where draft = false").Scan(&num)
    pe(err)
    return num
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////



///////////////////////////////////////////////////////////////////////////////////////////
// Similar queries in GetArticle and GetRaw
///////////////////////////////////////////////////////////////////////////////////////////

// GetArticle returns the article with the specified id
func GetArticle(loggedin bool, uid, id int64) (Article, bool) {
    var art Article
    var text string
    var mid = " "
    if !loggedin {
        mid = " public = true and "
    }
    err := db.QueryRow("select id, title, text, created_at, user_id from article where ((draft = false) or (draft = true and user_id = ?)) and" + mid + "id = ?", uid, id).Scan(
        &art.ID,
        &art.Title,
        &text,
        &art.CreatedAt,
        &art.UserID,
    )
    art.Text = template.HTML(common.GetMarkdown(text))
    if err != nil {
        return Article {}, false
    }
    art.Username = GetUsernameFromID(art.UserID)
    return art, true
}

// GetRaw returns the raw text
func GetRaw(loggedin bool, uid, id int64) (string, bool) {
    var text string
    var mid = " "
    if !loggedin {
        mid = " public = true and "
    }
    err := db.QueryRow(
        "select text from article where ((draft = false) or (draft = true and user_id = ?)) and" + mid + "id = ?",
        uid, id,
    ).Scan(&text)
    if err != nil {
        return "", false
    }
    return text, true
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////////////////////////////////
// Similar queries in SearchArticles and GetSearchCount
///////////////////////////////////////////////////////////////////////////////////////////

// SearchArticles searches in the database for articles that match
// the passed query, and return article objects.
func SearchArticles(loggedin bool, page int, query string) ([]Article) {
    var mid = " "   // space is intentional
    if !loggedin {
        mid = " public = true and "
    }
    offset := (page - 1) * common.ArticlesPerPage
    limit := common.ArticlesPerPage
    like := "%" + query + "%"
    stmt, err := db.Prepare("select id, title, text, created_at, user_id from article where draft = false and" + mid + "(title like ? or text like ?) order by created_at desc limit ? offset ?")
    pe(err)
    defer stmt.Close()
    rows, err2 := stmt.Query(like, like, limit, offset)
    pe(err2)    
    defer rows.Close()
    return getArticlesFromRows(rows)
}

// GetSearchCount returns the count of the number of search items
func GetSearchCount(loggedin bool, query string) int64 {
    var mid = " "   // space is intentional
    var count int64
    if !loggedin {
        mid = " public = true and "
    }
    like := "%" + query + "%"
    err3 := db.QueryRow(
        "select count(*) from article where draft = false and" + mid + "(title like ? or text like ?)",
        like, like,
    ).Scan(&count)
    pe(err3)
    return count
}

////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////

// GetArticleUserID returns the user ID of the creator of that article
func GetArticleUserID(id int64) int64 {
    var uid int64
    err := db.QueryRow("select user_id from article where id = ?", id).Scan(&uid)
    pe(err)
    return uid
}

// FetchForEdit fetches the data required by the edit page
func FetchForEdit(id int64) (string, string, bool, bool) {
    var title, text string
    var pub, draft bool
    err := db.QueryRow("select title, text, public, draft from article where id = ?", id).Scan(&title, &text, &pub, &draft)
    pe(err)
    return title, text, pub, draft
}

// EditArticle writes the edits back to the database
func EditArticle(id int64, title, text string, pub, draft bool) bool {
    _, err := db.Exec("update article set title = ?, text = ?, public = ?, draft = ? where id = ?", title, text, pub, draft, id)
    return pe(err)
}

///////////////////////////////////////////////////////////////////////////////////////////
// Similar queries in GetDraftList and GetDraftCount
///////////////////////////////////////////////////////////////////////////////////////////

// GetDraftList returns a list of drafts for a user
func GetDraftList(uid int64, page int) []Article {
    offset := (page - 1) * common.ArticlesPerPage
    limit := common.ArticlesPerPage
    rows, err := db.Query(
        "select id, title, text, created_at, user_id from article where draft = true and user_id = ? order by created_at desc limit ? offset ?",
        uid, limit, offset,
    )
    pe(err)
    return getArticlesFromRows(rows)
}

// GetDraftCount returns the number of articles saved as draft for a particular user
func GetDraftCount(uid int64) int64 {
    var count int64
    err := db.QueryRow("select count(*) from article where draft = true and user_id = ?", uid).Scan(&count)
    pe(err)
    return count
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

// DeleteArticle deletes an article from the article table
func DeleteArticle(id int64) bool {
    _, err := db.Exec("delete from article where id = ?", id)
    if err != nil {
        return false
    }
    return true
}