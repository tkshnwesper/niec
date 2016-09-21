package common

import (
    "encoding/json"
    "io/ioutil"
)

// Config is used to unmarshal the config file
type Config struct {
    DB ConfigDB
    SMTP ConfigSMTP
}

// ConfigDB is used to unmarshal the DB portion of the config
type ConfigDB struct {
    Name, User, Password string
}

// ConfigSMTP is used to unmarshal the SMTP portion of the config
type ConfigSMTP struct {
    Host string
    Port int
    Username, Password, FromAlias string
}

func parse() Config {
    data, err := ioutil.ReadFile("config.json")
    Pe(err)
    var config Config
    json.Unmarshal(data, &config)
    return config
}