package router

import (
    "github.com/kataras/iris"
    "fmt"
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

func buildErrorPage(c *iris.Context, err ErrorContainer) {
    c.Render("error.html", struct {
        Title string
        ErrorTitle string
        ErrorMessage string
    } {
        fmt.Sprintf("Error: %v", err.Title),
        err.Title,
        err.Message,
    })
}

func initErrorPages() {

    iris.OnError(iris.StatusNotFound, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "404",
            "Sorry! The page you requested could not be found.",
        })
    })
    
    iris.OnError(iris.StatusInternalServerError, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "503",
            "Our server encountered an internal server error.",
        })
    })
    
    iris.OnError(iris.StatusForbidden, func(c *iris.Context) {
        buildErrorPage(c, ErrorContainer {
            "403",
            "Please sign in first.",
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
            "This email address already exists on our database.",
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
            "The user with that email address does not exist on our database.",
        },
        "incorrect-captcha": ErrorContainer {
            "Incorrect Captcha",
            "The captcha that you entered was not correct.",
        },
    }
    
    for s, ec := range errTypes {
        iris.Handle("GET", fmt.Sprintf("/error/%v", s), ErrorContainer {
            Title: ec.Title,
            Message: ec.Message,
        })(s)
    }
}