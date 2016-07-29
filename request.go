package gonoc

import (
    //"fmt"
    "net/http"
    //"reflect"
    //"strings"
)

type Request struct {
    W http.ResponseWriter
    R *http.Request
}

//get和post都可以获取
func (this *Request) Get(key string) string {
    this.R.ParseForm()
    if len(this.R.Form[key]) > 0 {
        return this.R.Form[key][0]
    }
    return ""
}