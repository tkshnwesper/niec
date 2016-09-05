package common

import (
    "log"
    "fmt"
    "time"
)

func Pe(err error) bool {
    if err != nil {
        log.Fatal(err)
        return false
    }
    return true
}

func GetDatetime() string {
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