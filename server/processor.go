package server

import (
    "github.com/gorilla/mux"
    "github.com/mkzz115/zserve"
    "github.com/mkzz115/zserve/common/log"
    "net/http"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        //fmt.Printf("Request method: %s\n", r.Method)
        //for key, value := range r.Header {
        //    fmt.Printf("%s: %s\n", key, value)
        //}
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        //// 1. [必须]接受指定域的请求，可以使用*不加以限制，但不安全
        //// w.Header().Set("Access-Control-Allow-Origin", "*")
        ////w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
        ////// 2. [必须]设置服务器支持的所有跨域请求的方法
        ////w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
        ////// 3. [可选]服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
        ////w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length,Token")
        ////// 4. [可选]设置XMLHttpRequest的响应对象能拿到的额外字段
        w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers,Token")
        ////// 5. [可选]是否允许后续请求携带认证信息Cookir，该值只能是true，不需要则不设置
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        //// 检查是否为预检请求
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func (h *HttpServe) Init() error {
    log.Info("httpserver init")
    return nil
}

func (h *HttpServe) Driver() (string, interface{}) {
    //r := httprouter.New()
    r := mux.NewRouter()
    handle := NewProcessorHandle(h.config)
    r.HandleFunc("/r", corsMiddleware(handle.GptChat)).Methods(http.MethodOptions, http.MethodPost)
    r.HandleFunc("/chat", corsMiddleware(handle.GptChatWithHistory)).Methods(http.MethodOptions, http.MethodPost)
    return zserve.PROCESSOR_HTTP, r
}
