package server

import (
    "bytes"
    "context"
    "encoding/json"
    "github.com/mkzz115/zserve/common/log"
    "github.com/sashabaranov/go-openai"
    "io"
    "net/http"
)

var (
    tokenLimit int = 1024
)

type ChatsReq struct {
    Content    string        `json:"content"`
    HistoryMsg []HistoryChat `json:"history_msg"`
    SystemRole string        `json:"system_role,omitempty"`
    InUse      int           `json:"in_use,omitempty"`
    Lang       string        `json:"lang,omitempty"`
}

type ChatsRes struct {
    Result    string `json:"result"`
    ErrorCode string `json:"error_code"`
}

type HistoryChat struct {
    UserChat string `json:"user_chat"`
    GPTChat  string `json:"gpt_chat"`
}

func (p *ProecssorHandle) GptChatWithHistory(w http.ResponseWriter, r *http.Request) {
    buf, err := io.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    log.Info("hello is called, data[%v]", string(buf))
    req := &ChatsReq{}
    res := &ChatsRes{}
    err = json.Unmarshal(buf, req)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if req.InUse == 2 {
        res, err = p.SenseChatWithHistory(req)
    } else {
        res, err = p.ChatGPTChatWithHistory(req)
    }
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    resData, err := json.Marshal(res)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Write(resData)
}

func (p *ProecssorHandle) SenseChatWithHistory(req *ChatsReq) (*ChatsRes, error) {
    res := &ChatsRes{}
    history := []SMessage{}
    for i := 0; i < len(req.HistoryMsg); i++ {
        role1 := SMessage{
            Role:    SMessageRoleUser,
            Content: req.HistoryMsg[i].UserChat,
        }
        role2 := SMessage{
            Role:    SMessageRoleAssistant,
            Content: req.HistoryMsg[i].GPTChat,
        }
        history = append(history, role1, role2)
    }
    history = append(history, SMessage{
        Role:    SMessageRoleUser,
        Content: req.Content,
    })
    sReq := &SReq{
        Messages:     history,
        MaxNewTokens: 2048,
    }
    sReqBuff, err := json.Marshal(sReq)
    if err != nil {
        log.Error("invalid req:%v, err:%s", sReq, err.Error())
        return res, err
    }
    url := "https://lm_experience.sensetime.com/v1/nlp/chat/completions"
    httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(sReqBuff))
    if err != nil {
        log.Error("invalid req:%v, err:%s", sReqBuff, err.Error())
        return res, err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", p.SKey)
    client := &http.Client{}
    httpRes, err := client.Do(httpReq)
    if err != nil {
        log.Error("invalid req:%v, err:%s", sReqBuff, err.Error())
        return res, err
    }
    defer httpRes.Body.Close()
    httpResBuff, err := io.ReadAll(httpRes.Body)
    if err != nil {
        log.Error("invalid req:%v, err:%s", sReqBuff, err.Error())
        return res, err
    }
    log.Info("resp:%v", string(httpResBuff))
    println("resp: ", string(httpResBuff))
    sRes := &SRes{}
    err = json.Unmarshal(httpResBuff, sRes)
    if err != nil {
        log.Error("invalid req:%v, err:%s", httpResBuff, err.Error())
        return res, err
    }
    if sRes.Code != 200 {
        res.ErrorCode = sRes.Msg
        res.Result = ""
        return res, nil
    }
    for i := 0; i < len(sRes.Data.Choices); i++ {
        if sRes.Data.Choices[i].FinishReason == "stop" {
            res.Result = sRes.Data.Choices[i].Message
        }
    }

    res.ErrorCode = ""
    return res, nil
}

func (p *ProecssorHandle) ChatGPTChatWithHistory(req *ChatsReq) (*ChatsRes, error) {
    res := &ChatsRes{}

    history := []openai.ChatCompletionMessage{}
    if len(req.SystemRole) > 0 {
        history = append(history, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleSystem,
            Content: req.SystemRole,
        })
    }
    for i := 0; i < len(req.HistoryMsg); i++ {
        role1 := openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleUser,
            Content: req.HistoryMsg[i].UserChat,
        }
        role2 := openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleAssistant,
            Content: req.HistoryMsg[i].GPTChat,
        }
        history = append(history, role1, role2)
    }
    history = append(history, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: req.Content,
    })

    currentTokens := 0
    for _, h := range history {
        currentTokens += CountTokens(h.Content)
        if currentTokens > tokenLimit {
            for i := 0; i < len(history); i++ {
                currentTokens -= CountTokens(history[i].Content)
                history = history[i+1:]
                if currentTokens <= tokenLimit {
                    break
                }
            }
            break
        }
    }

    aiReq := openai.ChatCompletionRequest{
        Model:       openai.GPT3Dot5Turbo,
        Messages:    history,
        Temperature: 1.2,
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
    return res, nil
}
