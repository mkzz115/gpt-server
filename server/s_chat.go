package server

const (
    SMessageRoleUser      = "user"
    SMessageRoleAssistant = "assistant"
)

type SReq struct {
    Messages []SMessage `json:"messages"`
    //Temperature       float64 `json:"temperature"`
    //TopP              float64 `json:"top_p"`
    MaxNewTokens int `json:"max_new_tokens,omitempty"`
    //RepetitionPenalty int     `json:"repetition_penalty,omitempty"`
    //Stream            bool    `json:"stream,omitempty"`
    //User              string  `json:"user"`
}

type SRes struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data SData  `json:"data"`
}

type SMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type SData struct {
    Id      string `json:"id"`
    Choices []struct {
        Message      string `json:"message"`
        FinishReason string `json:"finish_reason"`
    } `json:"choices"`
    Model string `json:"model"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
    Status int `json:"status"`
}
