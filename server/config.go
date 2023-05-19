package server

import (
    "fmt"
    "github.com/BurntSushi/toml"
)

type Config struct {
    LogPath string     `toml:"log_path"`
    Mysql   MysqlConf  `toml:"mysql"`
    Server  ServerConf `toml:"server"`
}

type MysqlConf struct {
    Address  string `toml:"address"`
    User     string `toml:"user"`
    Passwd   string `toml:"password"`
    DataBase string `toml:"database"`
}

type ServerConf struct {
    Name string `toml:"name"`
    Host string `toml:"host"`
    Port int32  `toml:"port"`
    Key  string `toml:"key"`
    SKey string `toml:"s_key"`
}

func ReadConfig(path string) *Config {
    conf := &Config{}
    _, err := toml.DecodeFile(path, conf)
    if err != nil {
        panic(fmt.Sprintf("%s read error: %s", path, err.Error()))
    }
    return conf
}
