package server

import (
    "context"
    "encoding/json"
    "github.com/mkzz115/zserve/common/log"
    "github.com/sashabaranov/go-openai"
    "io"
    "net/http"
)

type ChatReq struct {
    Content string `json:"content"`
}

type ChatRes struct {
    Result    string `json:"result"`
    ErrorCode string `json:"error_code"`
}

func (p *ProecssorHandle) GptChat(w http.ResponseWriter, r *http.Request) {

    buf, err := io.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    log.Info("hello is called, data[%v]", string(buf))
    req := &ChatReq{}
    res := &ChatRes{}
    err = json.Unmarshal(buf, req)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    println("req: ", req.Content)
    // request openai
    aiReq := openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleUser,
                Content: req.Content,
            },
        },
    }
    ctx := context.Background()
    aiResp, err := p.OpenaiClient.CreateChatCompletion(ctx, aiReq)
    if err != nil {
        log.Info("request open ai failed: %v\n", err)
        res.ErrorCode = err.Error()
        res.Result = ""
    } else {
        res.Result = aiResp.Choices[0].Message.Content
        res.ErrorCode = ""
        log.Info("resp: %v\n", aiResp.Choices[0])
    }
    println("resp: ", aiResp.Choices[0].Message.Content)

    resData, err := json.Marshal(res)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Write(resData)
}
