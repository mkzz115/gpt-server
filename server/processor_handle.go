package server

import (
    "github.com/sashabaranov/go-openai"
    "log"
    "os"
)

var (
    noneKeyCode = -1
)

type ProecssorHandle struct {
    OpenaiClient *openai.Client
    Key          string
    SKey         string
    config       *Config
}

func NewProcessorHandle(c *Config) *ProecssorHandle {
    p := &ProecssorHandle{
        config: c,
    }
    if len(c.Server.Key) > 0 || len(c.Server.SKey) > 0 {
        p.Key = c.Server.Key
        p.SKey = c.Server.SKey
    } else {
        log.Panicf("invalid key")
        os.Exit(noneKeyCode)
    }
    p.OpenaiClient = openai.NewClient(p.Key)
    return p
}
