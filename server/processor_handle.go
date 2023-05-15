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
    config       *Config
}

func NewProcessorHandle(c *Config) *ProecssorHandle {
    p := &ProecssorHandle{
        config: c,
    }
    if len(c.Server.Key) > 0 {
        p.Key = c.Server.Key
    } else {
        log.Panicf("invalid key")
        os.Exit(noneKeyCode)
    }
    p.OpenaiClient = openai.NewClient(p.Key)
    return p
}
