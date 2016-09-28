package router

import (
    "github.com/kataras/iris"
    "fmt"
    "html/template"
)

// ErrorContainer is a struct which is responsible for holding information about errors
type ErrorContainer struct {
    Title string
    Message string
}

// Serve is an interface function of iris
func (ec ErrorContainer) Serve(c *iris.Context) {
    buildErrorPage(c, ec)
}

func getCode(code int) string {
    return fmt.Sprintf("%v", code)
}

func buildErrorPage(c *iris.Context, err ErrorContainer) {
    c.Render("error.html", struct {
        Title string
        Property Property
        ErrorTitle string
        ErrorMessage template.HTML
    } {
        fmt.Sprintf("Error: %v", err.Title),
        getProperty(c),
        err.Title,
        template.HTML(err.Message),
    })
}

func initErrorPages() {

    iris.OnError(iris.StatusNotFound, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            getCode(iris.StatusNotFound),
            iris.StatusText(iris.StatusNotFound),
        })
    })
    
    iris.OnError(iris.StatusInternalServerError, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            getCode(iris.StatusInternalServerError),
            iris.StatusText(iris.StatusInternalServerError),
        })
    })
    
    iris.OnError(iris.StatusUnauthorized, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            getCode(iris.StatusUnauthorized),
            iris.StatusText(iris.StatusUnauthorized),
        })
    })
    
    iris.OnError(iris.StatusForbidden, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            getCode(iris.StatusForbidden),
            iris.StatusText(iris.StatusForbidden),
        })
    })
    
    iris.OnError(iris.StatusNoContent, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            getCode(iris.StatusNoContent),
            iris.StatusText(iris.StatusNoContent),
        })
    })
    
    errTypes := map[string]ErrorContainer {
        "blank-field": ErrorContainer {
            "Blank Field(s)",
            "Kindly enter all the required fields.",
        },
        "password-mismatch": ErrorContainer {
            "Password Mismatch",
            "Your passwords don't match, please try again.",
        },
        "email-already-exists": ErrorContainer {
            "Email Already Exists",
            "This email address already exists in our database.",
        },
        "username-already-taken": ErrorContainer {
            "Username already taken",
            "This username has already been taken by another user.",
        },
        "invalid-credentials": ErrorContainer {
            "Invalid Credentials",
            "The email or password that you entered are not valid.",
        },
        "user-does-not-exist": ErrorContainer {
            "User does not exist",
            "The user with that email address does not exist in our database.",
        },
        "incorrect-captcha": ErrorContainer {
            "Incorrect Captcha",
            "The captcha that you entered was not correct.",
        },
        "invalid-verification": ErrorContainer {
            "Invalid Verification",
            "The verification url was not valid.",
        },
        "email-not-verified": ErrorContainer {
            "Email Not Verified",
            "Please verify your email by clicking on the link in the email that we sent to you.",
        },
    }
    
    for s, ec := range errTypes {
        iris.Handle("GET", fmt.Sprintf("/error/%v", s), ec)(s)
    }
}