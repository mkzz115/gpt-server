package server

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "testing"
    "time"
)

func TestProecssorHandle_GptChatWithHistory(t *testing.T) {
    c := http.Client{}
    url := "http://127.0.0.1:12308/chat"
    req1 := ChatsReq{}
    //req1.Content = "请为我写一个国企风格的工作汇报材料，字数在30字左右"
    req1.Content = "告诉我今天是几月几号，成都的天气怎么样？"
    req1.HistoryMsg = []HistoryChat{}

    resp1, err := singleReq(url, &c, req1, t)
    if err != nil {
        t.Error(err)
        return
    }
    //
    t.Logf("rep1:\n %s", resp1.Result)
    req2 := ChatsReq{}
    req2.Content = "我穿什么衣服合适呢"
    req2.HistoryMsg = []HistoryChat{
        {
            UserChat: req1.Content,
            GPTChat:  resp1.Result,
        },
    }
    resp2, err := singleReq(url, &c, req2, t)
    if err != nil {
        t.Error(err)
        return
    }
    t.Logf("\nresp2\n%s", resp2.Result)
    //req3 := ChatsReq{}
    //req3.Content = "单位年轻人比较多，风格上时髦一些"
    //req3.HistoryMsg = req2.HistoryMsg
    //req3.HistoryMsg = append(req3.HistoryMsg,
    //    HistoryChat{
    //        UserChat: req2.Content,
    //        GPTChat:  resp2.Result,
    //    })
    //resp3, err := singleReq(url, &c, req3, t)
    //if err != nil {
    //    t.Error(err)
    //    return
    //}
    ////t.Logf("\nresp3:\n %s", resp3.Result)
    //req4 := ChatsReq{}
    //req4.Content = "多引用一些诗词，文艺一点的"
    //req4.HistoryMsg = req3.HistoryMsg
    //req4.HistoryMsg = append(req4.HistoryMsg,
    //    HistoryChat{
    //        UserChat: req3.Content,
    //        GPTChat:  resp3.Result,
    //    })
    //resp4, err := singleReq(url, &c, req4, t)
    //if err != nil {
    //    t.Error(err)
    //    return
    //}
    //t.Logf("\nresp3:\n %s", resp4.Result)
}

func TestProecssorHandle_SenseChatWithHistory(t *testing.T) {
    c := http.Client{}
    url := "http://127.0.0.1:12308/chat"
    req1 := ChatsReq{}
    //req1.Content = "请为我写一个国企风格的工作汇报材料，字数在30字左右"
    req1.Content = "告诉我今天是几月几号，成都的天气怎么样？"
    req1.HistoryMsg = []HistoryChat{}
    req1.InUse = 2

    resp1, err := singleReq(url, &c, req1, t)
    if err != nil {
        t.Error(err)
        return
    }
    //
    t.Logf("rep1:\n %s", resp1.Result)
    req2 := ChatsReq{}
    req2.Content = "我穿什么衣服合适呢"
    req2.InUse = 2
    req2.HistoryMsg = []HistoryChat{
        {
            UserChat: req1.Content,
            GPTChat:  resp1.Result,
        },
    }
    time.Sleep(time.Duration(4 * time.Second))
    resp2, err := singleReq(url, &c, req2, t)
    if err != nil {
        t.Error(err)
        return
    }
    t.Logf("\nresp2\n%s", resp2.Result)
}

func singleReq(url string, c *http.Client, req2 ChatsReq, t *testing.T) (ChatsRes, error) {
    resp2 := ChatsRes{}
    buf, err := json.Marshal(req2)
    t.Logf("req: ==> %s\n", string(buf))
    res2, err := c.Post(url, "application/json", bytes.NewBuffer(buf))
    if err != nil {
        t.Error(err)
        return resp2, err
    }
    defer res2.Body.Close()
    body, err := io.ReadAll(res2.Body)
    if err != nil {
        t.Error(err)
        return resp2, err
    }
    t.Logf("raw resp: %s", string(body))
    err = json.Unmarshal(body, &resp2)
    return resp2, err
}
