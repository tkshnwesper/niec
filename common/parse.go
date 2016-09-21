package db

import (
    "encoding/json"
    "io/ioutil"
)

// Config is used to unmarshal the config file
type Config struct {
    DB ConfigDB
}

// ConfigDB is used to unmarshal the DB portion of the config
type ConfigDB struct {
    Name, User, Password string
}

func parse() Config {
    data, err := ioutil.ReadFile("config.json")
    pe(err)
    var config Config
    json.Unmarshal(data, &config)
    return config
}