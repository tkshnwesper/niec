package router

import (
    "niec/db"
    "github.com/kataras/iris"
    "github.com/dchest/captcha"
)

// Property iss a container that carries all the information required by a user
type Property struct {
    LoggedIn bool
    Username string
    UserID int64
}

func initSignPages() {
    iris.Get("/sign/up", func(c *iris.Context) {
        if !isLoggedIn(c) {
            renderSign(c, "Niec :: Sign Up", "Sign Up")
        } else {
            c.RedirectTo("landing")
        }
    })("signup")
    
    iris.Post("/sign/up", func(c *iris.Context) {
        if verifyCaptcha(c) {
            if res, email, password := getCreds(c); res {
                if db.CheckEmailExists(email) {
                    c.RedirectTo("email-already-exists")
                } else {
                    c.Session().Set("email", email)
                    c.Session().Set("password", password)
                    c.RedirectTo("signup-next")
                }
            } else {
                c.RedirectTo("blank-Field")
            }
        } else {
            c.RedirectTo("incorrect-captcha")
        }
    })
    
    iris.Get("/sign/in", func(c *iris.Context) {
        if !isLoggedIn(c) {
            renderSign(c, "Niec :: Sign In", "Sign In")
        } else {
            c.RedirectTo("landing")
        }
    })("signin")
    
    iris.Post("/sign/in", func(c *iris.Context) {
        if verifyCaptcha(c) {
            if res, email, password := getCreds(c); res {
                if db.CheckEmailExists(email) {
                    creds, verified := db.VerifyCreds(email, password)
                    if creds && verified {
                        // Session cleared after successful signin
                        c.Session().Clear()
                        // c.Session().Set(common.UserIdentificationAttribute, db.GetUsername(email))
                        un, id := db.GetUsernameAndID(email)
                        c.Session().Set("property", Property { true, un, id })
                        c.RedirectTo("landing")
                    } else if !verified {
                        c.RedirectTo("email-not-verified")
                    } else {
                        c.RedirectTo("invalid-credentials")
                    }
                } else {
                    c.RedirectTo("user-does-not-exist")
                }
            } else {
                c.RedirectTo("blank-Field")
            }
        } else {
            c.RedirectTo("incorrect-captcha")
        }
    })
    
    iris.Get("/sign/up/next", func(c *iris.Context) {
        Fields := []Field {
            {
                "text",
                "username",
                "Username",
                25,
                "",
                true,
            },
            {
                "password",
                "retype",
                "Retype password",
                -1,
                "",
                true,
            },
            {
                "url",
                "dp",
                "Display picture URL (Optional)",
                255,
                "",
                true,
            },
            {
                "url",
                "website",
                "Website (Optional)",
                255,
                "",
                true,
            },
        }
        cb := []Checkbox {
            {
                "Public profile",
                "privacy",
                "public",
                "If your profile is public, it can be viewed without logging in",
                "globe",
                false,
            },
        }
        c.Render("sign.up.next.html", struct{
            Title string
            Property Property
            Fields []Field
            Checkboxes []Checkbox
        }{
            "Niec :: Sign Up - Next",
            getProperty(c),
            Fields,
            cb,
        })
    })("signup-next")
    
    iris.Post("/sign/up/next", func(c *iris.Context) {
        bio := c.FormValueString("bio")
        username := c.FormValueString("username")
        dp := c.FormValueString("dp")
        retype := c.FormValueString("retype")
        website := c.FormValueString("website")
        privacy := c.FormValueString("privacy")
        pub := false
        if privacy == "public" {
            pub = true
        }
        if c.Session().GetString("password") != retype {
            c.RedirectTo("password-mismatch")
        } else if db.CheckUsernameExists(username) {
            c.RedirectTo("username-already-taken")
        } else if c.Session().GetString("email") == "" {
            c.EmitError(iris.StatusInternalServerError)
        } else {
            if !db.InsertUser(
                c.Session().GetString("email"),
                retype,
                username,
                dp,
                bio,
                website,
                pub,
            ) {
                c.EmitError(iris.StatusInternalServerError)
            } else {
                // Session cleared after successful signup
                c.Session().Clear()
                c.SetFlash("message", "Awesome! You have signed up. Just one more step, click on the verification link that we sent you.")
                c.SetFlash("messageType", "info")
                c.RedirectTo("landing")
            }
        }
    })
    
    iris.Get("/verify/:id/:hash", func(c *iris.Context) {
        id, err := c.ParamInt64("id")
        hash := c.Param("hash")
        if pe(err) {
            if db.VerifyEmail(id, hash) {
                c.SetFlash("message", "Email verified successfully! You can now log into your account.")
                c.SetFlash("messageType", "success")
                c.RedirectTo("landing")
            } else {
                c.RedirectTo("invalid-verification")
            }
        } else {
            c.EmitError(iris.StatusNotFound)
        }
    })
}

func getCreds(c *iris.Context) (bool, string, string) {
    email := c.FormValueString("email")
    password := c.FormValueString("password")
    if email == "" || password == "" {
        return false, "", ""
    }
    return true, email, password
}

func renderSign(c *iris.Context, title, action string) {
    c.SetHeader("Cache-Control", "no-store, must-revalidate, max-age=0")
    c.SetHeader("Pragma", "no-cache");
    capid := captcha.New()
    c.Session().Set("capid", capid)
    Fields := []Field {
        {
            "email",
            "email",
            "Email Address",
            255,
            "",
            true,
        },
        {
            "password",
            "password",
            "Password",
            -1,
            "",
            true,
        },
    }
    c.Render("sign.html", struct{
        Title string
        Property Property
        Action string
        Fields []Field
        CapID string
    }{
        title,
        getProperty(c),
        action,
        Fields,
        capid,
    })
}

func verifyCaptcha(c *iris.Context) bool {
    cap := c.FormValueString("captcha")
    return captcha.VerifyString(c.Session().Get("capid").(string), cap)
}
