package server

import (
    "github.com/mkzz115/zserve"
    "github.com/mkzz115/zserve/common/log"
)

type HttpServe struct {
    config *Config
}

func Start(c *Config) {
    name := "Gpt Server"
    address := ":12308"

    fn := func(cf zserve.Configer) error {
        log.Info("init config")
        return nil
    }
    proc := &HttpServe{
        config: c,
    }
    zserve := zserve.NewZServer(name, address)

    zserve.Init(nil, fn, proc)
    zserve.Start()
}
