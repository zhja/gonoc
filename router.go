package gonoc

import (
    //"fmt"
    "net/http"
    "reflect"
    "strings"
)

type Router struct {
    Host string
    Controller string
    Action string
}

func (this *Router) Init(w http.ResponseWriter, r *http.Request){
    //初始化request
    Requests.W = w
    Requests.R = r
    //支持两级，/api/test,第一级为controller,第二级为action
    path := r.URL.Path
    path = Substr(path, 1, len(path))
    pathSlice := strings.Split(path, "/")
    if len(pathSlice) < 2 {
        panic("Route error")
    }
    controllerTitle := StringFC(pathSlice[0])
    actionTitle := StringFC(pathSlice[1])

    conCase := RouterMap[controllerTitle]
    conCaseModels := reflect.ValueOf(conCase)
    ///httpRequest := make([]reflect.Value, 2)
    //httpRequest[0] = reflect.ValueOf(w)
    //httpRequest[1] = reflect.ValueOf(r)
    _ = conCaseModels.MethodByName(actionTitle).Call(nil)
}

func Route(w http.ResponseWriter, r *http.Request) {
    rt := &Router{}
    rt.Init(w, r)
}

func RegisterRouter(key string, con interface{}){
    if len(RouterMap) == 0 {
        RouterMap = make(map[string]interface{})
    }
    RouterMap[key] = con
}