package main

import (
    "flag"
    "github.com/mkzz115/gpt-server/m/server"
    "github.com/mkzz115/zserve/common/log"
)

var (
    configPath string
    listenHost string
    listenPort int
    serverName string
    logPath    string
)

func init() {
    flag.StringVar(&configPath, "c", "./conf/config.toml", "set config file path, default use conf/config.toml")
    flag.StringVar(&logPath, "d", "", "set log file path")
    flag.StringVar(&serverName, "name", "", "set server name, default use config file's name")
    flag.StringVar(&listenHost, "l", "", "set listen ip, default use config file's host")
    flag.IntVar(&listenPort, "p", 0, "set listen port, default use config's port")
}

func main() {
    conf := server.ReadConfig(configPath)
    err := log.Init(conf.LogPath)
    if err != nil {
        panic(err)
    }
    if len(logPath) > 0 {
        conf.LogPath = logPath
    }
    if len(listenHost) > 0 {
        conf.Server.Host = listenHost
    }
    if listenPort > 0 {
        conf.Server.Port = int32(listenPort)
    }
    if len(serverName) > 0 {
        conf.Server.Name = serverName
    }

    server.Start(conf)
}
