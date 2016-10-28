package db

import (
    "crypto/md5"
    "niec/common"
    "fmt"
    "bytes"
    "html/template"
)

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
    public bool,
) bool {
    hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
    stmt, err := db.Prepare("insert into user(email, password, username, dp, bio, created_at, website, verified, verifyhash, public) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
    a := pe(err)
    hash := common.GenerateRandomString(username + email)
    res, err1 := stmt.Exec(email, hashedPassword, username, dp, bio, getDatetime(), website, false, hash, public)
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
    // common.SendMail([]string{email}, mail)
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

// GetUser returns a user object of the specified username
func GetUser(id int64, loggedin bool) (User, bool) {
    var text string
    var user User
    var mid = " "
    if !loggedin {
        mid = " public = true and "
    }
    user.ID = id
    err := db.QueryRow("select username, bio, dp, created_at, website from user where" + mid + "id = ?", id).Scan(
        &user.Username,
        &text,
        &user.DP,
        &user.CreatedAt,
        &user.Website,
    )
    if err != nil {
        return User {}, false
    }
    user.Bio = template.HTML(common.GetMarkdown(text))
    return user, true
}

// VerifyCreds verifies whether the email and password match
func VerifyCreds(email string, password string) (bool, bool) {
    var pass string
    var verified bool
    err := db.QueryRow("select password, verified from user where email = ?", email).Scan(&pass, &verified)
    pe(err)
    if fmt.Sprintf("%x", md5.Sum([]byte(password))) == pass && verified {
        return true, true
    } else if !verified {
        return true, false
    }
    return false, false
}

// GetUsernameFromID returns the username of a user from their ID
func GetUsernameFromID(id int64) string {
    var un string
    err3 := db.QueryRow("select username from user where id = ?", id).Scan(&un)
    pe(err3)
    return un
}

// GetUsernameAndID returns both username and ID from user table based on email
func GetUsernameAndID(email string) (string, int64) {
    var un string
    var id int64
    err := db.QueryRow("select username, id from user where email = ?", email).Scan(&un, &id)
    pe(err)
    return un, id
}

// FetchForEditProfile returns fields that can be changed in the edit profile page
func FetchForEditProfile(id int64) (string, string, string, bool) {
    var dp, website, bio string
    var public bool
    err := db.QueryRow("select dp, website, bio, public from user where id = ?", id).Scan(
        &dp, &website, &bio, &public,
    )
    if err != nil {}
    return dp, website, bio, public
}

// EditProfile commits changes to the profile
func EditProfile(id int64, dp, website, bio string, public bool) bool {
    stmt, _ := db.Prepare("update user set dp = ?, website = ?, bio = ?, public = ? where id = ?")
    _, err := stmt.Exec(dp, website, bio, public, id)
    return pe(err)
}